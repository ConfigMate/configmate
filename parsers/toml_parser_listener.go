package parsers

import (
	// "strconv"
	// "strings"
	"fmt"

	parser_toml "github.com/ConfigMate/configmate/parsers/gen/parser_toml/parsers/grammars"
	"github.com/antlr4-go/antlr/v4"
	"github.com/golang-collections/collections/stack"
)

type TomlParser struct{}

type TomlErrorListener struct {
	antlr.DefaultErrorListener
	errors []error
}

// Error handling
func (s *TomlErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	s.errors = append(s.errors, fmt.Errorf("line %d:%d %s", line, column, msg))
}

// Custom TOML parser
func (p *TomlParser) Parse(data []byte) (*Node, error) {
	// Initialize the error listener
	errorListener := &TomlErrorListener{}

	input := antlr.NewInputStream(string(data))
	lexer := parser_toml.NewTOMLLexer(input)

	// Attach the error listener to the lexer
	lexer.AddErrorListener(errorListener)

	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := parser_toml.NewTOMLParser(tokenStream)

	// Attach the error listener to the parser
	parser.AddErrorListener(errorListener)

	tree := parser.Toml()

	// Check for errors after parsing
	if len(errorListener.errors) > 0 {
		return nil, fmt.Errorf("Syntax errors: %v", errorListener.errors)
	}

	walker := antlr.NewParseTreeWalker()
	tomlListener := &TomlParserListener{}
	walker.Walk(tomlListener, tree)

	return tomlListener.configFile, nil
}

type TomlParserListener struct {
	*parser_toml.BaseTOMLListener

	configFile *Node
	stack      stack.Stack
}
