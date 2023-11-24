package parsers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ConfigMate/configmate/parsers/gen/parser_toml"
	"github.com/antlr4-go/antlr/v4"
	"github.com/golang-collections/collections/stack"
)

type tomlParser struct {
	*parser_toml.BaseTOMLParserListener

	configFile            *Node
	stack                 stack.Stack
	directlyDefinedTables map[string]bool
	definedTables         map[string]tomlTableDefinitionMode
	errs                  []CMParserError
}

type tomlTableDefinitionMode int

const (
	tomlDefModeDirect tomlTableDefinitionMode = iota
	tomlDefModeTable
	tomlDefModeInlineTable
)

// Custom TOML parser
func (p *tomlParser) Parse(data []byte) (*Node, []CMParserError) {
	// Initialize the error listener
	errorListener := &CMErrorListener{}

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
		return nil, errorListener.errors
	}

	// Initialize config file to an object
	p.configFile = &Node{
		Type:  Object,
		Value: map[string]*Node{},
	}

	// Initialize maps
	p.directlyDefinedTables = make(map[string]bool)
	p.definedTables = make(map[string]tomlTableDefinitionMode)

	// Initialize Stack
	p.stack = stack.Stack{}
	p.stack.Push(p.configFile)

	walker := antlr.NewParseTreeWalker()
	walker.Walk(p, tree)

	// Check for errors
	if len(p.errs) > 0 {
		return nil, p.errs
	}

	return p.configFile, nil
}

// EnterKey_value is called when production key_value is entered.
func (p *tomlParser) EnterKey_value(ctx *parser_toml.Key_valueContext) {
	// Parse key
	fieldKey := p.parseKey(ctx.Key())

	// Get parent node in stack
	parentNode := p.stack.Peek().(*Node)

	// Get or create parent key node
	fieldNode, err := p.getOrCreateNode(parentNode, fieldKey, tomlDefModeDirect)
	if err != nil {
		p.errs = append(p.errs, CMParserError{
			Message: err.Error(),
			Location: TokenLocation{
				Start: CharLocation{
					Line:   ctx.Key().GetStart().GetLine() - 1,
					Column: ctx.Key().GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.Key().GetStop().GetLine() - 1,
					Column: ctx.Key().GetStop().GetColumn() + len(ctx.Key().GetStop().GetText()),
				},
			},
		})
		return
	}

	// Check if this key was already defined
	if fieldNode.Type != Null {
		p.errs = append(p.errs, CMParserError{
			Message: fmt.Errorf("can't redefine existing key: '%s'", fieldKey.String()).Error(),
			Location: TokenLocation{
				Start: CharLocation{
					Line:   ctx.Key().GetStart().GetLine() - 1,
					Column: ctx.Key().GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.Key().GetStop().GetLine() - 1,
					Column: ctx.Key().GetStop().GetColumn() + len(ctx.Key().GetStop().GetText()),
				},
			},
		})
		return
	}

	// Set name location
	fieldNode.NameLocation = TokenLocation{
		Start: CharLocation{
			Line:   ctx.Key().GetStart().GetLine() - 1,
			Column: ctx.Key().GetStart().GetColumn(),
		},
		End: CharLocation{
			Line:   ctx.Key().GetStop().GetLine() - 1,
			Column: ctx.Key().GetStop().GetColumn() + len(ctx.Key().GetStop().GetText()),
		},
	}

	// Add fieldnode to stack
	p.stack.Push(fieldNode)
}

// ExitKey_value is called when production key_value is exited.
func (p *tomlParser) ExitKey_value(ctx *parser_toml.Key_valueContext) {
	// Pop field node from stack
	p.stack.Pop()
}

