package compressor

import (
	"archive/tar"
	"bytes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-billy.v2"
	"gopkg.in/src-d/go-billy.v2/memfs"
	"path/filepath"
	"fmt"
	"testing"
)

func prepareMemFsFilesystem(base string, t *testing.T) billy.Filesystem {
	fs := memfs.New()

	f1, err := fs.Create(base+"/flat/file1")
	assert.NoError(t, err)
	_, err = f1.Write([]byte("a\n"))
	assert.NoError(t, err)

	f2, err := fs.Create(base+"/flat/file2")
	assert.NoError(t, err)
	_, err = f2.Write([]byte("b\n"))
	assert.NoError(t, err)

	f3, err := fs.Create(base+"/flat/file3")
	assert.NoError(t, err)
	_, err = f3.Write([]byte("c\n"))
	assert.NoError(t, err)

	return fs
}

func TestMemFs(t *testing.T) {
	tarBuf := new(bytes.Buffer)
	memTarBuf := new(bytes.Buffer)

	tarWriter := tar.NewWriter(tarBuf)
	defer tarWriter.Close()

	memFsTarWriter := tar.NewWriter(memTarBuf)
	defer memFsTarWriter.Close()

	fixturePath, err := filepath.Abs("./fixture")
	assert.NoError(t, err)

	var source = fixturePath + "/flat"
	var dest = "whatever"


	fmt.Println("------------------------")
	fmt.Println("TARFS!")
	fmt.Println("------------------------")

	err = tarFile(tarWriter, source, dest)
	assert.NoError(t, err)


	fmt.Println("------------------------")
	fmt.Println("MEMFS!")
	fmt.Println("------------------------")

	err = memFsTarDir(memFsTarWriter, fixturePath+"/flat", prepareMemFsFilesystem(fixturePath, t))
	assert.NoError(t, err)

	//	assert.Equal(t, tarBuf, memTarBuf)
}
