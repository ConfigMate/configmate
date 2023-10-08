package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldType is the type of a field in a configuration file.
type FieldType int

const (
	Null FieldType = iota
	Bool
	Int
	Float
	String
	Array
	Object
)

// Node is a node in a configuration file. Fields of type Object will be encoded as
// a map[string]*Node and fields of type Array will be encoded as a []*Node.
type Node struct {
	Type      FieldType   // Type of field
	ArrayType FieldType   // Type of elements in array (if Type == Array)
	Value     interface{} // Value of field

	NameLocation struct { // Location of field name in configuration file
		Line   int
		Column int
		Length int
	}
	ValueLocation struct { // Location of field value in configuration file
		Line   int
		Column int
		Length int
	}
}

// FieldType is the type of a field in a configuration file.
type FieldType int

const (
	Bool FieldType = iota
	Int
	Float
	String
	Array
	Object
)

func (ft FieldType) String() string {
	switch ft {
	case Bool:
		return "bool"
	case Int:
		return "int"
	case Float:
		return "float"
	case String:
		return "string"
	case Array:
		return "array"
	case Object:
		return "object"
	default:
		return "unknown"
	}
}

// Get returns the node at the given path from the current node.
// Paths look like these: "server.port", "settings.users[0].name", "logLevel"
func (n *Node) Get(path string) (*Node, error) {
	// Split the key
	segments := splitKey(path)
	currentNode := n

	for i, segment := range segments {
		switch currentNode.Type {
		case Object:
			objMap, ok := currentNode.Value.(map[string]*Node)
			if !ok {
				return nil, fmt.Errorf("failed to cast %s to object value in key %s", segment, path)
			}

			if nextNode, exists := objMap[segment]; exists {
				currentNode = nextNode
			} else {
				return nil, fmt.Errorf("field %s does not exist in key %s", segment, path)
			}

		case Array:
			index, err := strconv.Atoi(segment)
			if err != nil {
				return nil, fmt.Errorf("failed to convert [%s] to int value in key %s", segment, path)
			}

			// Try to cast the value to a slice of arrayValue
			arrayValue, ok := currentNode.Value.([]*Node)
			if !ok {
				return nil, fmt.Errorf("failed to cast %s to array value in key %s", segments[i], path)
			}

			if index >= len(arrayValue) {
				return nil, fmt.Errorf("index [%d] out of bounds in key %s", index, path)
			}

			currentNode = arrayValue[index]

		default:
			// If we are here, it means we're trying to traverse a leaf node
			return nil, fmt.Errorf("cannot traverse leaf node %s in key %s", segment, path)
		}
	}

	return currentNode, nil
}

// splitPath splits the given path into a list of segments appropiate for traversing a ConfigFile.
// Paths look like these: "server.port", "settings.users[0].name", "logLevel".
// The segments are split based on the dot and the square brackets.
func splitKey(path string) []string {
	// Split the key based on the dot
	segments := strings.Split(path, ".")

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
