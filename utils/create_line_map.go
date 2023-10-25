package utils

import (
	"fmt"
	"strings"
)

func CreateLineMap(fileContent []byte) (map[int]string, error) {
	// Check if fileContent is nil or empty
	if len(fileContent) == 0 {
		return nil, fmt.Errorf("fileContent is empty or nil")
	}

	// Create a map of line numbers to lines (content)
	lineMap := make(map[int]string)
	for i, line := range strings.Split(string(fileContent), "\n") {
		lineMap[i+1] = line
	}

	return lineMap, nil
}
