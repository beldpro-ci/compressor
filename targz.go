package compressor

import (
	"archive/tar"
	gzip "github.com/klauspost/pgzip"
	"io"
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

func (tarGzFormat) MakeBytes(filePaths []string, writer io.Writer, skipperFn func(string) bool) error {
	const targzPath = "/123132121231133231"
	gzWriter := gzip.NewWriter(writer)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return tarball(filePaths, tarWriter, targzPath, skipperFn)
}

// Open untars source and puts the contents into destination.
func (tarGzFormat) OpenBytes(source io.Reader, destination string) error {
	gzr, err := gzip.NewReader(source)
	if err != nil {
		return err
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
