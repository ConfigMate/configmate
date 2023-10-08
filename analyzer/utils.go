package analyzer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ConfigMate/configmate/parsers"
)

// getNodeFromConfigFileNode returns the value of the given key from the given config file.
// If the key does not exist, the first return value will be false.
// Keys look like these: "server.port", "settings.users[0].name", "logLevel"
func getNodeFromConfigFileNode(configFile parsers.ConfigFile, key string) (*parsers.Node, error) {
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

	return currentNode, nil
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

// decodeFileValue returns the alias and path of the given file value.
// File values look like these: "file_alias.server.port", "file_alias.settings.users[0].name".
func decodeFileValue(value string) (string, string) {
	// Split the value based on the dot
	segments := strings.SplitN(value, ".", 2)

	return segments[0], segments[1]
}

// equalType returns true if the given node type is equal to the given argument type.
func equalType(nodeType parsers.FieldType, argType CheckArgType) bool {
	switch argType {
	case Int:
		return nodeType == parsers.Int
	case Float:
		return nodeType == parsers.Float
	case Bool:
		return nodeType == parsers.Bool
	case String:
		return nodeType == parsers.String
	default:
		return false
	}
}

// interpretLiteralOutput attempts to convert the given literal argument to the given type.
func interpretLiteralOutput(argType CheckArgType, argValue string) (interface{}, error) {
	switch argType {
	case Int:
		// Verify value is an integer
		if value, err := strconv.Atoi(argValue); err != nil {
			return nil, fmt.Errorf("failed to interpret %s as int: %s", argValue, err.Error())
		} else {
			return value, nil
		}
	case Float:
		// Verify value is a float
		if value, err := strconv.ParseFloat(argValue, 64); err != nil {
			return nil, fmt.Errorf("failed to interpret %s as float: %s", argValue, err.Error())
		} else {
			return value, nil
		}
	case Bool:
		// Verify value is a bool
		if value, err := strconv.ParseBool(argValue); err != nil {
			return nil, fmt.Errorf("failed to interpret %s as bool: %s", argValue, err.Error())
		} else {
			return value, nil
		}
	case String:
		// Add value
		return argValue, nil
	default:
		return nil, fmt.Errorf("unknown argument type: %s", argType.String())
	}
}