// EnterStandard_table is called when production standard_table is entered.
func (p *tomlParser) EnterStandard_table(ctx *parser_toml.Standard_tableContext) {
	// Parse key
	fieldKey := p.parseKey(ctx.Key())

	// Get or create parent key node from root
	fieldNode, err := p.getOrCreateNode(p.configFile, fieldKey, tomlDefModeTable)
	if err != nil {
		p.errs = append(p.errs, CMParserError{
			Message: err.Error(),
			Location: TokenLocation{
				Start: CharLocation{
					Line:   ctx.Key().GetStart().GetLine() - 1,
					Column: ctx.Key().GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.Key().GetStop().GetLine() - 1,
					Column: ctx.Key().GetStop().GetColumn() + len(ctx.Key().GetStop().GetText()),
				},
			},
		})
		return
	}

	// Check if table was already directly defined
	if p.directlyDefinedTables[fmt.Sprintf("%p", fieldNode)] {
		p.errs = append(p.errs, CMParserError{
			Message: fmt.Errorf("can't redefine existing table: '%s'", fieldKey.String()).Error(),
			Location: TokenLocation{
				Start: CharLocation{
					Line:   ctx.Key().GetStart().GetLine() - 1,
					Column: ctx.Key().GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.Key().GetStop().GetLine() - 1,
					Column: ctx.Key().GetStop().GetColumn() + len(ctx.Key().GetStop().GetText()),
				},
			},
		})
		return
	}

	// Set table as directly defined
	p.directlyDefinedTables[fmt.Sprintf("%p", fieldNode)] = true

	// Check if table is new
	if fieldNode.Type == Null {
		// Set as object node
		fieldNode.Type = Object
		fieldNode.Value = map[string]*Node{}
	} else if fieldNode.Type != Object {
		p.errs = append(p.errs, CMParserError{
			Message: fmt.Errorf("can't redefine existing key: '%s'", fieldKey.String()).Error(),
			Location: TokenLocation{
				Start: CharLocation{
					Line:   ctx.Key().GetStart().GetLine() - 1,
					Column: ctx.Key().GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.Key().GetStop().GetLine() - 1,
					Column: ctx.Key().GetStop().GetColumn() + len(ctx.Key().GetStop().GetText()),
				},
			},
		})
		return
	}

	// Add location info
	fieldNode.NameLocation = TokenLocation{
		Start: CharLocation{
			Line:   ctx.GetStart().GetLine() - 1,
			Column: ctx.GetStart().GetColumn(),
		},
		End: CharLocation{
			Line:   ctx.GetStop().GetLine() - 1,
			Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
		},
	}

	// We cannot find value location, using name location to
	// guarantee better display result in case this is used
	fieldNode.ValueLocation = fieldNode.NameLocation

	// Add table node to stack
	p.stack.Push(fieldNode)
}

// EnterArray_table is called when production array_table is entered.
func (p *tomlParser) EnterArray_table(ctx *parser_toml.Array_tableContext) {
	// Parse key
	fieldKey := p.parseKey(ctx.Key())

	// Get or create parent key node from root
	arrayNode, err := p.getOrCreateNode(p.configFile, fieldKey, tomlDefModeTable)
	if err != nil {
		p.errs = append(p.errs, CMParserError{
			Message: err.Error(),
			Location: TokenLocation{
				Start: CharLocation{
					Line:   ctx.Key().GetStart().GetLine() - 1,
					Column: ctx.Key().GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.Key().GetStop().GetLine() - 1,
					Column: ctx.Key().GetStop().GetColumn() + len(ctx.Key().GetStop().GetText()),
				},
			},
		})
		return
	}

	// Check if array is new
	if arrayNode.Type == Null {
		// Set as array node
		arrayNode.Type = Array
		arrayNode.Value = []*Node{}
	} else if arrayNode.Type != Array {
		p.errs = append(p.errs, CMParserError{
			Message: fmt.Errorf("can't redefine existing key: '%s'", fieldKey.String()).Error(),
			Location: TokenLocation{
				Start: CharLocation{
					Line:   ctx.Key().GetStart().GetLine() - 1,
					Column: ctx.Key().GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.Key().GetStop().GetLine() - 1,
					Column: ctx.Key().GetStop().GetColumn() + len(ctx.Key().GetStop().GetText()),
				},
			},
		})
		return
	}

	// Create table node
	newInArrayTable := &Node{
		Type:  Object,
		Value: map[string]*Node{},
		NameLocation: TokenLocation{
			Start: CharLocation{
				Line:   ctx.GetStart().GetLine() - 1,
				Column: ctx.GetStart().GetColumn(),
			},
			End: CharLocation{
				Line:   ctx.GetStop().GetLine() - 1,
				Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
			},
		},
	}

	// We cannot find value location, using name location to
	// guarantee better display result in case this is used
	newInArrayTable.ValueLocation = newInArrayTable.NameLocation

	// Add new node to array
	arrayNode.Value = append(arrayNode.Value.([]*Node), newInArrayTable)

	// Add table node to stack
	p.stack.Push(newInArrayTable)
}

