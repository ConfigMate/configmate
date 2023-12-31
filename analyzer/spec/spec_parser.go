package spec

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ConfigMate/configmate/parsers"
	"github.com/ConfigMate/configmate/parsers/gen/parser_cmsl"
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
				Line:   line - 1,
				Column: column,
			},
			End: parsers.CharLocation{
				Line:   line - 1,
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
func (p *specParserImpl) EnterConfigDeclaration(ctx *parser_cmsl.ConfigDeclarationContext) {
	// Set values of file and fileLocation in spec
	p.spec.File = removeStrQuotesAndCleanSpaces(ctx.SHORT_STRING().GetText())
	p.spec.FileLocation = parsers.TokenLocation{
		Start: parsers.CharLocation{
			Line:   ctx.SHORT_STRING().GetSymbol().GetLine() - 1,
			Column: ctx.SHORT_STRING().GetSymbol().GetColumn(),
		},
		End: parsers.CharLocation{
			Line:   ctx.SHORT_STRING().GetSymbol().GetLine() - 1,
			Column: ctx.SHORT_STRING().GetSymbol().GetColumn() + len(ctx.SHORT_STRING().GetText()),
		},
	}

	// Set values of fileFormat and fileFormatLocation in spec
	p.spec.FileFormat = ctx.IDENTIFIER().GetText()
	p.spec.FileFormatLocation = parsers.TokenLocation{
		Start: parsers.CharLocation{
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine() - 1,
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn(),
		},
		End: parsers.CharLocation{
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine() - 1,
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
			Line:   ctx.SHORT_STRING().GetSymbol().GetLine() - 1,
			Column: ctx.SHORT_STRING().GetSymbol().GetColumn(),
		},
		End: parsers.CharLocation{
			Line:   ctx.SHORT_STRING().GetSymbol().GetLine() - 1,
			Column: ctx.SHORT_STRING().GetSymbol().GetColumn() + len(ctx.SHORT_STRING().GetText()),
		},
	}

	// Add import alias location
	p.spec.ImportsAliasLocation[ctx.IDENTIFIER().GetText()] = parsers.TokenLocation{
		Start: parsers.CharLocation{
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine() - 1,
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn(),
		},
		End: parsers.CharLocation{
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine() - 1,
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn() + len(ctx.IDENTIFIER().GetText()),
		},
	}
}

// EnterSpecificationItem is called when production specificationItem is entered.
func (p *specParserImpl) EnterSpecificationItem(ctx *parser_cmsl.SpecificationItemContext) {
	// Compute the field key
	fieldKey := parseFieldName(ctx.FieldName())
	if p.itemFieldStack.Len() != 0 {
		parentFieldKey := p.itemFieldStack.Peek().(*parsers.NodeKey)
		fieldKey = parentFieldKey.Join(fieldKey)
	}

	// Add item to stack
	p.itemFieldStack.Push(fieldKey)

	// Add item to spec
	fieldSpecification := FieldSpec{
		Field: fieldKey,
		FieldLocation: parsers.TokenLocation{
			Start: parsers.CharLocation{
				Line:   ctx.FieldName().GetStart().GetLine() - 1,
				Column: ctx.FieldName().GetStart().GetColumn(),
			},
			End: parsers.CharLocation{
				Line:   ctx.FieldName().GetStop().GetLine() - 1,
				Column: ctx.FieldName().GetStop().GetColumn() + len(ctx.FieldName().GetStop().GetText()),
			},
		},
		Checks: make([]CheckWithLocation, 0),
	}

	foundType := false
	foundDefault := false
	foundOptional := false
	foundNotes := false

	if ctx.ShortMetadataExpression() != nil {
		foundType = true
		typeExpr := ctx.ShortMetadataExpression().TypeExpr()

		// Add type to field
		fieldSpecification.Type = typeExpr.GetText()
		fieldSpecification.TypeLocation = parsers.TokenLocation{
			Start: parsers.CharLocation{
				Line:   typeExpr.GetStart().GetLine() - 1,
				Column: typeExpr.GetStart().GetColumn(),
			},
			End: parsers.CharLocation{
				Line:   typeExpr.GetStop().GetLine() - 1,
				Column: typeExpr.GetStop().GetColumn() + len(typeExpr.GetStop().GetText()),
			},
		}

		if ctx.ShortMetadataExpression().OPTIONAL_METAD_KW() != nil {
			fieldSpecification.Optional = true
			fieldSpecification.OptionalLocation = parsers.TokenLocation{
				Start: parsers.CharLocation{
					Line:   ctx.ShortMetadataExpression().OPTIONAL_METAD_KW().GetSymbol().GetLine() - 1,
					Column: ctx.ShortMetadataExpression().OPTIONAL_METAD_KW().GetSymbol().GetColumn(),
				},
				End: parsers.CharLocation{
					Line:   ctx.ShortMetadataExpression().OPTIONAL_METAD_KW().GetSymbol().GetLine() - 1,
					Column: ctx.ShortMetadataExpression().OPTIONAL_METAD_KW().GetSymbol().GetColumn() + len(ctx.ShortMetadataExpression().OPTIONAL_METAD_KW().GetSymbol().GetText()),
				},
			}
		}
	} else if ctx.LongMetadataExpression() != nil {
		// For each metadata item
		for _, metadataItem := range ctx.LongMetadataExpression().AllMetadataItem() {
			switch item := metadataItem.(type) {
			case *parser_cmsl.TypeMetadataContext:
				// Check if type has already been found
				if foundType {
					p.errs = append(p.errs, SpecParserError{
						ErrorMessage: fmt.Sprintf("duplicate type metadata for field %s", fieldKey.String()),
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{
								Line:   item.GetStart().GetLine() - 1,
								Column: item.GetStart().GetColumn(),
							},
							End: parsers.CharLocation{
								Line:   item.GetStop().GetLine() - 1,
								Column: item.GetStop().GetColumn() + len(item.GetStop().GetText()),
							},
						},
					})
					continue
				}
				foundType = true

				// Add type to field
				fieldSpecification.Type = item.TypeExpr().GetText()
				fieldSpecification.TypeLocation = parsers.TokenLocation{
					Start: parsers.CharLocation{
						Line:   item.TypeExpr().GetStart().GetLine() - 1,
						Column: item.TypeExpr().GetStart().GetColumn(),
					},
					End: parsers.CharLocation{
						Line:   item.TypeExpr().GetStop().GetLine() - 1,
						Column: item.TypeExpr().GetStop().GetColumn() + len(item.TypeExpr().GetStop().GetText()),
					},
				}
			case *parser_cmsl.OptionalMetadataContext:
				// Check if optional has already been found
				if foundOptional {
					p.errs = append(p.errs, SpecParserError{
						ErrorMessage: fmt.Sprintf("duplicate optional metadata for field %s", fieldKey.String()),
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{
								Line:   item.GetStart().GetLine() - 1,
								Column: item.GetStart().GetColumn(),
							},
							End: parsers.CharLocation{
								Line:   item.GetStop().GetLine() - 1,
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
						Line:   item.BOOL().GetSymbol().GetLine() - 1,
						Column: item.BOOL().GetSymbol().GetColumn(),
					},
					End: parsers.CharLocation{
						Line:   item.BOOL().GetSymbol().GetLine() - 1,
						Column: item.BOOL().GetSymbol().GetColumn() + len(item.BOOL().GetSymbol().GetText()),
					},
				}
			case *parser_cmsl.DefaultMetadataContext:
				// Check if default has already been found
				if foundDefault {
					p.errs = append(p.errs, SpecParserError{
						ErrorMessage: fmt.Sprintf("duplicate default metadata for field %s", fieldKey.String()),
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{
								Line:   item.GetStart().GetLine() - 1,
								Column: item.GetStart().GetColumn(),
							},
							End: parsers.CharLocation{
								Line:   item.GetStop().GetLine() - 1,
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
						Line:   item.Primitive().GetStart().GetLine() - 1,
						Column: item.Primitive().GetStart().GetColumn(),
					},
					End: parsers.CharLocation{
						Line:   item.Primitive().GetStop().GetLine() - 1,
						Column: item.Primitive().GetStop().GetColumn() + len(item.Primitive().GetStop().GetText()),
					},
				}
			case *parser_cmsl.NotesMetadataContext:
				// Check if notes has already been found
				if foundNotes {
					p.errs = append(p.errs, SpecParserError{
						ErrorMessage: fmt.Sprintf("duplicate notes metadata for field %s", fieldKey.String()),
						Location: parsers.TokenLocation{
							Start: parsers.CharLocation{
								Line:   item.GetStart().GetLine() - 1,
								Column: item.GetStart().GetColumn(),
							},
							End: parsers.CharLocation{
								Line:   item.GetStop().GetLine() - 1,
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
						Line:   item.StringExpr().GetStart().GetLine() - 1,
						Column: item.StringExpr().GetStart().GetColumn(),
					},
					End: parsers.CharLocation{
						Line:   item.StringExpr().GetStop().GetLine() - 1,
						Column: item.StringExpr().GetStop().GetColumn() + len(item.StringExpr().GetStop().GetText()),
					},
				}

			default:
				panic(fmt.Sprintf("unknown metadata item: %s; this error should have been cought in a previous stage", item.GetText()))
			}
		}
	} else {
		panic(fmt.Sprintf("unknown metadata expression: %s; this error should have been cought in a previous stage", ctx.GetText()))
	}

	if !foundType {
		p.errs = append(p.errs, SpecParserError{
			ErrorMessage: fmt.Sprintf("missing type metadata for field %s", fieldKey.String()),
			Location: parsers.TokenLocation{
				Start: parsers.CharLocation{
					Line:   ctx.GetStart().GetLine() - 1,
					Column: ctx.GetStart().GetColumn(),
				},
				End: parsers.CharLocation{
					Line:   ctx.GetStop().GetLine() - 1,
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
					Line:   check.GetStart().GetLine() - 1,
					Column: check.GetStart().GetColumn(),
				},
				End: parsers.CharLocation{
					Line:   check.GetStop().GetLine() - 1,
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

// EnterObjectDefinition is called when production objectDefinition is entered.
func (p *specParserImpl) EnterObjectDefinition(ctx *parser_cmsl.ObjectDefinitionContext) {
	// Create object definition structure
	objectDefinition := ObjectDef{}

	// Get object type name and locations
	objectDefinition.Name = ctx.IDENTIFIER().GetText()
	objectDefinition.NameLocation = parsers.TokenLocation{
		Start: parsers.CharLocation{
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine() - 1,
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn(),
		},
		End: parsers.CharLocation{
			Line:   ctx.IDENTIFIER().GetSymbol().GetLine() - 1,
			Column: ctx.IDENTIFIER().GetSymbol().GetColumn() + len(ctx.IDENTIFIER().GetText()),
		},
	}

	// Add object properties
	for _, propertyDef := range ctx.AllObjectPropertyDefinition() {
		// Create object property definition structure
		objectPropertyDefinition := ObjectPropertyDef{}

		// Get property name
		objectPropertyDefinition.Name = removeSingleQuotesInKeys(propertyDef.SimpleName().GetText())
		objectPropertyDefinition.NameLocation = parsers.TokenLocation{
			Start: parsers.CharLocation{
				Line:   propertyDef.SimpleName().GetStart().GetLine() - 1,
				Column: propertyDef.SimpleName().GetStart().GetColumn(),
			},
			End: parsers.CharLocation{
				Line:   propertyDef.SimpleName().GetStop().GetLine() - 1,
				Column: propertyDef.SimpleName().GetStop().GetColumn() + len(propertyDef.SimpleName().GetStop().GetText()),
			},
		}

		// Get property type
		objectPropertyDefinition.Type = propertyDef.ShortMetadataExpression().TypeExpr().GetText()
		objectPropertyDefinition.TypeLocation = parsers.TokenLocation{
			Start: parsers.CharLocation{
				Line:   propertyDef.ShortMetadataExpression().TypeExpr().GetStart().GetLine() - 1,
				Column: propertyDef.ShortMetadataExpression().TypeExpr().GetStart().GetColumn(),
			},
			End: parsers.CharLocation{
				Line:   propertyDef.ShortMetadataExpression().TypeExpr().GetStop().GetLine() - 1,
				Column: propertyDef.ShortMetadataExpression().TypeExpr().GetStop().GetColumn() + len(propertyDef.ShortMetadataExpression().TypeExpr().GetStop().GetText()),
			},
		}

		// Get property optional
		if propertyDef.ShortMetadataExpression().OPTIONAL_METAD_KW() != nil {
			objectPropertyDefinition.Optional = true
			objectPropertyDefinition.OptionalLocation = parsers.TokenLocation{
				Start: parsers.CharLocation{
					Line:   propertyDef.ShortMetadataExpression().OPTIONAL_METAD_KW().GetSymbol().GetLine() - 1,
					Column: propertyDef.ShortMetadataExpression().OPTIONAL_METAD_KW().GetSymbol().GetColumn(),
				},
				End: parsers.CharLocation{
					Line:   propertyDef.ShortMetadataExpression().OPTIONAL_METAD_KW().GetSymbol().GetLine() - 1,
					Column: propertyDef.ShortMetadataExpression().OPTIONAL_METAD_KW().GetSymbol().GetColumn() + len(propertyDef.ShortMetadataExpression().OPTIONAL_METAD_KW().GetSymbol().GetText()),
				},
			}
		}

		// Add property to object definition
		objectDefinition.Properties = append(objectDefinition.Properties, objectPropertyDefinition)
	}

	// Add object definition to spec
	p.spec.Objects = append(p.spec.Objects, objectDefinition)
}

func parseFieldName(ctx parser_cmsl.IFieldNameContext) *parsers.NodeKey {
	if ctx.SimpleName() != nil {
		return &parsers.NodeKey{Segments: []string{removeSingleQuotesInKeys(ctx.SimpleName().GetText())}}
	} else if ctx.DottedName() != nil {
		return parseDottedName(ctx.DottedName())
	}

	panic(fmt.Sprintf("unknown field name: %s; this error should have been cought in a previous stage", ctx.GetText()))
}

func parseDottedName(ctx parser_cmsl.IDottedNameContext) *parsers.NodeKey {
	segments := make([]string, 0)

	// For each segment
	for _, segment := range ctx.AllSimpleName() {
		segments = append(segments, removeSingleQuotesInKeys(segment.GetText()))
	}

	return &parsers.NodeKey{Segments: segments}
}

func removeSingleQuotesInKeys(str string) string {
	if strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'") {
		// Remove quotes
		str = str[1 : len(str)-1]
	}
	return str
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
