package files

import (
	"io"
	"os"
)

// FileFetcher is an interface that returns the
// content of files given a path
type FileFetcher interface {
	FetchFile(filename string) ([]byte, error)
}

// FileGetterImpl is an implementation of FileGetter
type fileFetcherImpl struct{}

// NewFileFetcher creates a new FileGetter
func NewFileFetcher() FileFetcher {
	return &fileFetcherImpl{}
}

// GetFile gets a file
func (f *fileFetcherImpl) FetchFile(filename string) ([]byte, error) {
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