// EnterInline_table is called when production inline_table is entered.
func (p *tomlParser) EnterInline_table(ctx *parser_toml.Inline_tableContext) {
	// Get parent node in stack
	parentNode := p.stack.Peek().(*Node)

	// If parent node is an array, append the inline table to the array
	if parentNode.Type == Array {
		parentNode.Value = append(parentNode.Value.([]*Node), &Node{
			Type:  Object,
			Value: map[string]*Node{},
			ValueLocation: TokenLocation{
				Start: CharLocation{
					Line:   ctx.GetStart().GetLine() - 1,
					Column: ctx.GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.GetStop().GetLine() - 1,
					Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
				},
			},
		})

		// Add inline table node to stack
		p.stack.Push(parentNode.Value.([]*Node)[len(parentNode.Value.([]*Node))-1])
	} else { // Set parent node as inline table (node created when key was found)
		parentNode.Type = Object
		parentNode.Value = map[string]*Node{}
		parentNode.ValueLocation = TokenLocation{
			Start: CharLocation{
				Line:   ctx.GetStart().GetLine() - 1,
				Column: ctx.GetStart().GetColumn(),
			},
			End: CharLocation{
				Line:   ctx.GetStop().GetLine() - 1,
				Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
			},
		}

		// Set as inline table definition
		p.definedTables[fmt.Sprintf("%p", parentNode)] = tomlDefModeInlineTable

		// Push again (redundant) to keep stack consistent
		p.stack.Push(parentNode)
	}
}

// ExitInline_table is called when production inline_table is exited.
func (p *tomlParser) ExitInline_table(ctx *parser_toml.Inline_tableContext) {
	// Pop inline table node from stack
	p.stack.Pop()
}

// EnterArray is called when production array is entered.
func (p *tomlParser) EnterArray(ctx *parser_toml.ArrayContext) {
	// Get parent node in stack
	parentNode := p.stack.Peek().(*Node)

	// If parent node is an array, append
	if parentNode.Type == Array {
		parentNode.Value = append(parentNode.Value.([]*Node), &Node{
			Type:  Array,
			Value: []*Node{},
			ValueLocation: TokenLocation{
				Start: CharLocation{
					Line:   ctx.GetStart().GetLine() - 1,
					Column: ctx.GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.GetStop().GetLine() - 1,
					Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
				},
			},
		})

		// Add array node to stack
		p.stack.Push(parentNode.Value.([]*Node)[len(parentNode.Value.([]*Node))-1])
	} else { // Set parent node as array (node created when key was found)
		parentNode.Type = Array
		parentNode.Value = []*Node{}
		parentNode.ValueLocation = TokenLocation{
			Start: CharLocation{
				Line:   ctx.GetStart().GetLine() - 1,
				Column: ctx.GetStart().GetColumn(),
			},
			End: CharLocation{
				Line:   ctx.GetStop().GetLine() - 1,
				Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
			},
		}

		// Push again (redundant) to keep stack consistent
		p.stack.Push(parentNode)
	}
}

// ExitArray is called when production array is exited.
func (p *tomlParser) ExitArray(ctx *parser_toml.ArrayContext) {
	// Pop array node from stack
	p.stack.Pop()
}

// EnterString is called when production string is entered.
func (p *tomlParser) EnterString(ctx *parser_toml.StringContext) {
	// Get parent node in stack
	parentNode := p.stack.Peek().(*Node)

	// If parent node is an array, append the string to the array
	if parentNode.Type == Array {
		parentNode.Value = append(parentNode.Value.([]*Node), &Node{
			Type:  String,
			Value: p.cleanString(ctx.GetText()),
			ValueLocation: TokenLocation{
				Start: CharLocation{
					Line:   ctx.GetStart().GetLine() - 1,
					Column: ctx.GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.GetStop().GetLine() - 1,
					Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
				},
			},
		})
	} else { // Set parent node as string (node created when key was found)
		parentNode.Type = String
		parentNode.Value = p.cleanString(ctx.GetText())
		parentNode.ValueLocation = TokenLocation{
			Start: CharLocation{
				Line:   ctx.GetStart().GetLine() - 1,
				Column: ctx.GetStart().GetColumn(),
			},
			End: CharLocation{
				Line:   ctx.GetStop().GetLine() - 1,
				Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
			},
		}
	}
}

