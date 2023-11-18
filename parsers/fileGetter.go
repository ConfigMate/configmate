package parsers

import (
	"io"
	"os"
)

// FileGetter is an interface that gets files
type FileGetter interface {
	GetFile(filename string) ([]byte, error)
}

// FileGetterImpl is an implementation of FileGetter
type fileGetterImpl struct{}

// NewFileGetter creates a new FileGetter
func NewFileGetter() FileGetter {
	return &fileGetterImpl{}
}

// GetFile gets a file
func (f *fileGetterImpl) GetFile(filename string) ([]byte, error) {
	// Get file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Read file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
