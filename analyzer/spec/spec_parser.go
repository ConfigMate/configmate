package spec

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ConfigMate/configmate/parsers"
	parser_cmsl "github.com/ConfigMate/configmate/parsers/gen/parser_cmsl/parsers/grammars"
	"github.com/antlr4-go/antlr/v4"
	"github.com/golang-collections/collections/stack"
)

type SpecParser interface {
	Parse(spec []byte) (*Specification, []SpecParserError)
}

type SpecParserError struct {
	ErrorMessage string
	Location     parsers.TokenLocation
}

func NewSpecParser() SpecParser {
	return &specParserImpl{}
}

type cmslErrorListener struct {
	*antlr.DefaultErrorListener
	errors []SpecParserError
}

func (d *cmslErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{},
	line, column int, msg string, e antlr.RecognitionException) {
	d.errors = append(d.errors, SpecParserError{
		ErrorMessage: msg,
		Location: parsers.TokenLocation{
			Start: parsers.CharLocation{
				Line:   line,
				Column: column,
			},
			End: parsers.CharLocation{
				Line:   line,
				Column: column + 1,
			},
		},
	})
}

type specParserImpl struct {
	*parser_cmsl.BaseCMSLListener

	spec           Specification
	itemFieldStack stack.Stack
	errs           []SpecParserError
}

func (p *specParserImpl) Parse(spec []byte) (*Specification, []SpecParserError) {
	// Create error listener
	errorListener := &cmslErrorListener{}

	// Create lexer
	input := antlr.NewInputStream(string(spec))
	lexer := parser_cmsl.NewCMSLLexer(input)

	// Add error listener
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := parser_cmsl.NewCMSLParser(stream)

	// Add error listener
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errorListener)

	tree := parser.Cmsl()

	// Check for errors
	if len(errorListener.errors) > 0 {
		return nil, errorListener.errors
	}

	// Zero out the spec and errs
	p.spec = Specification{
		Imports:              make(map[string]string),
		ImportsAliasLocation: make(map[string]parsers.TokenLocation),
		ImportsLocation:      make(map[string]parsers.TokenLocation),
		Fields:               make([]FieldSpec, 0),
	}
	p.errs = nil

	// Prepare stack
	p.itemFieldStack = stack.Stack{}

	// Walk the tree
	walker := antlr.NewParseTreeWalker()
	walker.Walk(p, tree)

	// Check for errors
	if len(p.errs) > 0 {
		return nil, p.errs
	}

	return &p.spec, nil
}

// EnterFileDeclaration is called when production fileDeclaration is entered.
func (p *specParserImpl) EnterFileDeclaration(ctx *parser_cmsl.FileDeclarationContext) {
	// Set values of file and fileLocation in spec
	p.spec.File = removeStrQuotesAndCleanSpaces(ctx.SHORT_STRING().GetText())
	p.spec.FileLocation = parsers.TokenLocation{
		Start: parsers.CharLocation{
			Line:   ctx.SHORT_STRING().GetSymbol().GetLine(),
			Column: ctx.SHORT_STRING().GetSymbol().GetColumn(),
		},
		End: parsers.CharLocation{
			Line:   ctx.SHORT_STRING().GetSymbol().GetLine(),
			Column: ctx.SHORT_STRING().GetSymbol().GetColumn() + len(ctx.SHORT_STRING().GetText()),
		},
	}

	// Set values of fileFormat and fileFormatLocation in spec
	p.spec.FileFormat = ctx.IDENTIFIER().GetText()
	p.spec.FileFormatLocation = parsers.TokenLocation{
		Start: parsers.CharLocation{
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine(),
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn(),
		},
		End: parsers.CharLocation{
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine(),
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn() + len(ctx.IDENTIFIER().GetText()),
		},
	}
}

// EnterImportStatement is called when production importStatement is entered.
func (p *specParserImpl) EnterImportItem(ctx *parser_cmsl.ImportItemContext) {
	// Add import to spec
	p.spec.Imports[ctx.IDENTIFIER().GetText()] = removeStrQuotesAndCleanSpaces(ctx.SHORT_STRING().GetText())
	p.spec.ImportsLocation[ctx.IDENTIFIER().GetText()] = parsers.TokenLocation{
		Start: parsers.CharLocation{
			Line:   ctx.SHORT_STRING().GetSymbol().GetLine(),
			Column: ctx.SHORT_STRING().GetSymbol().GetColumn(),
		},
		End: parsers.CharLocation{
			Line:   ctx.SHORT_STRING().GetSymbol().GetLine(),
			Column: ctx.SHORT_STRING().GetSymbol().GetColumn() + len(ctx.SHORT_STRING().GetText()),
		},
	}

	// Add import alias location
	p.spec.ImportsAliasLocation[ctx.IDENTIFIER().GetText()] = parsers.TokenLocation{
		Start: parsers.CharLocation{
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine(),
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn(),
		},
		End: parsers.CharLocation{
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine(),
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn() + len(ctx.IDENTIFIER().GetText()),
		},
	}
}