// EnterInteger is called when production integer is entered.
func (p *tomlParser) EnterInteger(ctx *parser_toml.IntegerContext) {
	// Get parent node in stack
	parentNode := p.stack.Peek().(*Node)

	intValue, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		panic(err)
	}

	// If parent node is an array, append the integer to the array
	if parentNode.Type == Array {
		parentNode.Value = append(parentNode.Value.([]*Node), &Node{
			Type:  Int,
			Value: intValue,
			ValueLocation: TokenLocation{
				Start: CharLocation{
					Line:   ctx.GetStart().GetLine() - 1,
					Column: ctx.GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.GetStop().GetLine() - 1,
					Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
				},
			},
		})
	} else { // Set parent node as integer (node created when key was found)
		parentNode.Type = Int
		parentNode.Value = intValue
		parentNode.ValueLocation = TokenLocation{
			Start: CharLocation{
				Line:   ctx.GetStart().GetLine() - 1,
				Column: ctx.GetStart().GetColumn(),
			},
			End: CharLocation{
				Line:   ctx.GetStop().GetLine() - 1,
				Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
			},
		}
	}
}

// EnterFloating_point is called when production floating_point is entered.
func (p *tomlParser) EnterFloating_point(ctx *parser_toml.Floating_pointContext) {
	// Get parent node in stack
	parentNode := p.stack.Peek().(*Node)

	floatValue, err := strconv.ParseFloat(ctx.GetText(), 64)
	if err != nil {
		panic(err)
	}

	// If parent node is an array, append the float to the array
	if parentNode.Type == Array {
		parentNode.Value = append(parentNode.Value.([]*Node), &Node{
			Type:  Float,
			Value: floatValue,
			ValueLocation: TokenLocation{
				Start: CharLocation{
					Line:   ctx.GetStart().GetLine() - 1,
					Column: ctx.GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.GetStop().GetLine() - 1,
					Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
				},
			},
		})
	} else { // Set parent node as float (node created when key was found)
		parentNode.Type = Float
		parentNode.Value = floatValue
		parentNode.ValueLocation = TokenLocation{
			Start: CharLocation{
				Line:   ctx.GetStart().GetLine() - 1,
				Column: ctx.GetStart().GetColumn(),
			},
			End: CharLocation{
				Line:   ctx.GetStop().GetLine() - 1,
				Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
			},
		}
	}
}

// EnterBool is called when production bool is entered.
func (p *tomlParser) EnterBool(ctx *parser_toml.BoolContext) {
	// Get parent node in stack
	parentNode := p.stack.Peek().(*Node)

	boolValue, err := strconv.ParseBool(ctx.GetText())
	if err != nil {
		panic(err)
	}

	// If parent node is an array, append the bool to the array
	if parentNode.Type == Array {
		parentNode.Value = append(parentNode.Value.([]*Node), &Node{
			Type:  Bool,
			Value: boolValue,
			ValueLocation: TokenLocation{
				Start: CharLocation{
					Line:   ctx.GetStart().GetLine() - 1,
					Column: ctx.GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.GetStop().GetLine() - 1,
					Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
				},
			},
		})
	} else { // Set parent node as bool (node created when key was found)
		parentNode.Type = Bool
		parentNode.Value = boolValue
		parentNode.ValueLocation = TokenLocation{
			Start: CharLocation{
				Line:   ctx.GetStart().GetLine() - 1,
				Column: ctx.GetStart().GetColumn(),
			},
			End: CharLocation{
				Line:   ctx.GetStop().GetLine() - 1,
				Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
			},
		}
	}
}

// EnterDate_time is called when production date_time is entered.
// Parsed as string
func (p *tomlParser) EnterDate_time(ctx *parser_toml.Date_timeContext) {
	// Get parent node in stack
	parentNode := p.stack.Peek().(*Node)

	// If parent node is an array, append the string to the array
	if parentNode.Type == Array {
		parentNode.Value = append(parentNode.Value.([]*Node), &Node{
			Type:  String,
			Value: p.cleanString(ctx.GetText()),
			ValueLocation: TokenLocation{
				Start: CharLocation{
					Line:   ctx.GetStart().GetLine() - 1,
					Column: ctx.GetStart().GetColumn(),
				},
				End: CharLocation{
					Line:   ctx.GetStop().GetLine() - 1,
					Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
				},
			},
		})
	} else { // Set parent node as string (node created when key was found)
		parentNode.Type = String
		parentNode.Value = p.cleanString(ctx.GetText())
		parentNode.ValueLocation = TokenLocation{
			Start: CharLocation{
				Line:   ctx.GetStart().GetLine() - 1,
				Column: ctx.GetStart().GetColumn(),
			},
			End: CharLocation{
				Line:   ctx.GetStop().GetLine() - 1,
				Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
			},
		}
	}
}

