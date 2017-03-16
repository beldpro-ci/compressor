package compressor

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var flatFiles = []string{
	"/flat/file1",
	"/flat/file2",
	"/flat/file3",
}

var nestedFiles = []string{
	"/nested/level1/file1",
	"/nested/level1/level12/file1",
	"/nested/level1/level12/file2",
	"/nested/level1/level13/file1",
	"/nested/level1/level13/file2",
}

var compressors = []Compressor{
	TarGz,
	TarBz2,
	Tar,
}

func TestCompressors_canCompressFlatDirectory(t *testing.T) {
	t.Parallel()
	for _, compressor := range compressors {
		fixturePath, err := filepath.Abs("./fixture")
		assert.NoError(t, err)

		destDir, err := ioutil.TempDir("", "")
		assert.NoError(t, err)
		defer os.RemoveAll(destDir)

		r, w := io.Pipe()
		go func() {
			defer w.Close()
			err = compressor.MakeBytes([]string{fixturePath + "/flat"}, w)
			assert.NoError(t, err)
		}()
		err = compressor.OpenBytes(r, destDir)
		assert.NoError(t, err)

		for _, file := range flatFiles {
			_, err = os.Stat(destDir + file)
			assert.NoError(t, err)
		}
	}
}

func TestCompressors_canCompressNestedDirectory(t *testing.T) {
	t.Parallel()
	for _, compressor := range compressors {
		fixturePath, err := filepath.Abs("./fixture")
		assert.NoError(t, err)

		destDir, err := ioutil.TempDir("", "")
		assert.NoError(t, err)
		defer os.RemoveAll(destDir)

		r, w := io.Pipe()
		go func() {
			defer w.Close()
			err := compressor.MakeBytes([]string{fixturePath + "/nested"}, w)
			assert.NoError(t, err)
		}()
		err = compressor.OpenBytes(r, destDir)
		assert.NoError(t, err)

		for _, file := range nestedFiles {
			_, err = os.Stat(destDir + file)
			assert.NoError(t, err)
		}
	}
}
