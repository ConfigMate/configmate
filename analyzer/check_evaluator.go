package analyzer

import (
	"fmt"

	"github.com/ConfigMate/configmate/analyzer/types"
	parser_cmcl "github.com/ConfigMate/configmate/parsers/gen/parser_cmcl/parsers/grammars"
	"github.com/antlr4-go/antlr/v4"
	"github.com/golang-collections/collections/stack"
)

type CMCLErrorListener struct {
	*antlr.DefaultErrorListener
	errors []error
}

func (d *CMCLErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{},
	line, column int, msg string, e antlr.RecognitionException) {
	d.errors = append(d.errors, fmt.Errorf("line %d:%d %s", line, column, msg))
}

// skipCheck is a type that is returned at any
// point it is determined the check should be skipped.
type skipCheck struct{}

// evaluatorError is a type that is returned at any
// point an error is encountered.
type evaluatorError struct {
	err error
}

type CheckEvaluator struct {
	*parser_cmcl.BaseCMCLListener

	fields           map[string]types.IType
	optMissingFields map[string]bool
	stack            stack.Stack
}

func NewCheckEvaluator(primaryField types.IType, fields map[string]types.IType, optMissingFields map[string]bool) *CheckEvaluator {
	evaluator := &CheckEvaluator{
		fields:           fields,
		optMissingFields: optMissingFields,
	}
	evaluator.stack.Push(primaryField)
	return evaluator
}

func (v *CheckEvaluator) Evaluate(check string) (res types.IType, skipped bool, err error) {
	// Parse check
	input := antlr.NewInputStream(check)
	lexer := parser_cmcl.NewCMCLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser_cmcl.NewCMCLParser(stream)

	// Add error listener
	errorListener := &CMCLErrorListener{}
	p.RemoveErrorListeners()
	p.AddErrorListener(errorListener)

	tree := p.Check()

	// Check for errors
	if len(errorListener.errors) > 0 {
		return nil, false, fmt.Errorf("syntax errors: %v", errorListener.errors)
	}

	// Evaluate check
	walker := antlr.NewParseTreeWalker()
	walker.Walk(v, tree)

	return
}

// EnterCheck is called when production check is entered.
func (s *CheckEvaluator) EnterCheck(ctx *parser_cmcl.CheckContext) {
	fmt.Println("EnterCheck: ", ctx.GetText())
}

// ExitCheck is called when production check is exited.
func (s *CheckEvaluator) ExitCheck(ctx *parser_cmcl.CheckContext) {
	fmt.Println("ExitCheck: ", ctx.GetText())
}

// EnterFunction is called when production function is entered.
func (s *CheckEvaluator) EnterFunction(ctx *parser_cmcl.FunctionContext) {
	fmt.Println("EnterFunction: ", ctx.GetText())
}

// ExitFunction is called when production function is exited.
func (s *CheckEvaluator) ExitFunction(ctx *parser_cmcl.FunctionContext) {
	fmt.Println("ExitFunction: ", ctx.GetText())
}

// EnterArgument is called when production argument is entered.
func (s *CheckEvaluator) EnterArgument(ctx *parser_cmcl.ArgumentContext) {
	fmt.Println("EnterArgument: ", ctx.GetText())
}

// ExitArgument is called when production argument is exited.
func (s *CheckEvaluator) ExitArgument(ctx *parser_cmcl.ArgumentContext) {
	fmt.Println("ExitArgument: ", ctx.GetText())
}

// EnterString is called when production string is entered.
func (s *CheckEvaluator) EnterString(ctx *parser_cmcl.StringContext) {
	fmt.Println("EnterString: ", ctx.GetText())
}

// ExitString is called when production string is exited.
func (s *CheckEvaluator) ExitString(ctx *parser_cmcl.StringContext) {
	fmt.Println("ExitString: ", ctx.GetText())
}

// EnterInt is called when production int is entered.
func (s *CheckEvaluator) EnterInt(ctx *parser_cmcl.IntContext) {
	fmt.Println("EnterInt: ", ctx.GetText())
}

// ExitInt is called when production int is exited.
func (s *CheckEvaluator) ExitInt(ctx *parser_cmcl.IntContext) {
	fmt.Println("ExitInt: ", ctx.GetText())
}

// EnterFloat is called when production float is entered.
func (s *CheckEvaluator) EnterFloat(ctx *parser_cmcl.FloatContext) {
	fmt.Println("EnterFloat: ", ctx.GetText())
}

// ExitFloat is called when production float is exited.
func (s *CheckEvaluator) ExitFloat(ctx *parser_cmcl.FloatContext) {
	fmt.Println("ExitFloat: ", ctx.GetText())
}

// EnterBoolean is called when production boolean is entered.
func (s *CheckEvaluator) EnterBoolean(ctx *parser_cmcl.BooleanContext) {
	fmt.Println("EnterBoolean: ", ctx.GetText())
}

// ExitBoolean is called when production boolean is exited.
func (s *CheckEvaluator) ExitBoolean(ctx *parser_cmcl.BooleanContext) {
	fmt.Println("ExitBoolean: ", ctx.GetText())
}

// EnterField is called when production field is entered.
func (s *CheckEvaluator) EnterField(ctx *parser_cmcl.FieldContext) {
	fmt.Println("EnterField: ", ctx.GetText())
}

// ExitField is called when production field is exited.
func (s *CheckEvaluator) ExitField(ctx *parser_cmcl.FieldContext) {
	fmt.Println("ExitField: ", ctx.GetText())
}

// EnterField_function is called when production field_function is entered.
func (s *CheckEvaluator) EnterField_function(ctx *parser_cmcl.Field_functionContext) {
	fmt.Println("EnterField_function: ", ctx.GetText())
}

// ExitField_function is called when production field_function is exited.
func (s *CheckEvaluator) ExitField_function(ctx *parser_cmcl.Field_functionContext) {
	fmt.Println("ExitField_function: ", ctx.GetText())
}

func removeStringQuotes(s string) string {
	return s[1 : len(s)-1]
}
