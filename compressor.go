package compressor

import (
	"bytes"
)

// Compressor todo
type Compressor interface {
	// MakeBytes todo
	MakeBytes(filePaths []string) (*bytes.Buffer, error)

	// Make todo
	Make(tarPath string, filePaths []string) error

	// Open todo
	Open(source, destination string) error

	// OpenBytes todo
	OpenBytes(source *bytes.Buffer, destination string) error

	// Match todo
	Match(filename string) bool
}

// TarBz2 is an instantiation of the tarBz2Format
// struct that implements the Compressor interface
var TarBz2 = &tarBz2Format{}

// TarGz is an instantiation of the tarGzFormat
// struct that implements the Compressor interface
var TarGz = &tarGzFormat{}

// Tar is an instantiation of the Tar struct
// that implements the Compressor interface
var Tar = &tarFormat{}
