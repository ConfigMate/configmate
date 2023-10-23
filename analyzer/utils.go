package analyzer

import (
	"errors"
	"strings"

	"github.com/ConfigMate/configmate/parsers"
)

// splitFileAliasAndPath splits a field represented as a string into a file alias and a path.
// The separator used is a colon (:).
func splitFileAliasAndPath(field string) (fileAlias string, path string, err error) {
	split := strings.Split(field, ":")
	if len(split) != 2 {
		return "", "", errors.New("invalid field format: " + field)
	}

	return split[0], split[1], nil
}

// makeValueTokenLocation returns a TokenLocation object from a given file alias and a parsers.Node
// using the ValueLocation of the node.
func makeValueTokenLocation(fileAlias string, node *parsers.Node) TokenLocationWithFileAlias {
	return TokenLocationWithFileAlias{
		File:     fileAlias,
		Location: node.ValueLocation,
	}
}

// // makeNameTokenLocation returns a TokenLocation object from a given file alias and a parsers.Node
// // using the NameLocation of the node.
// func makeNameTokenLocation(fileAlias string, node *parsers.Node) TokenLocationWithFileAlias {
// 	return TokenLocationWithFileAlias{
// 		File:     fileAlias,
// 		Location: node.NameLocation,
// 	}
// }

// // makeTOFTokenLocation returns a TokenLocation object from a given file alias without any specific
// // line, column or length information; it sets them all to 0.
// func makeTOFTokenLocation(fileAlias string) TokenLocationWithFileAlias {
// 	return TokenLocationWithFileAlias{
// 		File: fileAlias,
// 		Location: parsers.TokenLocation{
// 			Start: parsers.CharLocation{Line: 0, Column: 0},
// 			End:   parsers.CharLocation{Line: 0, Column: 0},
// 		},
// 	}
// }
