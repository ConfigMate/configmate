package utils

import (
	"fmt"
	"os"
)

func ReadFile(filePath string) ([]byte, error) {
	// Read the TOML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return data, nil
}
