package parsers

import (
	"strconv"
	"strings"

	"github.com/ConfigMate/configmate/parsers/gen/parser_json"
	"github.com/antlr4-go/antlr/v4"
	"github.com/golang-collections/collections/stack"
)

type jsonParser struct {
	*parser_json.BaseJSONListener

	configFile *Node
	stack      stack.Stack
}

// Custom Json Parser
func (p *jsonParser) Parse(data []byte) (*Node, []CMParserError) {
	// Initialize the error listener
	errorListener := &CMErrorListener{}

	input := antlr.NewInputStream(string(data))
	lexer := parser_json.NewJSONLexer(input)

	// Attach the error listener to the lexer
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)

	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := parser_json.NewJSONParser(tokenStream)

	// Attach the error listener to the parser
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errorListener)
	tree := parser.Json()

	// Check for errors after parsing
	if len(errorListener.errors) > 0 {
		return nil, errorListener.errors
	}

	// Initialize config file
	p.configFile = nil

	// Initialize stack
	p.stack = stack.Stack{}

	walker := antlr.NewParseTreeWalker()
	walker.Walk(p, tree)

	return p.configFile, nil
}

func (l *jsonParser) EnterObj(ctx *parser_json.ObjContext) {
	if l.configFile == nil { // This object is the top level entity
		// Create new node for object
		node := &Node{Type: Object, Value: map[string]*Node{}}

		// Set config file to point to this object
		l.configFile = node

		// Push object node to stack
		l.stack.Push(node)

	} else if l.getTOS().Type == Null { // Is the value of a pair
		// Set pair node to correct type
		l.getTOS().Type = Object

		// Create value for object
		l.getTOS().Value = map[string]*Node{}

		// Push object node to stack. This is redundant, since
		// the object node is already on the stack, but the stack will
		// be popped in ExitObj, so we need to push it again.
		l.stack.Push(l.getTOS())

	} else if l.getTOS().Type == Array { // Is an element of an array
		// Create new node for object
		node := &Node{Type: Object, Value: map[string]*Node{}}

		// Add object node to array
		l.getTOS().Value = append(l.getTOSArray(), node)

		// Push object node to stack
		l.stack.Push(node)

	} else {
		panic("Invalid state")
	}

	// Add start location information of the object
	l.getTOS().ValueLocation.Start.Line = ctx.GetStart().GetLine() - 1
	l.getTOS().ValueLocation.Start.Column = ctx.GetStart().GetColumn()

	// Add end location information of the object
	l.getTOS().ValueLocation.End.Line = ctx.GetStop().GetLine() - 1
	l.getTOS().ValueLocation.End.Column = ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText())
}

func (l *jsonParser) ExitObj(ctx *parser_json.ObjContext) {
	l.stack.Pop()
}

func (l *jsonParser) EnterPair(ctx *parser_json.PairContext) {
	// Create new node for pair. The type will be set in Enter<Value>
	node := &Node{Type: Null, Value: nil}

	// Get pair key without quotes
	key := removeQuotes(ctx.STRING().GetText())

	// Add pair node to parent object node
	l.getTOSObject()[key] = node

	// Push pair node to stack
	l.stack.Push(node)

	// Add start location information of the pair key
	node.NameLocation.Start.Line = ctx.GetStart().GetLine() - 1
	node.NameLocation.Start.Column = ctx.GetStart().GetColumn()

	// Add end location information of the pair key
	node.NameLocation.End.Line = ctx.GetStart().GetLine() - 1
	node.NameLocation.End.Column = ctx.GetStart().GetColumn() + len(ctx.GetStart().GetText())
}

func (l *jsonParser) ExitPair(ctx *parser_json.PairContext) {
	l.stack.Pop()
}

func (l *jsonParser) EnterArr(ctx *parser_json.ArrContext) {
	if l.configFile == nil { // This array is the top level entity
		// Create new node for array
		node := &Node{Type: Array, Value: []*Node{}}

		// Set config file to point to this array
		l.configFile = node

		// Push array node to stack
		l.stack.Push(node)

	} else if l.getTOS().Type == Null { // Is the value of a pair
		// Set pair node to correct type
		l.getTOS().Type = Array

		// Create value for array
		l.getTOS().Value = []*Node{}

		// Push array node to stack. This is redundant, since
		// the array node is already on the stack, but the stack will
		// be popped in ExitArr, so we need to push it again.
		l.stack.Push(l.getTOS())

	} else if l.getTOS().Type == Array { // Is an element of an array
		// Create new node for array
		node := &Node{Type: Array, Value: []*Node{}}

		// Add array node to array
		l.getTOS().Value = append(l.getTOSArray(), node)

		// Push array node to stack
		l.stack.Push(node)

	} else {
		panic("Invalid state")
	}

	// Add start location information of the array
	l.getTOS().ValueLocation.Start.Line = ctx.GetStart().GetLine() - 1
	l.getTOS().ValueLocation.Start.Column = ctx.GetStart().GetColumn()

	// Add end location information of the array
	l.getTOS().ValueLocation.End.Line = ctx.GetStop().GetLine() - 1
	l.getTOS().ValueLocation.End.Column = ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText())
}

