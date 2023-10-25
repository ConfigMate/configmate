package analyzer

import (
	"fmt"

	parser_cmcl "github.com/ConfigMate/configmate/parsers/gen/parser_cmcl/parsers/grammars"
	"github.com/antlr4-go/antlr/v4"
	"github.com/golang-collections/collections/stack"
	"go.uber.org/multierr"
)

type CMCLErrorListener struct {
	*antlr.DefaultErrorListener
	errors []error
}

func (d *CMCLErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{},
	line, column int, msg string, e antlr.RecognitionException) {
	d.errors = append(d.errors, fmt.Errorf("line %d:%d %s", line, column, msg))
}

type CheckParser struct {
	*parser_cmcl.BaseCMCLListener

	stack         stack.Stack
	executionTree *cmclNode
}

func (p *CheckParser) parse(check string) (*cmclNode, error) {
	// Parse check
	input := antlr.NewInputStream(check)
	lexer := parser_cmcl.NewCMCLLexer(input)

	// Add error listener
	errorListener := &CMCLErrorListener{}
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)

	// Check for errors
	if len(errorListener.errors) > 0 {
		return nil, fmt.Errorf("syntax errors: %v", multierr.Combine(errorListener.errors...))
	}

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := parser_cmcl.NewCMCLParser(stream)

	// Add error listener
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errorListener)

	tree := parser.Check()

	// Check for errors
	if len(errorListener.errors) > 0 {
		return nil, fmt.Errorf("syntax errors: %v", multierr.Combine(errorListener.errors...))
	}

	// Walk the tree
	walker := antlr.NewParseTreeWalker()
	walker.Walk(p, tree)

	return p.executionTree, nil
}

// EnterIfCheck is called when production ifCheck is entered.
func (p *CheckParser) EnterIfCheck(ctx *parser_cmcl.IfCheckContext) {
	// Create new node for if statement
	newNode := &cmclNode{
		nodeType:         cmclIfCheck,
		children:         make([]*cmclNode, 0),
		elseIfStatements: make([]*cmclNode, 0),
	}

	// Add node to execution tree
	if p.executionTree == nil { // Root node
		p.executionTree = newNode
	} else { // Child node
		parentNode := p.stack.Peek().(*cmclNode)
		parentNode.children = append(parentNode.children, newNode)
	}

	// Add node the the stack
	p.stack.Push(newNode)
}

// ExitIfCheck is called when production ifCheck is exited.
func (p *CheckParser) ExitIfCheck(ctx *parser_cmcl.IfCheckContext) {
	// Pop node from the stack
	p.stack.Pop()
}

// EnterForeachCheck is called when production foreachCheck is entered.
func (p *CheckParser) EnterForeachCheck(ctx *parser_cmcl.ForeachCheckContext) {
	// Create new node for foreach statement
	newNode := &cmclNode{
		nodeType: cmclForeachCheck,
		children: make([]*cmclNode, 0),
	}

	// Add list item alias to the node as a child
	newNode.children = append(newNode.children, &cmclNode{
		nodeType: cmclForeachItemAlias,
		value:    ctx.Foreach().NAME().GetText(),
	})

	// Add field being iterated over to the node as a child
	newNode.children = append(newNode.children, &cmclNode{
		nodeType: cmclForeachListArg,
		value:    ctx.Foreach().Field().GetText(),
	})

	// Add node to execution tree
	if p.executionTree == nil { // Root node
		p.executionTree = newNode
	} else { // Child node
		parentNode := p.stack.Peek().(*cmclNode)
		parentNode.children = append(parentNode.children, newNode)
	}

	// Add node the the stack
	p.stack.Push(newNode)
}

// ExitForeachCheck is called when production foreachCheck is exited.
func (p *CheckParser) ExitForeachCheck(ctx *parser_cmcl.ForeachCheckContext) {
	// Pop node from the stack
	p.stack.Pop()
}

// EnterOrExpr is called when production orExpr is entered.
func (p *CheckParser) EnterOrExpr(ctx *parser_cmcl.OrExprContext) {
	// Create new node for or expression
	newNode := &cmclNode{
		nodeType: cmclOrExpr,
		children: make([]*cmclNode, 0),
	}

	// Add node to execution tree
	if p.executionTree == nil { // Root node
		p.executionTree = newNode
	} else { // Child node
		parentNode := p.stack.Peek().(*cmclNode)
		parentNode.children = append(parentNode.children, newNode)
	}

	// Add node the the stack
	p.stack.Push(newNode)
}

