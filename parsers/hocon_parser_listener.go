package parsers

import (
	"fmt"

	parser_hocon "github.com/ConfigMate/configmate/parsers/gen/parser_hocon/parsers/grammars"
	"github.com/antlr4-go/antlr/v4"
	"github.com/golang-collections/collections/stack"
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
	stack      stack.Stack
}

// EnterPath is called when production path is entered.
func (l *HoconParserListener) EnterPath(ctx *parser_hocon.PathContext) {}

// ExitPath is called when production path is exited.
func (l *HoconParserListener) ExitPath(ctx *parser_hocon.PathContext) {}

// EnterKey is called when production key is entered.
func (l *HoconParserListener) EnterKey(ctx *parser_hocon.KeyContext) {}

// ExitKey is called when production key is exited.
func (l *HoconParserListener) ExitKey(ctx *parser_hocon.KeyContext) {
	fmt.Println("KEY", ctx.Path().GetText())
}

// EnterObj is called when production obj is entered.
func (l *HoconParserListener) EnterObj(ctx *parser_hocon.ObjContext) {}

// ExitObj is called when production obj is exited.
func (l *HoconParserListener) ExitObj(ctx *parser_hocon.ObjContext) {}

// EnterProperty is called when production property is entered.
func (l *HoconParserListener) EnterProperty(ctx *parser_hocon.PropertyContext) {}

// ExitProperty is called when production property is exited.
func (l *HoconParserListener) ExitProperty(ctx *parser_hocon.PropertyContext) {}

// EnterRawstring is called when production rawstring is entered.
func (l *HoconParserListener) EnterRawstring(ctx *parser_hocon.RawstringContext) {}

// ExitRawstring is called when production rawstring is exited.
func (l *HoconParserListener) ExitRawstring(ctx *parser_hocon.RawstringContext) {}

// EnterV_string is called when production v_string is entered.
func (l *HoconParserListener) EnterV_string(ctx *parser_hocon.V_stringContext) {}

// ExitV_string is called when production v_string is exited.
func (l *HoconParserListener) ExitV_string(ctx *parser_hocon.V_stringContext) {}

// EnterV_rawstring is called when production v_rawstring is entered.
func (l *HoconParserListener) EnterV_rawstring(ctx *parser_hocon.V_rawstringContext) {}

// ExitV_rawstring is called when production v_rawstring is exited.
func (l *HoconParserListener) ExitV_rawstring(ctx *parser_hocon.V_rawstringContext) {}

// EnterV_reference is called when production v_reference is entered.
func (l *HoconParserListener) EnterV_reference(ctx *parser_hocon.V_referenceContext) {}

// ExitV_reference is called when production v_reference is exited.
func (l *HoconParserListener) ExitV_reference(ctx *parser_hocon.V_referenceContext) {}

// EnterObject_begin is called when production object_begin is entered.
func (l *HoconParserListener) EnterObject_begin(ctx *parser_hocon.Object_beginContext) {}

// ExitObject_begin is called when production object_begin is exited.
func (l *HoconParserListener) ExitObject_begin(ctx *parser_hocon.Object_beginContext) {}

// EnterObject_end is called when production object_end is entered.
func (l *HoconParserListener) EnterObject_end(ctx *parser_hocon.Object_endContext) {}

// ExitObject_end is called when production object_end is exited.
func (l *HoconParserListener) ExitObject_end(ctx *parser_hocon.Object_endContext) {}

// EnterObject_data is called when production object_data is entered.
func (l *HoconParserListener) EnterObject_data(ctx *parser_hocon.Object_dataContext) {
	fmt.Println("ENTER")
}

// ExitObject_data is called when production object_data is exited.
func (l *HoconParserListener) ExitObject_data(ctx *parser_hocon.Object_dataContext) {
	fmt.Println("EXIT")
}

// EnterArray_data is called when production array_data is entered.
func (l *HoconParserListener) EnterArray_data(ctx *parser_hocon.Array_dataContext) {}