// parseKey parses a key and returns the key after removing in-between spaces
// and cleaning each simple key string
func (p *tomlParser) parseKey(ctx parser_toml.IKeyContext) NodeKey {
	if ctx.Simple_key() != nil {
		return NodeKey{Segments: []string{p.cleanString(ctx.Simple_key().GetText())}}
	} else if ctx.Dotted_key() != nil {
		return p.parseDottedKey(ctx.Dotted_key())
	}

	panic("This should never happen")
}

// parseDottedKey parses a dotted key and returns the parent key and the field key.
func (p *tomlParser) parseDottedKey(ctx parser_toml.IDotted_keyContext) NodeKey {
	result := NodeKey{
		Segments: []string{},
	}
	for _, key := range ctx.AllSimple_key() {
		result.Segments = append(result.Segments, p.cleanString(key.GetText()))
	}

	return result
}

// Gets a node from a path starting at parentNode. If the node doesn't exist
// it gets created as a Null node. Intermediate nodes are created as Object nodes.
// If an array is found in the path, the last item in the array is used.
func (p *tomlParser) getOrCreateNode(parentNode *Node, key NodeKey, mode tomlTableDefinitionMode) (*Node, error) {
	currentNode := parentNode
	segments := key.Segments

	for index := 0; index < len(segments); index++ {
		segment := segments[index]

		switch currentNode.Type {
		case Object:
			// Cast value as map[string]*Node (unsafe)
			objMap := currentNode.Value.(map[string]*Node)

			// Check if the segment exists in the map
			if nextNode, exists := objMap[segment]; exists {
				// Check if inline table, these can't be traversed
				if p.definedTables[fmt.Sprintf("%p", nextNode)] == tomlDefModeInlineTable {
					return nil, fmt.Errorf("can't redefine existing key: '%s'", segment)
				}

				// If last segment, check mode
				if index == len(segments)-1 && p.definedTables[fmt.Sprintf("%p", nextNode)] != mode {
					return nil, fmt.Errorf("can't redefine existing key: '%s'", segment)
				}

				currentNode = nextNode
			} else {
				var newNode *Node

				// If we are at the last segment, create a null node
				if index == len(segments)-1 {
					newNode = &Node{
						Type:  Null,
						Value: nil,
					}
				} else { // Otherwise create an object node
					newNode = &Node{
						Type:  Object,
						Value: map[string]*Node{},
					}
				}

				p.definedTables[fmt.Sprintf("%p", newNode)] = mode

				// Add the new node to the map
				objMap[segment] = newNode

				// Set the current node to the new node
				currentNode = newNode
			}

		case Array:
			// Cast value as []*Node (unsafe)
			arrayValue := currentNode.Value.([]*Node)

			// Get last item in array, if empty throw error
			if len(arrayValue) == 0 {
				return nil, fmt.Errorf("cannot traverse empty array in path %s", key.String())
			}

			// Get last item in array
			lastItem := arrayValue[len(arrayValue)-1]

			currentNode = lastItem

			index-- // Decrement index to avoid skipping the next segment

		default:
			// If we are here, it means we're trying to traverse a leaf node
			return nil, fmt.Errorf("can't redefine existing key to contain: '%s'", segment)
		}
	}

	return currentNode, nil
}

func (p *tomlParser) cleanString(stringValue string) string {
	if strings.HasPrefix(stringValue, "\"\"\"") && strings.HasSuffix(stringValue, "\"\"\"") { // Check if it's ML basic string
		return stringValue[3 : len(stringValue)-3]
	} else if strings.HasPrefix(stringValue, "'''") && strings.HasSuffix(stringValue, "'''") { // Check if it's ML literal string
		return stringValue[3 : len(stringValue)-3]
	} else if strings.HasPrefix(stringValue, "\"") && strings.HasSuffix(stringValue, "\"") { // Check if it's basic string
		return stringValue[1 : len(stringValue)-1]
	} else if strings.HasPrefix(stringValue, "'") && strings.HasSuffix(stringValue, "'") { // Check if it's literal string
		return stringValue[1 : len(stringValue)-1]
	}

	return stringValue
}
