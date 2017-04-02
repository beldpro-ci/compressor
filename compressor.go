package compressor

import (
	"io"
)

// Compressor todo
type Compressor interface {
	// MakeBytes todo
	MakeBytes(filePaths []string, writer io.Writer, skipeprFn func(string) bool) error

	// OpenBytes todo
	OpenBytes(source io.Reader, destination string) error

	// Match todo
	Match(filename string) bool
}

// TarGz is an instantiation of the tarGzFormat
// struct that implements the Compressor interface
var TarGz = &tarGzFormat{}

// Tar is an instantiation of the Tar struct
// that implements the Compressor interface
var Tar = &tarFormat{}