// ExitArray_data is called when production array_data is exited.
func (l *HoconParserListener) ExitArray_data(ctx *parser_hocon.Array_dataContext) {}

// EnterString_data is called when production string_data is entered.
func (l *HoconParserListener) EnterString_data(ctx *parser_hocon.String_dataContext) {}

// ExitString_data is called when production string_data is exited.
func (l *HoconParserListener) ExitString_data(ctx *parser_hocon.String_dataContext) {}

// EnterReference_data is called when production reference_data is entered.
func (l *HoconParserListener) EnterReference_data(ctx *parser_hocon.Reference_dataContext) {}

// ExitReference_data is called when production reference_data is exited.
func (l *HoconParserListener) ExitReference_data(ctx *parser_hocon.Reference_dataContext) {}

// EnterNumber_data is called when production number_data is entered.
func (l *HoconParserListener) EnterNumber_data(ctx *parser_hocon.Number_dataContext) {}

// ExitNumber_data is called when production number_data is exited.
func (l *HoconParserListener) ExitNumber_data(ctx *parser_hocon.Number_dataContext) {}

// EnterArray_begin is called when production array_begin is entered.
func (l *HoconParserListener) EnterArray_begin(ctx *parser_hocon.Array_beginContext) {}

// ExitArray_begin is called when production array_begin is exited.
func (l *HoconParserListener) ExitArray_begin(ctx *parser_hocon.Array_beginContext) {}

// EnterArray_end is called when production array_end is entered.
func (l *HoconParserListener) EnterArray_end(ctx *parser_hocon.Array_endContext) {}

// ExitArray_end is called when production array_end is exited.
func (l *HoconParserListener) ExitArray_end(ctx *parser_hocon.Array_endContext) {}

// EnterArray is called when production array is entered.
func (l *HoconParserListener) EnterArray(ctx *parser_hocon.ArrayContext) {}

// ExitArray is called when production array is exited.
func (l *HoconParserListener) ExitArray(ctx *parser_hocon.ArrayContext) {}

// EnterArray_string is called when production array_string is entered.
func (l *HoconParserListener) EnterArray_string(ctx *parser_hocon.Array_stringContext) {}

// ExitArray_string is called when production array_string is exited.
func (l *HoconParserListener) ExitArray_string(ctx *parser_hocon.Array_stringContext) {}

// EnterArray_reference is called when production array_reference is entered.
func (l *HoconParserListener) EnterArray_reference(ctx *parser_hocon.Array_referenceContext) {}

// ExitArray_reference is called when production array_reference is exited.
func (l *HoconParserListener) ExitArray_reference(ctx *parser_hocon.Array_referenceContext) {}

// EnterArray_number is called when production array_number is entered.
func (l *HoconParserListener) EnterArray_number(ctx *parser_hocon.Array_numberContext) {}

// ExitArray_number is called when production array_number is exited.
func (l *HoconParserListener) ExitArray_number(ctx *parser_hocon.Array_numberContext) {}

// EnterArray_obj is called when production array_obj is entered.
func (l *HoconParserListener) EnterArray_obj(ctx *parser_hocon.Array_objContext) {}

// ExitArray_obj is called when production array_obj is exited.
func (l *HoconParserListener) ExitArray_obj(ctx *parser_hocon.Array_objContext) {}

// EnterArray_array is called when production array_array is entered.
func (l *HoconParserListener) EnterArray_array(ctx *parser_hocon.Array_arrayContext) {}

// ExitArray_array is called when production array_array is exited.
func (l *HoconParserListener) ExitArray_array(ctx *parser_hocon.Array_arrayContext) {}

// EnterArray_value is called when production array_value is entered.
func (l *HoconParserListener) EnterArray_value(ctx *parser_hocon.Array_valueContext) {}

// ExitArray_value is called when production array_value is exited.
func (l *HoconParserListener) ExitArray_value(ctx *parser_hocon.Array_valueContext) {}