package compressor

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// base comparison test
func TestTar(t *testing.T) {
	t.Parallel()
	sourceDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	destDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(destDir)

	file1, err := ioutil.TempFile(sourceDir, "")
	assert.NoError(t, err)
	file2, err := ioutil.TempFile(sourceDir, "")
	assert.NoError(t, err)

	err = exec.Command("tar", "-cf", destDir+"/tar.tar", "-C", sourceDir,
		filepath.Base(file1.Name()), filepath.Base(file2.Name())).Run()
	assert.NoError(t, err)

	files, err := ioutil.ReadDir(destDir)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(files))

	err = exec.Command("tar", "-xf", destDir+"/tar.tar", "-C", destDir).Run()
	assert.NoError(t, err)
	files, err = ioutil.ReadDir(destDir)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(files))
}

func TestTarOpen_canUntarNormalFile(t *testing.T) {
	t.Parallel()
	sourceDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	destDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(destDir)

	file1, err := ioutil.TempFile(sourceDir, "")
	assert.NoError(t, err)
	file2, err := ioutil.TempFile(sourceDir, "")
	assert.NoError(t, err)

	err = exec.Command("tar", "-cf", destDir+"/tar.tar", "-C", sourceDir,
		filepath.Base(file1.Name()), filepath.Base(file2.Name())).Run()
	assert.NoError(t, err)

	files, err := ioutil.ReadDir(destDir)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(files))

	err = Tar.Open(destDir+"/tar.tar", destDir)
	assert.NoError(t, err)
	files, err = ioutil.ReadDir(destDir)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(files))
}

func TestTarMake_shouldCorrectlyProduceTar(t *testing.T) {
	t.Parallel()
	sourceDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	destDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(destDir)

	file1, err := ioutil.TempFile(sourceDir, "")
	assert.NoError(t, err)
	file2, err := ioutil.TempFile(sourceDir, "")
	assert.NoError(t, err)

	err = Tar.Make(destDir+"/tar.tar", []string{file1.Name(), file2.Name()})
	assert.NoError(t, err)

	files, err := ioutil.ReadDir(destDir)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(files))

	err = exec.Command("tar", "-xf", destDir+"/tar.tar", "-C", destDir).Run()
	assert.NoError(t, err)
	files, err = ioutil.ReadDir(destDir)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(files))
}

func trueSkipper(string) bool {
	return false
}

func TestTarMakeBytes_shouldCorrectlyProduce(t *testing.T) {
	t.Parallel()
	sourceDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	destDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(destDir)

	file1, err := ioutil.TempFile(sourceDir, "")
	assert.NoError(t, err)
	file2, err := ioutil.TempFile(sourceDir, "")
	assert.NoError(t, err)

	var buffer = new(bytes.Buffer)
	err = Tar.MakeBytes([]string{file1.Name(), file2.Name()}, buffer, trueSkipper)
	assert.NoError(t, err)

	out, err := os.Create(destDir + "/tar.tar")
	assert.NoError(t, err)
	defer out.Close()

	_, err = io.Copy(out, buffer)
	assert.NoError(t, err)
	out.Sync()

	files, err := ioutil.ReadDir(destDir)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(files))

	err = exec.Command("tar", "-xf", destDir+"/tar.tar", "-C", destDir).Run()
	assert.NoError(t, err)
	files, err = ioutil.ReadDir(destDir)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(files))
}

func TestTarOpenBytes_shouldCorrectlyUntarBytes(t *testing.T) {
	t.Parallel()
	sourceDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	destDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(destDir)

	file1, err := ioutil.TempFile(sourceDir, "")
	assert.NoError(t, err)
	file2, err := ioutil.TempFile(sourceDir, "")
	assert.NoError(t, err)

	var buffer = new(bytes.Buffer)
	err = Tar.MakeBytes([]string{file1.Name(), file2.Name()}, buffer, trueSkipper)
	assert.NoError(t, err)

	err = Tar.OpenBytes(buffer, destDir)
	assert.NoError(t, err)
	files, err := ioutil.ReadDir(destDir)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(files))
}
