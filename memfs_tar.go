package compressor

import (
	"archive/tar"
	"fmt"
	"gopkg.in/src-d/go-billy.v2"
	"path/filepath"
	"strings"
	"io"
)

// path is abspath to current file
func memFsTarFile(tarWriter *tar.Writer,source string, baseDir string, path string, fs billy.Filesystem) error {
	fmt.Printf("looking-at=%s\n", path, baseDir)
	info,err := fs.Stat(path)
	if err != nil {
		return fmt.Errorf("error walking to %s: %v", path, err)
	}

	header, err := tar.FileInfoHeader(info, path)
	if err != nil {
		return fmt.Errorf("%s: making header: %v", path, err)
	}

	if baseDir != "" {
		header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
	}

	if info.IsDir() {
		header.Name += "/"
	}

	fmt.Printf("header-name=%s\n", header.Name)
	fmt.Printf("fileinfo=%+v\n", info)
	fmt.Printf("header=%+v\n", header)
	err = tarWriter.WriteHeader(header)
	if err != nil {
		panic(err)
	}

	fmt.Printf("header written!\n")
	if info.IsDir() {
		fmt.Printf("isdirectory=true; recursing\n")
		// recurse here
		files, err := fs.ReadDir(source)
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			if err = memFsTarFile(tarWriter, source, baseDir, filepath.Join(source, f.Name()), fs); err != nil {
				panic(err)
			}
		}

		return nil
	}

	if header.Typeflag == tar.TypeReg {
		file, err := fs.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		fmt.Printf("openned the file %s\n", path)

		bytesWritten, err := io.CopyN(tarWriter, file, info.Size())
		if err != nil && err != io.EOF {
			return fmt.Errorf("%s: copying contents: %v", path, err)
		}

		fmt.Printf("contents copied to tarWriter! bytesWritten=%d\n", bytesWritten)
	}

	return nil
}

// func memFsTarFile(tarWriter *tar.Writer, baseDir string, path string, fs billy.Filesystem) error {

// source = already in memory
func memFsTarDir(tarWriter *tar.Writer, source string, fs billy.Filesystem) error {
	sourceInfo, err := fs.Stat(source)
	if err != nil {
		return fmt.Errorf("%s: stat: %v", source, err)
	}

	var baseDir string
	if sourceInfo.IsDir() {
		baseDir = filepath.Base(source)
		fmt.Printf("baseDir=%s,directory=true\n", baseDir)
	}


	if err = memFsTarFile(tarWriter, source, baseDir, source, fs); err != nil {
		fmt.Printf("just errored :(%v", err)
		return err
	}

	return nil
}
