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

type cmclNodeType int

const (
	cmclIfCheck cmclNodeType = iota
	cmclForeachCheck
	cmclFieldExpr
	cmclOrExpr
	cmclAndExpr
	cmclNotExpr
)

type cmclNode struct {
	nodeType cmclNodeType
	children []*cmclNode

	// Used by cmclIfCheck
	elseIfStatements []*cmclNode
	elseStatement    *cmclNode
}

type CheckEvaluator struct {
	*parser_cmcl.BaseCMCLListener

	primaryField     string
	fields           map[string]FieldInfo
	optMissingFields map[string]bool

	stack         stack.Stack
	executionTree *cmclNode

	res      types.IType
	skipping bool
	err      error
}

func NewCheckEvaluator(primaryField string, fields map[string]FieldInfo, optMissingFields map[string]bool) *CheckEvaluator {
	// Create evaluator
	evaluator := &CheckEvaluator{
		primaryField:     primaryField,
		fields:           fields,
		optMissingFields: optMissingFields,
	}

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

// EnterExprCheck is called when production exprCheck is entered.
func (ce *CheckEvaluator) EnterExprCheck(ctx *parser_cmcl.ExprCheckContext) {}

// ExitExprCheck is called when production exprCheck is exited.
func (ce *CheckEvaluator) ExitExprCheck(ctx *parser_cmcl.ExprCheckContext) {}

// EnterIfCheck is called when production ifCheck is entered.
func (ce *CheckEvaluator) EnterIfCheck(ctx *parser_cmcl.IfCheckContext) {
	// Create new node for if statement
	newNode := &cmclNode{
		nodeType: cmclIfCheck,
		children: make([]*cmclNode, 0),
	}

	// Add node the the stack
	ce.stack.Push(newNode)

	// Add node to execution tree
	if ce.executionTree == nil { // Root node
		ce.executionTree = newNode
	} else { // Child node
		parentNode := ce.stack.Peek().(*cmclNode)
		parentNode.children = append(parentNode.children, newNode)
	}
}

// ExitIfCheck is called when production ifCheck is exited.
func (ce *CheckEvaluator) ExitIfCheck(ctx *parser_cmcl.IfCheckContext) {
	// Pop node from the stack
	ce.stack.Pop()
}

// EnterForeachCheck is called when production foreachCheck is entered.
func (ce *CheckEvaluator) EnterForeachCheck(ctx *parser_cmcl.ForeachCheckContext) {}

// ExitForeachCheck is called when production foreachCheck is exited.
func (ce *CheckEvaluator) ExitForeachCheck(ctx *parser_cmcl.ForeachCheckContext) {}

// EnterOrExpr is called when production orExpr is entered.
func (ce *CheckEvaluator) EnterOrExpr(ctx *parser_cmcl.OrExprContext) {}

// ExitOrExpr is called when production orExpr is exited.
func (ce *CheckEvaluator) ExitOrExpr(ctx *parser_cmcl.OrExprContext) {}

// EnterAndExpr is called when production andExpr is entered.
func (ce *CheckEvaluator) EnterAndExpr(ctx *parser_cmcl.AndExprContext) {}

// ExitAndExpr is called when production andExpr is exited.
func (ce *CheckEvaluator) ExitAndExpr(ctx *parser_cmcl.AndExprContext) {}

// EnterNotExpr is called when production notExpr is entered.
func (ce *CheckEvaluator) EnterNotExpr(ctx *parser_cmcl.NotExprContext) {}

// ExitNotExpr is called when production notExpr is exited.
func (ce *CheckEvaluator) ExitNotExpr(ctx *parser_cmcl.NotExprContext) {}

// EnterFieldCheck is called when production fieldCheck is entered.
func (ce *CheckEvaluator) EnterFieldCheck(ctx *parser_cmcl.FieldCheckContext) {}

// ExitFieldCheck is called when production fieldCheck is exited.
func (ce *CheckEvaluator) ExitFieldCheck(ctx *parser_cmcl.FieldCheckContext) {}

// EnterParenExpr is called when production parenExpr is entered.
func (ce *CheckEvaluator) EnterParenExpr(ctx *parser_cmcl.ParenExprContext) {}

// ExitParenExpr is called when production parenExpr is exited.
func (ce *CheckEvaluator) ExitParenExpr(ctx *parser_cmcl.ParenExprContext) {}

// EnterIf is called when production if is entered.
func (ce *CheckEvaluator) EnterIf(ctx *parser_cmcl.IfContext) {}

// ExitIf is called when production if is exited.
func (ce *CheckEvaluator) ExitIf(ctx *parser_cmcl.IfContext) {}

// EnterElseif is called when production elseif is entered.
func (ce *CheckEvaluator) EnterElseif(ctx *parser_cmcl.ElseifContext) {}

// ExitElseif is called when production elseif is exited.
func (ce *CheckEvaluator) ExitElseif(ctx *parser_cmcl.ElseifContext) {}

// EnterElse is called when production else is entered.
func (ce *CheckEvaluator) EnterElse(ctx *parser_cmcl.ElseContext) {}

// ExitElse is called when production else is exited.
func (ce *CheckEvaluator) ExitElse(ctx *parser_cmcl.ElseContext) {}

// EnterForeach is called when production foreach is entered.
func (ce *CheckEvaluator) EnterForeach(ctx *parser_cmcl.ForeachContext) {}

// ExitForeach is called when production foreach is exited.
func (ce *CheckEvaluator) ExitForeach(ctx *parser_cmcl.ForeachContext) {}

// EnterNot is called when production not is entered.
func (ce *CheckEvaluator) EnterNot(ctx *parser_cmcl.NotContext) {}

// ExitNot is called when production not is exited.
func (ce *CheckEvaluator) ExitNot(ctx *parser_cmcl.NotContext) {}

// EnterFunction is called when production function is entered.
func (ce *CheckEvaluator) EnterFunction(ctx *parser_cmcl.FunctionContext) {}

// ExitFunction is called when production function is exited.
func (ce *CheckEvaluator) ExitFunction(ctx *parser_cmcl.FunctionContext) {}

// EnterArgument is called when production argument is entered.
func (ce *CheckEvaluator) EnterArgument(ctx *parser_cmcl.ArgumentContext) {}

// ExitArgument is called when production argument is exited.
func (ce *CheckEvaluator) ExitArgument(ctx *parser_cmcl.ArgumentContext) {}

// EnterString is called when production string is entered.
func (ce *CheckEvaluator) EnterString(ctx *parser_cmcl.StringContext) {}

// ExitString is called when production string is exited.
func (ce *CheckEvaluator) ExitString(ctx *parser_cmcl.StringContext) {}

// EnterInt is called when production int is entered.
func (ce *CheckEvaluator) EnterInt(ctx *parser_cmcl.IntContext) {}

// ExitInt is called when production int is exited.
func (ce *CheckEvaluator) ExitInt(ctx *parser_cmcl.IntContext) {}

// EnterFloat is called when production float is entered.
func (ce *CheckEvaluator) EnterFloat(ctx *parser_cmcl.FloatContext) {}

// ExitFloat is called when production float is exited.
func (ce *CheckEvaluator) ExitFloat(ctx *parser_cmcl.FloatContext) {}

// EnterBoolean is called when production boolean is entered.
func (ce *CheckEvaluator) EnterBoolean(ctx *parser_cmcl.BooleanContext) {}

// ExitBoolean is called when production boolean is exited.
func (ce *CheckEvaluator) ExitBoolean(ctx *parser_cmcl.BooleanContext) {}

// EnterField is called when production field is entered.
func (ce *CheckEvaluator) EnterField(ctx *parser_cmcl.FieldContext) {}

// ExitField is called when production field is exited.
func (ce *CheckEvaluator) ExitField(ctx *parser_cmcl.FieldContext) {}