// ExitOrExpr is called when production orExpr is exited.
func (p *CheckParser) ExitOrExpr(ctx *parser_cmcl.OrExprContext) {
	// Pop node from the stack
	p.stack.Pop()
}

// EnterAndExpr is called when production andExpr is entered.
func (p *CheckParser) EnterAndExpr(ctx *parser_cmcl.AndExprContext) {
	// Create new node for and expression
	newNode := &cmclNode{
		nodeType: cmclAndExpr,
		children: make([]*cmclNode, 0),
	}

	// Add node to execution tree
	if p.executionTree == nil { // Root node
		p.executionTree = newNode
	} else { // Child node
		parentNode := p.stack.Peek().(*cmclNode)
		parentNode.children = append(parentNode.children, newNode)
	}

	// Add node the the stack
	p.stack.Push(newNode)
}

// ExitAndExpr is called when production andExpr is exited.
func (p *CheckParser) ExitAndExpr(ctx *parser_cmcl.AndExprContext) {
	// Pop node from the stack
	p.stack.Pop()
}

// EnterNotExpr is called when production notExpr is entered.
func (p *CheckParser) EnterNotExpr(ctx *parser_cmcl.NotExprContext) {
	// Create new node for not expression
	newNode := &cmclNode{
		nodeType: cmclNotExpr,
		children: make([]*cmclNode, 0),
	}

	// Add node to execution tree
	if p.executionTree == nil { // Root node
		p.executionTree = newNode
	} else { // Child node
		parentNode := p.stack.Peek().(*cmclNode)
		parentNode.children = append(parentNode.children, newNode)
	}

	// Add node the the stack
	p.stack.Push(newNode)
}

// ExitNotExpr is called when production notExpr is exited.
func (p *CheckParser) ExitNotExpr(ctx *parser_cmcl.NotExprContext) {
	// Pop node from the stack
	p.stack.Pop()
}

// EnterFieldCheck is called when production fieldCheck is entered.
func (p *CheckParser) EnterFieldExpr(ctx *parser_cmcl.FieldExprContext) {
	// Create new node for field expression
	newNode := &cmclNode{
		nodeType: cmclFieldExpr,
		value:    ctx.FieldExpression().Field().GetText(),
		children: make([]*cmclNode, 0),
	}

	// Add node to execution tree
	if p.executionTree == nil { // Root node
		p.executionTree = newNode
	} else { // Child node
		parentNode := p.stack.Peek().(*cmclNode)
		parentNode.children = append(parentNode.children, newNode)
	}

	// Add node the the stack
	p.stack.Push(newNode)
}

// ExitFieldCheck is called when production fieldCheck is exited.
func (p *CheckParser) ExitFieldExpr(ctx *parser_cmcl.FieldExprContext) {
	// Pop node from the stack
	p.stack.Pop()
}

// EnterParenExpr is called when production parenExpr is entered.
func (p *CheckParser) EnterParenExpr(ctx *parser_cmcl.ParenExprContext) {
	// Create new node for paren expression
	newNode := &cmclNode{
		nodeType: cmclParenExpr,
		children: make([]*cmclNode, 0),
	}

	// Add node to execution tree
	if p.executionTree == nil { // Root node
		p.executionTree = newNode
	} else { // Child node
		parentNode := p.stack.Peek().(*cmclNode)
		parentNode.children = append(parentNode.children, newNode)
	}

	// Add node the the stack
	p.stack.Push(newNode)
}

// ExitParenExpr is called when production parenExpr is exited.
func (p *CheckParser) ExitParenExpr(ctx *parser_cmcl.ParenExprContext) {
	// Pop node from the stack
	p.stack.Pop()
}

// EnterElseif is called when production elseif is entered.
func (p *CheckParser) EnterElseif(ctx *parser_cmcl.ElseifContext) {
	// Create new node for else if statement
	newNode := &cmclNode{
		nodeType: cmclIfCheck,
		children: make([]*cmclNode, 0),
	}

	// Add node to the execution tree
	parentNode := p.stack.Peek().(*cmclNode)
	parentNode.elseIfStatements = append(parentNode.elseIfStatements, newNode)

	// Add node the the stack
	p.stack.Push(newNode)
}

