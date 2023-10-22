package parsers

import (
	"fmt"
	"strconv"
	"strings"
)

// Node is a node in a configuration file. Fields of type Object will be encoded as
// a map[string]*Node and fields of type Array will be encoded as a []*Node.
type Node struct {
	Type      FieldType   // Type of field
	ArrayType FieldType   // Type of elements in array (if Type == Array)
	Value     interface{} // Value of field

	NameLocation  TokenLocation // Location of field name in the file
	ValueLocation TokenLocation // Location of field value in the file
}

// CharLocation is the location of a character file.
type CharLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// TokenLocation is the location of a token in a file. It is made
// from the locations of the start and end character in the token.
type TokenLocation struct {
	Start CharLocation `json:"start"`
	End   CharLocation `json:"end"`
}

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

func (ft FieldType) String() string {
	switch ft {
	case Null:
		return "null"
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
// If the node is found it is returned, otherwise nil is returned.
// If the path is invalid, an error is returned.
func (n *Node) Get(field string) (*Node, error) {
	// Split the key
	segments := splitPathInSegments(field)
	currentNode := n

	for _, segment := range segments {
		switch currentNode.Type {
		case Object:
			// Cast value as map[string]*Node (unsafe)
			objMap := currentNode.Value.(map[string]*Node)

			// Check if the segment exists in the map
			if nextNode, exists := objMap[segment]; exists {
				currentNode = nextNode
			} else {
				return nil, nil
			}

		case Array:
			// Cast value as []*Node (unsafe)
			arrayValue := currentNode.Value.([]*Node)

			// Convert segment to integer index
			index, err := strconv.Atoi(segment)
			if err != nil {
				return nil, fmt.Errorf("failed to convert [%s] to int value in path %s", segment, field)
			}

			// Check if the index is out of bounds
			if index >= len(arrayValue) {
				return nil, nil
			}

			currentNode = arrayValue[index]

		default:
			// If we are here, it means we're trying to traverse a leaf node
			return nil, fmt.Errorf("cannot traverse leaf node %s in path %s", segment, field)
		}
	}

	return currentNode, nil
}

// splitPathInSegments splits the given path into a list of segments appropiate
// for traversing a ConfigFile. Paths look like these: "server.port", "settings.users[0].name", "logLevel".
// The segments are split based on the dot and the square brackets.
func splitPathInSegments(path string) []string {
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
