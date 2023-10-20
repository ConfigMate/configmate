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
	errors []error
}


// Error handling
func (s *TomlErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	s.errors = append(s.errors, fmt.Errorf("line %d:%d %s", line, column, msg))
}

func (s *TomlErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex int, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
    s.errors = append(s.errors, fmt.Errorf("Ambiguity detected between positions %d and %d", startIndex, stopIndex))
}

func (s *TomlErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex int, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
    s.errors = append(s.errors, fmt.Errorf("Attempting full context between positions %d and %d", startIndex, stopIndex))
}

func (s *TomlErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex int, stopIndex int, prediction int, configs *antlr.ATNConfigSet) {
    s.errors = append(s.errors, fmt.Errorf("Context sensitivity detected between positions %d and %d", startIndex, stopIndex))
}


// Custom TOML parser
func (p *TomlParser) Parse(data []byte) (*Node, error) {
	// Initialize the error listener
	errorListener := &TomlErrorListener{}

	input := antlr.NewInputStream(string(data))
	lexer := parser_toml.NewTOMLLexer(input)

	// Attach the error listener to the lexer
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)

	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := parser_toml.NewTOMLParser(tokenStream)

	// Attach the error listener to the parser
	parser.RemoveErrorListeners()
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


// EnterExpression is called when production expression is entered.
func (l *TomlParserListener) EnterExpression(ctx *parser_toml.ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (l *TomlParserListener) ExitExpression(ctx *parser_toml.ExpressionContext) {}

// EnterComment is called when production comment is entered.
func (l *TomlParserListener) EnterComment(ctx *parser_toml.CommentContext) {}

// ExitComment is called when production comment is exited.
func (l *TomlParserListener) ExitComment(ctx *parser_toml.CommentContext) {}

// EnterKey_value is called when production key_value is entered.
func (l *TomlParserListener) EnterKey_value(ctx *parser_toml.Key_valueContext) {}

// ExitKey_value is called when production key_value is exited.
func (l *TomlParserListener) ExitKey_value(ctx *parser_toml.Key_valueContext) {}

// EnterKey is called when production key is entered.
func (l *TomlParserListener) EnterKey(ctx *parser_toml.KeyContext) {}

// ExitKey is called when production key is exited.
func (l *TomlParserListener) ExitKey(ctx *parser_toml.KeyContext) {}

// EnterSimple_key is called when production simple_key is entered.
func (l *TomlParserListener) EnterSimple_key(ctx *parser_toml.Simple_keyContext) {}

// ExitSimple_key is called when production simple_key is exited.
func (l *TomlParserListener) ExitSimple_key(ctx *parser_toml.Simple_keyContext) {}

// EnterUnquoted_key is called when production unquoted_key is entered.
func (l *TomlParserListener) EnterUnquoted_key(ctx *parser_toml.Unquoted_keyContext) {}

// ExitUnquoted_key is called when production unquoted_key is exited.
func (l *TomlParserListener) ExitUnquoted_key(ctx *parser_toml.Unquoted_keyContext) {}

// EnterQuoted_key is called when production quoted_key is entered.
func (l *TomlParserListener) EnterQuoted_key(ctx *parser_toml.Quoted_keyContext) {}

// ExitQuoted_key is called when production quoted_key is exited.
func (l *TomlParserListener) ExitQuoted_key(ctx *parser_toml.Quoted_keyContext) {}

// EnterDotted_key is called when production dotted_key is entered.
func (l *TomlParserListener) EnterDotted_key(ctx *parser_toml.Dotted_keyContext) {}

// ExitDotted_key is called when production dotted_key is exited.
func (l *TomlParserListener) ExitDotted_key(ctx *parser_toml.Dotted_keyContext) {}

// EnterValue is called when production value is entered.
func (l *TomlParserListener) EnterValue(ctx *parser_toml.ValueContext) {}

// ExitValue is called when production value is exited.
func (l *TomlParserListener) ExitValue(ctx *parser_toml.ValueContext) {}

// EnterString is called when production string is entered.
func (l *TomlParserListener) EnterString(ctx *parser_toml.StringContext) {}

// ExitString is called when production string is exited.
func (l *TomlParserListener) ExitString(ctx *parser_toml.StringContext) {}

// EnterInteger is called when production integer is entered.
func (l *TomlParserListener) EnterInteger(ctx *parser_toml.IntegerContext) {}

// ExitInteger is called when production integer is exited.
func (l *TomlParserListener) ExitInteger(ctx *parser_toml.IntegerContext) {}

// EnterFloating_point is called when production floating_point is entered.
func (l *TomlParserListener) EnterFloating_point(ctx *parser_toml.Floating_pointContext) {}

// ExitFloating_point is called when production floating_point is exited.
func (l *TomlParserListener) ExitFloating_point(ctx *parser_toml.Floating_pointContext) {}

// EnterBool is called when production bool is entered.
func (l *TomlParserListener) EnterBool(ctx *parser_toml.BoolContext) {}

// ExitBool is called when production bool is exited.
func (l *TomlParserListener) ExitBool(ctx *parser_toml.BoolContext) {}

// EnterDate_time is called when production date_time is entered.
func (l *TomlParserListener) EnterDate_time(ctx *parser_toml.Date_timeContext) {}

// ExitDate_time is called when production date_time is exited.
func (l *TomlParserListener) ExitDate_time(ctx *parser_toml.Date_timeContext) {}

// EnterArray is called when production array is entered.
func (l *TomlParserListener) EnterArray(ctx *parser_toml.ArrayContext) {}

// ExitArray is called when production array is exited.
func (l *TomlParserListener) ExitArray(ctx *parser_toml.ArrayContext) {}

// EnterArray_values is called when production array_values is entered.
func (l *TomlParserListener) EnterArray_values(ctx *parser_toml.Array_valuesContext) {}

// ExitArray_values is called when production array_values is exited.
func (l *TomlParserListener) ExitArray_values(ctx *parser_toml.Array_valuesContext) {}

// EnterComment_or_nl is called when production comment_or_nl is entered.
func (l *TomlParserListener) EnterComment_or_nl(ctx *parser_toml.Comment_or_nlContext) {}

// ExitComment_or_nl is called when production comment_or_nl is exited.
func (l *TomlParserListener) ExitComment_or_nl(ctx *parser_toml.Comment_or_nlContext) {}

// EnterTable is called when production table is entered.
func (l *TomlParserListener) EnterTable(ctx *parser_toml.TableContext) {}

// ExitTable is called when production table is exited.
func (l *TomlParserListener) ExitTable(ctx *parser_toml.TableContext) {}

// EnterStandard_table is called when production standard_table is entered.
func (l *TomlParserListener) EnterStandard_table(ctx *parser_toml.Standard_tableContext) {}

// ExitStandard_table is called when production standard_table is exited.
func (l *TomlParserListener) ExitStandard_table(ctx *parser_toml.Standard_tableContext) {}

// EnterInline_table is called when production inline_table is entered.
func (l *TomlParserListener) EnterInline_table(ctx *parser_toml.Inline_tableContext) {}

// ExitInline_table is called when production inline_table is exited.
func (l *TomlParserListener) ExitInline_table(ctx *parser_toml.Inline_tableContext) {}

// EnterInline_table_keyvals is called when production inline_table_keyvals is entered.
func (l *TomlParserListener) EnterInline_table_keyvals(ctx *parser_toml.Inline_table_keyvalsContext) {}

// ExitInline_table_keyvals is called when production inline_table_keyvals is exited.
func (l *TomlParserListener) ExitInline_table_keyvals(ctx *parser_toml.Inline_table_keyvalsContext) {}

// EnterInline_table_keyvals_non_empty is called when production inline_table_keyvals_non_empty is entered.
func (l *TomlParserListener) EnterInline_table_keyvals_non_empty(ctx *parser_toml.Inline_table_keyvals_non_emptyContext) {}

// ExitInline_table_keyvals_non_empty is called when production inline_table_keyvals_non_empty is exited.
func (l *TomlParserListener) ExitInline_table_keyvals_non_empty(ctx *parser_toml.Inline_table_keyvals_non_emptyContext) {}

// EnterArray_table is called when production array_table is entered.
func (l *TomlParserListener) EnterArray_table(ctx *parser_toml.Array_tableContext) {}

// ExitArray_table is called when production array_table is exited.
func (l *TomlParserListener) ExitArray_table(ctx *parser_toml.Array_tableContext) {}


// Helper functions
func (l *TomlParserListener) getTOS() *Node {
	return l.stack.Peek().(*Node)
}

func (l *TomlParserListener) getTOSObject() map[string]*Node {
	return l.getTOS().Value.(map[string]*Node)
}

func (l *TomlParserListener) getTOSArray() []*Node {
	return l.getTOS().Value.([]*Node)
}