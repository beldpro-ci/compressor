package compressor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// base comparison test
func TestTar(t *testing.T) {
	sourceDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	destDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	fmt.Println("source=" + sourceDir)
	fmt.Println("dest=" + destDir)

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
	files, err = ioutil.ReadDir(destDir)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(files))
}

func TestTar_shouldCorrectlyProduce(t *testing.T) {
	sourceDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	destDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(sourceDir)

	fmt.Println("source=" + sourceDir)
	fmt.Println("dest=" + destDir)

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
	files, err = ioutil.ReadDir(destDir)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(files))

}
