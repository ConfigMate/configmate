package parsers

import (
	"fmt"

	"github.com/ConfigMate/configmate/parsers/gen/parser_hocon"
	"github.com/antlr4-go/antlr/v4"
	// "github.com/golang-collections/collections/stack"
)

type HoconParser struct{}

type HoconErrorListener struct {
	antlr.DefaultErrorListener
	errors []error
}

// Error handling
func (s *HoconErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	s.errors = append(s.errors, fmt.Errorf("line %d:%d %s", line, column, msg))
}

// Custom Json Parser
func (p *HoconParser) Parse(data []byte) (*Node, error) {
	// Initialize the error listener
	errorListener := &HoconErrorListener{}

	input := antlr.NewInputStream(string(data))
	lexer := parser_hocon.NewHOCONLexer(input)

	// Attach the error listener to the lexer
	lexer.AddErrorListener(errorListener)

	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := parser_hocon.NewHOCONParser(tokenStream)

	// Attach the error listener to the parser
	parser.AddErrorListener(errorListener)

	tree := parser.Hocon()

	// Check for errors after parsing
	if len(errorListener.errors) > 0 {
		return nil, fmt.Errorf("Syntax errors: %v", errorListener.errors)
	}

	walker := antlr.NewParseTreeWalker()
	hoconListener := &HoconParserListener{}
	walker.Walk(hoconListener, tree)

	return hoconListener.configFile, nil
}

type HoconParserListener struct {
	*parser_hocon.BaseHOCONListener

	configFile *Node
	// stack      stack.Stack
}
