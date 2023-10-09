package parsers

import (
	"strconv"
	"strings"

	parser_json "github.com/ConfigMate/configmate/parsers/gen/parser_JSON/parsers/grammars"
	"github.com/antlr4-go/antlr/v4"
	"github.com/golang-collections/collections/stack"
)

type JsonParser struct{}

func (p *JsonParser) Parse(data []byte) (*Node, error) {
	input := antlr.NewInputStream(string(data))
	lexer := parser_json.NewJSONLexer(input)
	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := parser_json.NewJSONParser(tokenStream)
	tree := parser.Json()

	walker := antlr.NewParseTreeWalker()

	jsonListener := &JsonParserListener{}
	walker.Walk(jsonListener, tree)

	return jsonListener.configFile, nil
}

type JsonParserListener struct {
	*parser_json.BaseJSONListener

	configFile *Node
	stack      stack.Stack
}

func (l *JsonParserListener) getTOS() *Node {
	return l.stack.Peek().(*Node)
}

func (l *JsonParserListener) getTOSObject() map[string]*Node {
	return l.getTOS().Value.(map[string]*Node)
}

func (l *JsonParserListener) EnterJson(ctx *parser_json.JsonContext) {
	l.configFile = &Node{Type: Object, Value: map[string]*Node{}}
	l.stack.Push(l.configFile)
}

func (l *JsonParserListener) EnterObject(ctx *parser_json.ObjContext) {
	l.getTOS().Type = Object
	l.getTOS().Value = map[string]*Node{}
}

func (l *JsonParserListener) ExitObject(ctx *parser_json.ObjContext) {
	l.stack.Pop()
}

func (l *JsonParserListener) EnterPair(ctx *parser_json.PairContext) {
	pairNode := &Node{Type: Null, Value: nil}
	l.getTOSObject()[ctx.STRING().GetText()] = pairNode
	l.stack.Push(pairNode)
}

func (l *JsonParserListener) ExitPair(ctx *parser_json.PairContext) {
	l.stack.Pop()
}

func (l *JsonParserListener) EnterArray(ctx *parser_json.ArrayContext) {
	l.getTOS().Type = Array
	l.getTOS().Value = make([]*Node, 0)
}

func (l *JsonParserListener) ExitArray(ctx *parser_json.ArrayContext) {
	l.stack.Pop()
}

func (l *JsonParserListener) EnterNumber(ctx *parser_json.NumberContext) {
	if strings.Contains(ctx.NUMBER().GetText(), ".") {
		l.getTOS().Type = Float
		value, err := strconv.ParseFloat(ctx.NUMBER().GetText(), 64)
		if err != nil {
			panic(err)
		}
		l.getTOS().Value = value
	} else {
		l.getTOS().Type = Int
		value, err := strconv.Atoi(ctx.NUMBER().GetText())
		if err != nil {
			panic(err)
		}
		l.getTOS().Value = value
	}
}

func (l *JsonParserListener) EnterString(ctx *parser_json.StringContext) {
	l.getTOS().Type = String
	l.getTOS().Value = ctx.STRING().GetText()
}

func (l *JsonParserListener) EnterTrue(ctx *parser_json.BooleanTrueContext) {
	l.getTOS().Type = Bool
	l.getTOS().Value = true
}

func (l *JsonParserListener) EnterFalse(ctx *parser_json.BooleanFalseContext) {
	l.getTOS().Type = Bool
	l.getTOS().Value = false
}

func (l *JsonParserListener) EnterNull(ctx *parser_json.NullContext) {
	l.getTOS().Type = Null
	l.getTOS().Value = nil
}