// ExitElseif is called when production elseif is exited.
func (p *CheckParser) ExitElseif(ctx *parser_cmcl.ElseifContext) {
	// Pop node from the stack
	p.stack.Pop()
}

// EnterElse is called when production else is entered.
func (p *CheckParser) EnterElse(ctx *parser_cmcl.ElseContext) {
	// Create new node for else statement
	newNode := &cmclNode{
		nodeType: cmclIfCheck,
		children: make([]*cmclNode, 0),
	}

	// Add node to the execution tree
	parentNode := p.stack.Peek().(*cmclNode)
	parentNode.elseStatement = newNode

	// Add node the the stack
	p.stack.Push(newNode)
}

// ExitElse is called when production else is exited.
func (p *CheckParser) ExitElse(ctx *parser_cmcl.ElseContext) {
	// Pop node from the stack
	p.stack.Pop()
}

// EnterFunction is called when production function is entered.
func (p *CheckParser) EnterFuncExpr(ctx *parser_cmcl.FuncExprContext) {
	// Create new node for function
	newNode := &cmclNode{
		nodeType: cmclFuncExpr,
		children: make([]*cmclNode, 0),
	}

	// Add node to execution tree
	if p.executionTree == nil { // Root node
		p.executionTree = newNode
	} else { // Child node
		parentNode := p.stack.Peek().(*cmclNode)
		parentNode.children = append(parentNode.children, newNode)
	}

	// Add node the the stack
	p.stack.Push(newNode)
}

func (p *CheckParser) ExitFuncExpr(ctx *parser_cmcl.FuncExprContext) {
	// Pop node from the stack
	p.stack.Pop()
}

// EnterFunction is called when production function is entered.
func (p *CheckParser) EnterFunction(ctx *parser_cmcl.FunctionContext) {
	// Create new node for function
	newNode := &cmclNode{
		nodeType: cmclFunction,
		value:    ctx.NAME().GetText(),
		children: make([]*cmclNode, 0),
	}

	// Add node to execution tree
	if p.executionTree == nil { // Root node
		p.executionTree = newNode
	} else { // Child node
		parentNode := p.stack.Peek().(*cmclNode)
		parentNode.children = append(parentNode.children, newNode)
	}

	// Add node the the stack
	p.stack.Push(newNode)
}

// ExitFunction is called when production function is exited.
func (p *CheckParser) ExitFunction(ctx *parser_cmcl.FunctionContext) {
	// Pop node from the stack
	p.stack.Pop()
}

// EnterString is called when production string is entered.
func (p *CheckParser) EnterString(ctx *parser_cmcl.StringContext) {
	// Create new node for string
	newNode := &cmclNode{
		nodeType: cmclString,
		value:    ctx.GetText(),
	}

	// Add node to the execution tree
	parentNode := p.stack.Peek().(*cmclNode)
	parentNode.children = append(parentNode.children, newNode)
}

// EnterInt is called when production int is entered.
func (p *CheckParser) EnterInt(ctx *parser_cmcl.IntContext) {
	// Create new node for int
	newNode := &cmclNode{
		nodeType: cmclInt,
		value:    ctx.GetText(),
	}

	// Add node to the execution tree
	parentNode := p.stack.Peek().(*cmclNode)
	parentNode.children = append(parentNode.children, newNode)
}

// EnterFloat is called when production float is entered.
func (p *CheckParser) EnterFloat(ctx *parser_cmcl.FloatContext) {
	// Create new node for float
	newNode := &cmclNode{
		nodeType: cmclFloat,
		value:    ctx.GetText(),
	}

	// Add node to the execution tree
	parentNode := p.stack.Peek().(*cmclNode)
	parentNode.children = append(parentNode.children, newNode)
}

// EnterBoolean is called when production boolean is entered.
func (p *CheckParser) EnterBoolean(ctx *parser_cmcl.BooleanContext) {
	// Create new node for boolean
	newNode := &cmclNode{
		nodeType: cmclBool,
		value:    ctx.GetText(),
	}

	// Add node to the execution tree
	parentNode := p.stack.Peek().(*cmclNode)
	parentNode.children = append(parentNode.children, newNode)
}
