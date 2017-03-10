package compressor

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"strings"
)

// tarGzFormat is for tarGzFormat format
type tarGzFormat struct{}

func (tarGzFormat) Match(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".tar.gz") ||
		strings.HasSuffix(strings.ToLower(filename), ".tgz") ||
		istarGzFormat(filename)
}

func (tarGzFormat) MakeBytes(filePaths []string) (*bytes.Buffer, error) {
	const targzPath = "/123132121231133231"
	buf := new(bytes.Buffer)
	gzWriter := gzip.NewWriter(buf)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return buf, tarball(filePaths, tarWriter, targzPath)
}

// Open untars source and puts the contents into destination.
func (tarGzFormat) OpenBytes(source *bytes.Buffer, destination string) error {
	gzr, err := gzip.NewReader(source)
	if err != nil {
		return err
	}
	defer gzr.Close()

	return untar(tar.NewReader(gzr), destination)
}

// Make creates a .tar.gz file at targzPath containing
// the contents of files listed in filePaths. It works
// the same way Tar does, but with gzip compression.
func (tarGzFormat) Make(targzPath string, filePaths []string) error {
	out, err := os.Create(targzPath)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", targzPath, err)
	}
	defer out.Close()

	gzWriter := gzip.NewWriter(out)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return tarball(filePaths, tarWriter, targzPath)
}

// Open untars source and decompresses the contents into destination.
func (tarGzFormat) Open(source, destination string) error {
	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("%s: failed to open archive: %v", source, err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("%s: create new gzip reader: %v", source, err)
	}
	defer gzr.Close()

	return untar(tar.NewReader(gzr), destination)
}

// istarGzFormat checks the file has the gzip compressed Tar format header by reading
// its beginning block.
func istarGzFormat(targzPath string) bool {
	f, err := os.Open(targzPath)
	if err != nil {
		return false
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return false
	}
	defer gzr.Close()

	buf := make([]byte, tarBlockSize)
	n, err := gzr.Read(buf)
	if err != nil || n < tarBlockSize {
		return false
	}

	return hasTarHeader(buf)
}