// EnterSpecificationItem is called when production specificationItem is entered.
func (p *specParserImpl) EnterSpecificationItem(ctx *parser_cmsl.SpecificationItemContext) {
	// Compute fieldName
	fieldName := ctx.FieldName().GetText()
	if p.itemFieldStack.Len() != 0 {
		parentField := p.itemFieldStack.Peek().(string)
		fieldName = parentField + "." + fieldName
	}

	// Add item to stack
	p.itemFieldStack.Push(fieldName)

	// Add item to spec
	fieldSpecification := FieldSpec{
		Field: fieldName,
		FieldLocation: parsers.TokenLocation{
			Start: parsers.CharLocation{
				Line:   ctx.FieldName().GetStart().GetLine(),
				Column: ctx.FieldName().GetStart().GetColumn(),
			},
			End: parsers.CharLocation{
				Line:   ctx.FieldName().GetStop().GetLine(),
				Column: ctx.FieldName().GetStop().GetColumn() + len(ctx.FieldName().GetStop().GetText()),
			},
		},
		Checks: make([]CheckWithLocation, 0),
	}

	foundType := false
	foundDefault := false
	foundOptional := false
	foundNotes := false

	// For each metadata item
	for _, metadataItem := range ctx.MetadataExpression().AllMetadataItem() {
		switch item := metadataItem.(type) {
		case *parser_cmsl.TypeMetadataContext:
			// Check if type has already been found
			if foundType {
				p.errs = append(p.errs, SpecParserError{
					ErrorMessage: fmt.Sprintf("duplicate type metadata for field %s", fieldName),
					Location: parsers.TokenLocation{
						Start: parsers.CharLocation{
							Line:   item.GetStart().GetLine(),
							Column: item.GetStart().GetColumn(),
						},
						End: parsers.CharLocation{
							Line:   item.GetStop().GetLine(),
							Column: item.GetStop().GetColumn() + len(item.GetStop().GetText()),
						},
					},
				})
				continue
			}
			foundType = true

			// Add type to field
			fieldSpecification.Type = condenseListExpressions(item.TypeExpr().GetText())
			fieldSpecification.TypeLocation = parsers.TokenLocation{
				Start: parsers.CharLocation{
					Line:   item.TypeExpr().GetStart().GetLine(),
					Column: item.TypeExpr().GetStart().GetColumn(),
				},
				End: parsers.CharLocation{
					Line:   item.TypeExpr().GetStop().GetLine(),
					Column: item.TypeExpr().GetStop().GetColumn() + len(item.TypeExpr().GetStop().GetText()),
				},
			}
		case *parser_cmsl.OptionalMetadataContext:
			// Check if optional has already been found
			if foundOptional {
				p.errs = append(p.errs, SpecParserError{
					ErrorMessage: fmt.Sprintf("duplicate optional metadata for field %s", fieldName),
					Location: parsers.TokenLocation{
						Start: parsers.CharLocation{
							Line:   item.GetStart().GetLine(),
							Column: item.GetStart().GetColumn(),
						},
						End: parsers.CharLocation{
							Line:   item.GetStop().GetLine(),
							Column: item.GetStop().GetColumn() + len(item.GetStop().GetText()),
						},
					},
				})
				continue
			}
			foundOptional = true

			// Add optional to field
			optional, err := strconv.ParseBool(item.BOOL().GetText())
			if err != nil {
				panic(fmt.Sprintf("optional must be a bool, found: %s; this error should have been cought in a previous stage", item.BOOL().GetText()))
			}

			fieldSpecification.Optional = optional
			fieldSpecification.OptionalLocation = parsers.TokenLocation{
				Start: parsers.CharLocation{
					Line:   item.BOOL().GetSymbol().GetLine(),
					Column: item.BOOL().GetSymbol().GetColumn(),
				},
				End: parsers.CharLocation{
					Line:   item.BOOL().GetSymbol().GetLine(),
					Column: item.BOOL().GetSymbol().GetColumn() + len(item.BOOL().GetSymbol().GetText()),
				},
			}
		case *parser_cmsl.DefaultMetadataContext:
			// Check if default has already been found
			if foundDefault {
				p.errs = append(p.errs, SpecParserError{
					ErrorMessage: fmt.Sprintf("duplicate default metadata for field %s", fieldName),
					Location: parsers.TokenLocation{
						Start: parsers.CharLocation{
							Line:   item.GetStart().GetLine(),
							Column: item.GetStart().GetColumn(),
						},
						End: parsers.CharLocation{
							Line:   item.GetStop().GetLine(),
							Column: item.GetStop().GetColumn() + len(item.GetStop().GetText()),
						},
					},
				})
				continue
			}
			foundDefault = true

			// Add default to field
			fieldSpecification.Default = removeStrQuotesAndCleanSpaces(item.Primitive().GetText())
			fieldSpecification.DefaultLocation = parsers.TokenLocation{
				Start: parsers.CharLocation{
					Line:   item.Primitive().GetStart().GetLine(),
					Column: item.Primitive().GetStart().GetColumn(),
				},
				End: parsers.CharLocation{
					Line:   item.Primitive().GetStop().GetLine(),
					Column: item.Primitive().GetStop().GetColumn() + len(item.Primitive().GetStop().GetText()),
				},
			}
		case *parser_cmsl.NotesMetadataContext:
			// Check if notes has already been found
			if foundNotes {
				p.errs = append(p.errs, SpecParserError{
					ErrorMessage: fmt.Sprintf("duplicate notes metadata for field %s", fieldName),
					Location: parsers.TokenLocation{
						Start: parsers.CharLocation{
							Line:   item.GetStart().GetLine(),
							Column: item.GetStart().GetColumn(),
						},
						End: parsers.CharLocation{
							Line:   item.GetStop().GetLine(),
							Column: item.GetStop().GetColumn() + len(item.GetStop().GetText()),
						},
					},
				})
				continue
			}
			foundNotes = true

			// Add notes to field
			fieldSpecification.Notes = removeStrQuotesAndCleanSpaces(item.StringExpr().GetText())
			fieldSpecification.NotesLocation = parsers.TokenLocation{
				Start: parsers.CharLocation{
					Line:   item.StringExpr().GetStart().GetLine(),
					Column: item.StringExpr().GetStart().GetColumn(),
				},
				End: parsers.CharLocation{
					Line:   item.StringExpr().GetStop().GetLine(),
					Column: item.StringExpr().GetStop().GetColumn() + len(item.StringExpr().GetStop().GetText()),
				},
			}

		default:
			panic(fmt.Sprintf("unknown metadata item: %s; this error should have been cought in a previous stage", item.GetText()))
		}
	}

	if !foundType {
		p.errs = append(p.errs, SpecParserError{
			ErrorMessage: fmt.Sprintf("missing type metadata for field %s", fieldName),
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{
					Line:   ctx.GetStart().GetLine(),
					Column: ctx.GetStart().GetColumn(),
				},
				End: parsers.CharLocation{
					Line:   ctx.GetStop().GetLine(),
					Column: ctx.GetStop().GetColumn() + len(ctx.GetStop().GetText()),
				},
			},
		})
	}

	// For each check statement
	for _, check := range ctx.AllCheck() {
		// Add check to field
		fieldSpecification.Checks = append(fieldSpecification.Checks, CheckWithLocation{
			Check: check.GetText(),
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{
					Line:   check.GetStart().GetLine(),
					Column: check.GetStart().GetColumn(),
				},
				End: parsers.CharLocation{
					Line:   check.GetStop().GetLine(),
					Column: check.GetStop().GetColumn() + len(check.GetStop().GetText()),
				},
			},
		})
	}

	p.spec.Fields = append(p.spec.Fields, fieldSpecification)
}

// ExitObjectField is called when production objectField is exited.
func (p *specParserImpl) ExitSpecificationItem(ctx *parser_cmsl.SpecificationItemContext) {
	// Pop field from stack
	p.itemFieldStack.Pop()
}

func removeStrQuotesAndCleanSpaces(str string) string {
	if strings.HasPrefix(str, "\"\"\"") {
		// Remove triple quotes
		str = str[3 : len(str)-3]
	} else if strings.HasPrefix(str, "\"") {
		// Remove quotes
		str = str[1 : len(str)-1]
	}

	// Remove consecutive spaces
	space := regexp.MustCompile(`\s+`)
	str = space.ReplaceAllString(str, " ")

	// Remove spaces at the beginning and end
	str = strings.TrimSpace(str)

	return str
}

func condenseListExpressions(typeStr string) string {
	result := ""

	for strings.HasPrefix(typeStr, "list<") && strings.HasSuffix(typeStr, ">") {
		result = result + "list:"
		typeStr = typeStr[5 : len(typeStr)-1]
	}

	result = result + typeStr
	return result
}
