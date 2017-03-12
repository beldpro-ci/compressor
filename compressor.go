package compressor

import (
	"bytes"
)

// Compressor todo
type Compressor interface {
	// MakeBytes todo
	MakeBytes(filePaths []string) (*bytes.Buffer, error)

	// OpenBytes todo
	OpenBytes(source *bytes.Buffer, destination string) error
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
