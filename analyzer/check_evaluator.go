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

type executionNode struct {
	FuncName string
	Value    types.IType
	CheckErr error
	Children []*executionNode
}

type CheckEvaluator struct {
	*parser_cmcl.BaseCMCLListener

	fields           map[string]FieldInfo
	optMissingFields map[string]bool
	stack            stack.Stack

	res      types.IType
	skipping bool
	err      error
}

func NewCheckEvaluator(primaryField types.IType, fields map[string]FieldInfo, optMissingFields map[string]bool) *CheckEvaluator {
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

func (ce *CheckEvaluator) Evaluate(check string) (types.IType, bool, error) {
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
	walker.Walk(ce, tree)

	return ce.res, ce.skipping, ce.err
}

func (ce *CheckEvaluator) ExitCheck(ctx *parser_cmcl.CheckContext) {
	// Return if error has already been encountered
	if ce.err != nil {
		return
	}

	// Check that the stack only contains the root node.
	// Otherwise this indicates the exiting check belongs
	// to some parameter, not the primary field.
	if ce.stack.Len() != 1 {
		return
	}

	// Get result node
	resNode := ce.stack.Pop().(*executionNode)

	// Check that the result is a boolean
	if resNode.Value.TypeName() != "bool" {
		ce.err = fmt.Errorf("check must return a boolean")
	} else {
		ce.res = resNode.Value
		ce.err = resNode.CheckErr
	}
}

// EnterFunction is called when production function is entered.
func (ce *CheckEvaluator) EnterFunction(ctx *parser_cmcl.FunctionContext) {
	// Return if error has already been encountered
	if ce.err != nil {
		return
	}

	// Get function name
	funcName := ctx.NAME().GetText()

	// Create a node for this function
	funcNode := &executionNode{
		FuncName: funcName,
		Children: make([]*executionNode, 0),
	}

	// Add function to the children of the current node
	currNode := ce.stack.Peek().(*executionNode)
	currNode.Children = append(currNode.Children, funcNode)

	// Add function to the stack
	ce.stack.Push(funcNode)
}

// ExitFunction is called when production function is exited.
func (ce *CheckEvaluator) ExitFunction(ctx *parser_cmcl.FunctionContext) {
	// Return if error has already been encountered
	if ce.err != nil {
		return
	}

	// Get and pop current node (function)
	funcNode := ce.stack.Pop().(*executionNode)

	// Get function name
	funcName := funcNode.FuncName

	// Get function arguments
	args := make([]types.IType, len(funcNode.Children))
	for i, child := range funcNode.Children {
		args[i] = child.Value
	}

	// Get current node (node the function is applied to)
	currNode := ce.stack.Peek().(*executionNode)

	// Find function
	f, exists := currNode.Value.Checks()[funcName]
	if !exists {
		ce.err = fmt.Errorf("function %s.%s not found", currNode.Value.TypeName(), funcName)
		return
	}

	// Apply function
	res, err := f(args)
	if res == nil {
		ce.err = err
		return
	}

	// Update current node with result
	currNode.FuncName = ""
	currNode.Value = res
	currNode.CheckErr = err
	currNode.Children = nil
}

// EnterField is called when production field is entered.
func (ce *CheckEvaluator) EnterField(ctx *parser_cmcl.FieldContext) {
	// Return if error has already been encountered
	if ce.err != nil {
		return
	}

	// Get field name
	fieldName := ctx.Fieldname().GetText()

	// Find field in fields map
	if field, exists := ce.fields[fieldName]; exists {
		// Create a node for this field
		fieldNode := &executionNode{
			Value:    field.Value,
			Children: make([]*executionNode, 0),
		}

		// Add to the childe of the current node
		currNode := ce.stack.Peek().(*executionNode)
		currNode.Children = append(currNode.Children, fieldNode)

		// Add field to the stack
		ce.stack.Push(fieldNode)
	} else if ce.optMissingFields[fieldName] {
		ce.skipping = true
		ce.err = fmt.Errorf("skipping check because referenced optional field %s is missing", fieldName)
	} else {
		ce.err = fmt.Errorf("field %s not found", fieldName)
	}
}

// ExitField is called when production field is exited.
func (ce *CheckEvaluator) ExitField(ctx *parser_cmcl.FieldContext) {
	// Return if error has already been encountered
	if ce.err != nil {
		return
	}

	// Pop field from the stack
	ce.stack.Pop()
}

// EnterString is called when production string is entered.
func (ce *CheckEvaluator) EnterString(ctx *parser_cmcl.StringContext) {
	// Return if error has already been encountered
	if ce.err != nil {
		return
	}

	// Get string value
	strValue := removeStringQuotes(ctx.GetText())

	// Make string type
	strType, err := types.MakeType("string", strValue)
	if err != nil {
		ce.err = err
		return
	}

	// Add string type to the children of the current node
	currNode := ce.stack.Peek().(*executionNode)
	currNode.Children = append(currNode.Children, &executionNode{
		Value: strType,
	})
}

// EnterInt is called when production int is entered.
func (ce *CheckEvaluator) EnterInt(ctx *parser_cmcl.IntContext) {
	// Return if error has already been encountered
	if ce.err != nil {
		return
	}

	// Get int value
	intValue, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		ce.err = err
		return
	}

	// Make int type
	intType, err := types.MakeType("int", intValue)
	if err != nil {
		ce.err = err
		return
	}

	// Add int type to the children of the current node
	currNode := ce.stack.Peek().(*executionNode)
	currNode.Children = append(currNode.Children, &executionNode{
		Value: intType,
	})
}

// EnterFloat is called when production float is entered.
func (ce *CheckEvaluator) EnterFloat(ctx *parser_cmcl.FloatContext) {
	// Return if error has already been encountered
	if ce.err != nil {
		return
	}

	// Get float value
	floatValue, err := strconv.ParseFloat(ctx.GetText(), 64)
	if err != nil {
		ce.err = err
		return
	}

	// Make float type
	floatType, err := types.MakeType("float", floatValue)
	if err != nil {
		ce.err = err
		return
	}

	// Add float type to the children of the current node
	currNode := ce.stack.Peek().(*executionNode)
	currNode.Children = append(currNode.Children, &executionNode{
		Value: floatType,
	})
}

// EnterBoolean is called when production boolean is entered.
func (ce *CheckEvaluator) EnterBoolean(ctx *parser_cmcl.BooleanContext) {
	// Return if error has already been encountered
	if ce.err != nil {
		return
	}

	// Get boolean value
	boolValue, err := strconv.ParseBool(ctx.GetText())
	if err != nil {
		ce.err = err
		return
	}

	// Make boolean type
	boolType, err := types.MakeType("bool", boolValue)
	if err != nil {
		ce.err = err
		return
	}

	// Add boolean type to the children of the current node
	currNode := ce.stack.Peek().(*executionNode)
	currNode.Children = append(currNode.Children, &executionNode{
		Value: boolType,
	})
}

func removeStringQuotes(s string) string {
	return s[1 : len(s)-1]
}
