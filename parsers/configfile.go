package parsers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Node is a node in a configuration file. Fields of type Object will be encoded as
// a map[string]*Node and fields of type Array will be encoded as a []*Node.
type Node struct {
	Type  FieldType   // Type of field
	Value interface{} // Value of field

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

type NodeKey struct {
	Segments []string
}

func (nk *NodeKey) String() string {
	if len(nk.Segments) == 0 {
		return ""
	}

	result := ""
	for _, segment := range nk.Segments {
		// Check if segment contains spaces or dots
		if strings.ContainsAny(segment, " .") {
			// Escape the segment with single quotes
			segment = fmt.Sprintf("'%s'", segment)
		}

		// Append the segment to the result
		result += fmt.Sprintf("%s.", segment)
	}

	return result[:len(result)-1]
}

func (nk *NodeKey) Join(otherKey *NodeKey) *NodeKey {
	return &NodeKey{Segments: append(nk.Segments, otherKey.Segments...)}
}

// MarshalJSON customizes the JSON output of NodeKey
func (nk *NodeKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(nk.String())
}

func (n *Node) Get(key *NodeKey) (*Node, error) {
	currentNode := n
	for _, segment := range key.Segments {
		if currentNode == nil {
			return nil, fmt.Errorf("cannot traverse nil node in path %s", key.String())
		}

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
			// We're trying to traverse an array node
			return nil, fmt.Errorf("cannot traverse array node %s in path %s", segment, key.String())

		default:
			// We're trying to traverse a leaf node
			return nil, fmt.Errorf("cannot traverse leaf node %s in path %s", segment, key.String())
		}
	}

	return currentNode, nil
}