func (l *jsonParser) ExitArr(ctx *parser_json.ArrContext) {
	l.stack.Pop()
}

func (l *jsonParser) EnterNumber(ctx *parser_json.NumberContext) {
	// Get value and type
	var numberType FieldType
	var value interface{}
	var err error
	if strings.Contains(ctx.NUMBER().GetText(), ".") {
		numberType = Float
		value, err = strconv.ParseFloat(ctx.NUMBER().GetText(), 64)
		if err != nil {
			panic(err)
		}
	} else {
		numberType = Int
		value, err = strconv.Atoi(ctx.NUMBER().GetText())
		if err != nil {
			panic(err)
		}
	}

	// Create holder to store the node where the location information
	// of the number value should be stored
	var locationInfoDest *Node
	if l.configFile == nil { // This number is the top level entity
		// Create new node for number
		node := &Node{Type: numberType, Value: value}

		// Set config file to point to this number
		l.configFile = node

		// Set location destination to the newly created node
		locationInfoDest = node

	} else if l.getTOS().Type == Null { // Is the value of a pair
		// Set pair node to correct type
		l.getTOS().Type = numberType

		// Create value for number
		l.getTOS().Value = value

		// Set location destination to the pair node
		locationInfoDest = l.getTOS()

	} else if l.getTOS().Type == Array { // Is an element of an array
		// Create new node for number
		node := &Node{Type: numberType, Value: value}

		// Add number node to array
		l.getTOS().Value = append(l.getTOSArray(), node)

		// Set location destination to the newly created node
		locationInfoDest = node

	} else {
		panic("Invalid state")
	}

	// Add start location information of the string value
	locationInfoDest.ValueLocation.Start.Line = ctx.GetStart().GetLine() - 1
	locationInfoDest.ValueLocation.Start.Column = ctx.GetStart().GetColumn()

	// Add end location information of the string value
	locationInfoDest.ValueLocation.End.Line = ctx.GetStop().GetLine() - 1
	locationInfoDest.ValueLocation.End.Column = ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText())
}

func (l *jsonParser) EnterString(ctx *parser_json.StringContext) {
	// Get value
	value := removeQuotes(ctx.STRING().GetText())

	// Create holder to store the node where the location information
	// of the string value should be stored
	var locationInfoDest *Node
	if l.configFile == nil { // This string is the top level entity
		// Create new node for string
		node := &Node{Type: String, Value: value}

		// Set config file to point to this string
		l.configFile = node

		// Set location destination to the newly created node
		locationInfoDest = node

	} else if l.getTOS().Type == Null { // Is the value of a pair
		// Set pair node to correct type
		l.getTOS().Type = String

		// Set value for string
		l.getTOS().Value = value

		// Set location destination to the pair node
		locationInfoDest = l.getTOS()

	} else if l.getTOS().Type == Array { // Is an element of an array
		// Create new node for string
		node := &Node{Type: String, Value: value}

		// Add string node to array
		l.getTOS().Value = append(l.getTOSArray(), node)

		// Set location destination to the newly created node
		locationInfoDest = node

	} else {
		panic("Invalid state")
	}

	// Add start location information of the string value
	locationInfoDest.ValueLocation.Start.Line = ctx.GetStart().GetLine() - 1
	locationInfoDest.ValueLocation.Start.Column = ctx.GetStart().GetColumn()

	// Add end location information of the string value
	locationInfoDest.ValueLocation.End.Line = ctx.GetStop().GetLine() - 1
	locationInfoDest.ValueLocation.End.Column = ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText())
}

