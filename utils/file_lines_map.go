package utils

import (
	"strings"
)

// CreateLinesMapForFiles creates a map of file paths to line maps.
// It takes in a map of file paths to file content.
func CreateLinesMapForFiles(files map[string][]byte) map[string]map[int]string {
	// Create a map of file paths to line maps
	filesLines := make(map[string]map[int]string)

	for path, fileContent := range files {
		filesLines[path] = createLineMap(fileContent)
	}

	return filesLines
}

func createLineMap(fileContent []byte) map[int]string {
	// Create a map of line numbers to lines (content)
	lineMap := make(map[int]string)

	if len(fileContent) == 0 {
		return lineMap
	}

	for i, line := range strings.Split(string(fileContent), "\n") {
		lineMap[i+1] = line
	}

	return lineMap
}
