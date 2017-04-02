package compressor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	Tar,
}

func trueSkipper(filePath string) bool {
	return false
}

func file1Skipper(filePath string) bool {
	return strings.HasSuffix(filePath, "file1")
}

func level12Skipper(filePath string) bool {
	fmt.Println("checking " + filePath)
	return strings.HasSuffix(filePath, "/nested/level1/level12")
}

func TestCompressors_canIgnoreCertainFilesInNestedStructure(t *testing.T) {
	var expectedFiles = []string{
		"/nested/level1/file1",
		"/nested/level1/level13/file1",
		"/nested/level1/level13/file2",
	}

	var expectedNotPresent = []string{
		"/nested/level1/level12/file1",
		"/nested/level1/level12/file2",
	}

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
			err := compressor.MakeBytes([]string{fixturePath + "/nested"}, w, level12Skipper)
			assert.NoError(t, err)
		}()
		err = compressor.OpenBytes(r, destDir)
		assert.NoError(t, err)

		for _, file := range expectedFiles {
			_, err = os.Stat(destDir + file)
			assert.NoError(t, err)
		}

		for _, file := range expectedNotPresent {
			_, err = os.Stat(destDir + file)
			assert.Error(t, err)
		}
	}
}

func TestCompressors_canIgnoreCertainFilesInFlatStructure(t *testing.T) {
	var expectedFiles = []string{
		"/flat/file2",
		"/flat/file3",
	}

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
			err = compressor.MakeBytes([]string{fixturePath + "/flat"}, w, file1Skipper)
			assert.NoError(t, err)
		}()
		err = compressor.OpenBytes(r, destDir)
		assert.NoError(t, err)

		for _, file := range expectedFiles {
			_, err = os.Stat(destDir + file)
			assert.NoError(t, err)
		}

		_, err = os.Stat(destDir + "/flat/file1")
		assert.Error(t, err)
	}
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
			err = compressor.MakeBytes([]string{fixturePath + "/flat"}, w, trueSkipper)
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
			err := compressor.MakeBytes([]string{fixturePath + "/nested"}, w, trueSkipper)
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
