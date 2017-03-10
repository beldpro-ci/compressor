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
}
