package analyzer

import (
	"github.com/ConfigMate/configmate/parsers"
)

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

// makeValueTokenLocation returns a TokenLocation object from a given file alias and a parsers.Node
// using the ValueLocation of the node.
func makeValueTokenLocation(fileAlias string, node *parsers.Node) TokenLocation {
	return TokenLocation{
		File:   fileAlias,
		Line:   node.ValueLocation.Line,
		Column: node.ValueLocation.Column,
		Length: node.ValueLocation.Length,
	}
}

// makeNameTokenLocation returns a TokenLocation object from a given file alias and a parsers.Node
// using the NameLocation of the node.
func makeNameTokenLocation(fileAlias string, node *parsers.Node) TokenLocation {
	return TokenLocation{
		File:   fileAlias,
		Line:   node.NameLocation.Line,
		Column: node.NameLocation.Column,
		Length: node.NameLocation.Length,
	}
}

// makeTOFTokenLocation returns a TokenLocation object from a given file alias without any specific
// line, column or length information; it sets them all to 0.
func makeTOFTokenLocation(fileAlias string) TokenLocation {
	return TokenLocation{
		File:   fileAlias,
		Line:   0,
		Column: 0,
		Length: 0,
	}
}
