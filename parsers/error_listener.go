package parsers

import (
	"github.com/antlr4-go/antlr/v4"
)

type CMErrorListener struct {
	antlr.DefaultErrorListener
	errors []CMParserError
}

type CMParserError struct {
	Message  string
	Location TokenLocation
}

func (d *CMErrorListener) GetErrors() []CMParserError {
	return d.errors
}

func (d *CMErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{},
	line, column int, msg string, e antlr.RecognitionException) {
	d.errors = append(d.errors, CMParserError{
		Message: msg,
		Location: TokenLocation{
			Start: CharLocation{
				Line:   line - 1,
				Column: column,
			},
			End: CharLocation{
				Line:   line - 1,
				Column: column + 1,
			},
		},
	})
}
