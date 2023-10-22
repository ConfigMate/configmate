package analyzer

import (
	"fmt"
	"strconv"

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

// errSkipCheck is an error type that is used to skip a check
type errSkipCheck struct{}

func (e errSkipCheck) Error() string {
	return "skip check"
}

type executionNode struct {
	FuncName string
	Value    types.IType
	Children []*executionNode
}

type CheckEvaluator struct {
	*parser_cmcl.BaseCMCLListener

	fields           map[string]types.IType
	optMissingFields map[string]bool

	stack stack.Stack
	err   error
}

func NewCheckEvaluator(primaryField types.IType, fields map[string]types.IType, optMissingFields map[string]bool) *CheckEvaluator {
	// Create evaluator
	evaluator := &CheckEvaluator{
		fields:           fields,
		optMissingFields: optMissingFields,
	}

	// Push primary field to stack
	evaluator.stack.Push(&executionNode{
		Value:    primaryField,
		Children: make([]*executionNode, 0),
	})

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
	// Return if error has already been encountered
	if s.err != nil {
		return
	}

	fmt.Println("EnterCheck: ", ctx.GetText())
}

// ExitCheck is called when production check is exited.
func (s *CheckEvaluator) ExitCheck(ctx *parser_cmcl.CheckContext) {
	// Return if error has already been encountered
	if s.err != nil {
		return
	}

	fmt.Println("ExitCheck: ", ctx.GetText())
}

// EnterFunction is called when production function is entered.
func (s *CheckEvaluator) EnterFunction(ctx *parser_cmcl.FunctionContext) {
	// Return if error has already been encountered
	if s.err != nil {
		return
	}

	fmt.Println("EnterFunction: ", ctx.GetText())

	// Get function name
	funcName := ctx.NAME().GetText()

	// Create a node for this function
	funcNode := &executionNode{
		FuncName: funcName,
		Children: make([]*executionNode, 0),
	}

	// Add function to the children of the current node
	currNode := s.stack.Peek().(*executionNode)
	currNode.Children = append(currNode.Children, funcNode)

	// Add function to the stack
	s.stack.Push(funcNode)
}

// ExitFunction is called when production function is exited.
func (s *CheckEvaluator) ExitFunction(ctx *parser_cmcl.FunctionContext) {
	// Return if error has already been encountered
	if s.err != nil {
		return
	}

	fmt.Println("ExitFunction: ", ctx.GetText())

	// Get and pop current node
	currNode := s.stack.Pop().(*executionNode)

	// Get function arguments
	args := make([]types.IType, len(currNode.Children))
	for i, child := range currNode.Children {
		args[i] = child.Value
	}

	// Get current node
	currNode = s.stack.Peek().(*executionNode)

	// Apply function
	res, err := currNode.Value.Checks()[currNode.FuncName](args)
	if err != nil {
		s.err = err
		return
	}

	// Update current node
	currNode.Value = res
	currNode.FuncName = ""
	currNode.Children = nil
}

// EnterString is called when production string is entered.
func (s *CheckEvaluator) EnterString(ctx *parser_cmcl.StringContext) {
	// Return if error has already been encountered
	if s.err != nil {
		return
	}

	fmt.Println("EnterString: ", ctx.GetText())

	// Get string value
	strValue := removeStringQuotes(ctx.GetText())

	// Make string type
	strType, err := types.MakeType("string", strValue)
	if err != nil {
		s.err = err
		return
	}

	// Add string type to the children of the current node
	currNode := s.stack.Peek().(*executionNode)
	currNode.Children = append(currNode.Children, &executionNode{
		Value: strType,
	})
}

// EnterInt is called when production int is entered.
func (s *CheckEvaluator) EnterInt(ctx *parser_cmcl.IntContext) {
	// Return if error has already been encountered
	if s.err != nil {
		return
	}

	fmt.Println("EnterInt: ", ctx.GetText())

	// Get int value
	intValue, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		s.err = err
		return
	}

	// Make int type
	intType, err := types.MakeType("int", intValue)
	if err != nil {
		s.err = err
		return
	}

	// Add int type to the children of the current node
	currNode := s.stack.Peek().(*executionNode)
	currNode.Children = append(currNode.Children, &executionNode{
		Value: intType,
	})
}

// EnterFloat is called when production float is entered.
func (s *CheckEvaluator) EnterFloat(ctx *parser_cmcl.FloatContext) {
	// Return if error has already been encountered
	if s.err != nil {
		return
	}

	fmt.Println("EnterFloat: ", ctx.GetText())

	// Get float value
	floatValue, err := strconv.ParseFloat(ctx.GetText(), 64)
	if err != nil {
		s.err = err
		return
	}

	// Make float type
	floatType, err := types.MakeType("float", floatValue)
	if err != nil {
		s.err = err
		return
	}

	// Add float type to the children of the current node
	currNode := s.stack.Peek().(*executionNode)
	currNode.Children = append(currNode.Children, &executionNode{
		Value: floatType,
	})
}

// EnterBoolean is called when production boolean is entered.
func (s *CheckEvaluator) EnterBoolean(ctx *parser_cmcl.BooleanContext) {
	// Return if error has already been encountered
	if s.err != nil {
		return
	}

	fmt.Println("EnterBoolean: ", ctx.GetText())

	// Get boolean value
	boolValue, err := strconv.ParseBool(ctx.GetText())
	if err != nil {
		s.err = err
		return
	}

	// Make boolean type
	boolType, err := types.MakeType("bool", boolValue)
	if err != nil {
		s.err = err
		return
	}

	// Add boolean type to the children of the current node
	currNode := s.stack.Peek().(*executionNode)
	currNode.Children = append(currNode.Children, &executionNode{
		Value: boolType,
	})
}

// EnterField is called when production field is entered.
func (s *CheckEvaluator) EnterField(ctx *parser_cmcl.FieldContext) {
	// Return if error has already been encountered
	if s.err != nil {
		return
	}

	fmt.Println("EnterField: ", ctx.GetText())

	// Get field name
	fieldName := ctx.Fieldname().GetText()

	// Find field in fields map
	if field, exists := s.fields[fieldName]; exists {
		// Create a node for this field
		fieldNode := &executionNode{
			Value:    field,
			Children: make([]*executionNode, 0),
		}

		// Add to the childe of the current node
		currNode := s.stack.Peek().(*executionNode)
		currNode.Children = append(currNode.Children, fieldNode)

		// Add field to the stack
		s.stack.Push(fieldNode)
	} else if s.optMissingFields[fieldName] {
		s.err = errSkipCheck{}
	} else {
		s.err = fmt.Errorf("field %s not found", fieldName)
	}
}

// ExitField is called when production field is exited.
func (s *CheckEvaluator) ExitField(ctx *parser_cmcl.FieldContext) {
	// Return if error has already been encountered
	if s.err != nil {
		return
	}

	fmt.Println("ExitField: ", ctx.GetText())

	// Pop field from the stack
	s.stack.Pop()
}

func removeStringQuotes(s string) string {
	return s[1 : len(s)-1]
}
