package analyzer

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ConfigMate/configmate/parsers"
)

type FileFormat int

const (
	Unknown FileFormat = iota
	HOCON
	JSON
	TOML
	YAML
)

// getFileFormat returns the file format of the given filename.
func getFileFormat(filename string) FileFormat {
	ext := filepath.Ext(filename)
	switch ext {
	case ".hocon":
		return HOCON
	case ".json":
		return JSON
	case ".toml":
		return TOML
	case ".yaml":
		return YAML
	case ".yml":
		return YAML
	default:
		return Unknown
	}
}

// getValueFromConfigFile returns the value of the given key from the given config file.
// If the key does not exist, the second return value will be false.
// Keys look like these: "server.port", "settings.users[0].name", "logLevel"
func getValueFromConfigFile(configFile parsers.ConfigFile, key string) (interface{}, error) {
	// Split the key
	segments := splitKey(key)
	currentNode := configFile

	for i, segment := range segments {
		switch currentNode.Type {
		case parsers.Object:
			objMap, ok := currentNode.Value.(map[string]*parsers.Node)
			if !ok {
				return nil, fmt.Errorf("failed to cast %s to object value in key %s", segment, key)
			}

			if nextNode, exists := objMap[segment]; exists {
				currentNode = nextNode
			} else {
				return nil, fmt.Errorf("field %s does not exist in key %s", segment, key)
			}

		case parsers.Array:
			index, err := strconv.Atoi(segment)
			if err != nil {
				return nil, fmt.Errorf("failed to convert [%s] to int value in key %s", segment, key)
			}

			// Try to cast the value to a slice of arrayValue
			arrayValue, ok := currentNode.Value.([]*parsers.Node)
			if !ok {
				return nil, fmt.Errorf("failed to cast %s to array value in key %s", segments[i], key)
			}

			if index >= len(arrayValue) {
				return nil, fmt.Errorf("index [%d] out of bounds in key %s", index, key)
			}

			currentNode = arrayValue[index]

		default:
			// If we are here, it means we're trying to traverse a leaf node
			return nil, fmt.Errorf("cannot traverse leaf node %s in key %s", segment, key)
		}
	}

	return currentNode.Value, nil
}

// splitKey splits the given key into a list of segments.
// Keys look like these: "server.port", "settings.users[0].name", "logLevel"
func splitKey(key string) []string {
	// Split the key based on the dot
	segments := strings.Split(key, ".")

	// Split array accesses into separate segments
	for i := 0; i < len(segments); i++ {
		segment := segments[i]
		if strings.Contains(segment, "[") {
			// Split the segment into two segments
			key := segment[:strings.Index(segment, "[")]
			indexStr := segment[strings.Index(segment, "[")+1 : len(segment)-1]

			// Check if there is an array name. In the case that a document has an
			// array as the top level object, the array name will be empty and the first
			// segment will be just the index in square brackets (eg. "[0].name").
			if len(key) == 0 {
				segments[i] = indexStr
				continue
			}

			// Insert the new segments into the slice
			segments[i] = key
			segments = append(segments, "")
			copy(segments[i+2:], segments[i+1:])
			segments[i+1] = indexStr
		}
	}

	return segments
}

// decodeFileValue returns true if the given value is a file value. It also returns the file alias
// and the key of the value in the file.
// File values look like these: "f:file_alias.server.port", "f:file_alias.settings.users[0].name".
func decodeFileValue(value string) (bool, string, string) {
	// Check if the value starts with "f:"
	if !strings.HasPrefix(value, "f:") {
		return false, "", ""
	}

	// Remove the "f:" prefix
	value = value[2:]

	// Split value into file alias and key
	segments := strings.SplitN(value, ".", 2)

	// Check if there are two none empty segments
	if len(segments) != 2 || len(segments[0]) == 0 || len(segments[1]) == 0 {
		return false, "", ""
	}

	return true, segments[0], segments[1]
}

// decodeLiteralValue returns true if the given value is a literal value.
// Literal values look like these: "l:100", "l:hello world", "l:true".
func decodeLiteralValue(value string) (bool, string) {
	// Check if the value starts with "l:"
	if !strings.HasPrefix(value, "l:") {
		return false, ""
	}

	// Remove the "l:" prefix
	value = value[2:]

	return true, value
}