func (l *jsonParser) EnterBooleanTrue(ctx *parser_json.BooleanTrueContext) {

	// Create holder to store the node where the location information
	// of the boolean value should be stored
	var locationInfoDest *Node
	if l.configFile == nil { // This boolean is the top level entity
		// Create new node for boolean
		node := &Node{Type: Bool, Value: true}

		// Set config file to point to this boolean
		l.configFile = node

		// Set location destination to the newly created node
		locationInfoDest = node

	} else if l.getTOS().Type == Null { // Is the value of a pair
		// Set pair node to correct type
		l.getTOS().Type = Bool

		// Set value for boolean
		l.getTOS().Value = true

		// Set location destination to the pair node
		locationInfoDest = l.getTOS()

	} else if l.getTOS().Type == Array { // Is an element of an array
		// Create new node for boolean
		node := &Node{Type: Bool, Value: true}

		// Add boolean node to array
		l.getTOS().Value = append(l.getTOSArray(), node)

		// Set location destination to the newly created node
		locationInfoDest = node

	} else {
		panic("Invalid state")
	}

	// Add start location information of the string value
	locationInfoDest.ValueLocation.Start.Line = ctx.GetStart().GetLine() - 1
	locationInfoDest.ValueLocation.Start.Column = ctx.GetStart().GetColumn()

	// Add end location information of the string value
	locationInfoDest.ValueLocation.End.Line = ctx.GetStop().GetLine() - 1
	locationInfoDest.ValueLocation.End.Column = ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText())
}

func (l *jsonParser) EnterBooleanFalse(ctx *parser_json.BooleanFalseContext) {

	// Create holder to store the node where the location information
	// of the boolean value should be stored
	var locationInfoDest *Node
	if l.configFile == nil { // This boolean is the top level entity
		// Create new node for boolean
		node := &Node{Type: Bool, Value: false}

		// Set config file to point to this boolean
		l.configFile = node

		// Set location destination to the newly created node
		locationInfoDest = node

	} else if l.getTOS().Type == Null { // Is the value of a pair
		// Set pair node to correct type
		l.getTOS().Type = Bool

		// Set value for boolean
		l.getTOS().Value = false

		// Set location destination to the pair node
		locationInfoDest = l.getTOS()

	} else if l.getTOS().Type == Array { // Is an element of an array
		// Create new node for boolean
		node := &Node{Type: Bool, Value: false}

		// Add boolean node to array
		l.getTOS().Value = append(l.getTOSArray(), node)

		// Set location destination to the newly created node
		locationInfoDest = node

	} else {
		panic("Invalid state")
	}

	// Add start location information of the string value
	locationInfoDest.ValueLocation.Start.Line = ctx.GetStart().GetLine() - 1
	locationInfoDest.ValueLocation.Start.Column = ctx.GetStart().GetColumn()

	// Add end location information of the string value
	locationInfoDest.ValueLocation.End.Line = ctx.GetStop().GetLine() - 1
	locationInfoDest.ValueLocation.End.Column = ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText())
}

func (l *jsonParser) EnterNull(ctx *parser_json.NullContext) {

	// Create holder to store the node where the location information
	// of the null value should be stored
	var locationInfoDest *Node
	if l.configFile == nil { // This null is the top level entity
		// Create new node for null
		node := &Node{Type: Null, Value: nil}

		// Set config file to point to this null
		l.configFile = node

		// Set location destination to the newly created node
		locationInfoDest = node

	} else if l.getTOS().Type == Null { // Is the value of a pair
		// Set pair node to correct type
		l.getTOS().Type = Null

		// Set value for null
		l.getTOS().Value = nil

		// Set location destination to the pair node
		locationInfoDest = l.getTOS()

	} else if l.getTOS().Type == Array { // Is an element of an array
		// Create new node for null
		node := &Node{Type: Null, Value: nil}

		// Add null node to array
		l.getTOS().Value = append(l.getTOSArray(), node)

		// Set location destination to the newly created node
		locationInfoDest = node

	} else {
		panic("Invalid state")
	}

	// Add start location information of the string value
	locationInfoDest.ValueLocation.Start.Line = ctx.GetStart().GetLine() - 1
	locationInfoDest.ValueLocation.Start.Column = ctx.GetStart().GetColumn()

	// Add end location information of the string value
	locationInfoDest.ValueLocation.End.Line = ctx.GetStop().GetLine() - 1
	locationInfoDest.ValueLocation.End.Column = ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText())
}

func (l *jsonParser) getTOS() *Node {
	return l.stack.Peek().(*Node)
}

func (l *jsonParser) getTOSObject() map[string]*Node {
	return l.getTOS().Value.(map[string]*Node)
}

func (l *jsonParser) getTOSArray() []*Node {
	return l.getTOS().Value.([]*Node)
}

func removeQuotes(s string) string {
	return s[1 : len(s)-1]
}
