package langsrv

import (
	"github.com/ConfigMate/configmate/parsers/gen/parser_cmsl"
	"github.com/antlr4-go/antlr/v4"
)

type SemanticTokenProvider interface {
	GetSemanticTokens(content []byte) ([]ParsedToken, error)
}

type ParsedToken struct {
	Line      int               `json:"line"`
	Column    int               `json:"column"`
	Length    int               `json:"length"`
	TokenType SemanticTokenType `json:"tokenType"`
}

type SemanticTokenType string

const (
	STTKeyword   SemanticTokenType = "keyword"
	STTVariable  SemanticTokenType = "variable"
	STTProperty  SemanticTokenType = "property"
	STTType      SemanticTokenType = "type"
	STTDecorator SemanticTokenType = "decorator"
	STTMethod    SemanticTokenType = "method"
	STTString    SemanticTokenType = "string"
	STTNumber    SemanticTokenType = "number"
	STTOperator  SemanticTokenType = "operator"
)

func NewSemanticTokenProvider() SemanticTokenProvider {
	return &semanticTokenProviderImpl{}
}

type semanticTokenProviderImpl struct {
	*parser_cmsl.BaseCMSLListener
	tokens []ParsedToken
}

func (p *semanticTokenProviderImpl) GetSemanticTokens(content []byte) ([]ParsedToken, error) {
	// Create lexer
	input := antlr.NewInputStream(string(content))
	lexer := parser_cmsl.NewCMSLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := parser_cmsl.NewCMSLParser(stream)
	tree := parser.Cmsl()

	// Zero out the tokens
	p.tokens = make([]ParsedToken, 0)

	// Walk the tree
	walker := antlr.NewParseTreeWalker()
	walker.Walk(p, tree)

	return p.tokens, nil
}

// EnterFileDeclaration is called when production fileDeclaration is entered.
func (s *semanticTokenProviderImpl) EnterConfigDeclaration(ctx *parser_cmsl.ConfigDeclarationContext) {
	// Add the config keyword token
	configKeyword := ctx.CONFIG_DCLR_KW()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      configKeyword.GetSymbol().GetLine() - 1,
		Column:    configKeyword.GetSymbol().GetColumn(),
		Length:    len(configKeyword.GetText()),
		TokenType: STTKeyword,
	})

	// Add file path token
	filePath := ctx.SHORT_STRING()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      filePath.GetSymbol().GetLine() - 1,
		Column:    filePath.GetSymbol().GetColumn(),
		Length:    len(filePath.GetText()),
		TokenType: STTString,
	})

	// Add the file format token
	fileFormat := ctx.IDENTIFIER()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      fileFormat.GetSymbol().GetLine() - 1,
		Column:    fileFormat.GetSymbol().GetColumn(),
		Length:    len(fileFormat.GetText()),
		TokenType: STTDecorator,
	})
}

// EnterImportStatement is called when production importStatement is entered.
func (s *semanticTokenProviderImpl) EnterImportStatement(ctx *parser_cmsl.ImportStatementContext) {
	// Add the import keyword token
	importKeyword := ctx.IMPORT_KW()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      importKeyword.GetSymbol().GetLine() - 1,
		Column:    importKeyword.GetSymbol().GetColumn(),
		Length:    len(importKeyword.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterImportItem is called when production importItem is entered.
func (s *semanticTokenProviderImpl) EnterImportItem(ctx *parser_cmsl.ImportItemContext) {
	// Add the import alias token
	importAlias := ctx.IDENTIFIER()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      importAlias.GetSymbol().GetLine() - 1,
		Column:    importAlias.GetSymbol().GetColumn(),
		Length:    len(importAlias.GetText()),
		TokenType: STTVariable,
	})

	// Add file path token
	filePath := ctx.SHORT_STRING()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      filePath.GetSymbol().GetLine() - 1,
		Column:    filePath.GetSymbol().GetColumn(),
		Length:    len(filePath.GetText()),
		TokenType: STTString,
	})
}

// EnterSpecificationBody is called when production specificationBody is entered.
func (s *semanticTokenProviderImpl) EnterSpecificationBody(ctx *parser_cmsl.SpecificationBodyContext) {
	// Add the specification keyword token
	specKeyword := ctx.SPEC_ROOT_KW()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      specKeyword.GetSymbol().GetLine() - 1,
		Column:    specKeyword.GetSymbol().GetColumn(),
		Length:    len(specKeyword.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterObjectDefinitions is called when production objectDefinitions is entered.
func (s *semanticTokenProviderImpl) EnterObjectDefinitions(ctx *parser_cmsl.ObjectDefinitionsContext) {
	// Add the object keyword token
	objectKeyword := ctx.OBJ_DEF_KW()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      objectKeyword.GetSymbol().GetLine() - 1,
		Column:    objectKeyword.GetSymbol().GetColumn(),
		Length:    len(objectKeyword.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterObjectDefinition is called when production objectDefinition is entered.
func (s *semanticTokenProviderImpl) EnterObjectDefinition(ctx *parser_cmsl.ObjectDefinitionContext) {
	// Add the object type name
	objectTypeName := ctx.IDENTIFIER()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      objectTypeName.GetSymbol().GetLine() - 1,
		Column:    objectTypeName.GetSymbol().GetColumn(),
		Length:    len(objectTypeName.GetText()),
		TokenType: STTType,
	})
}

// ObjectPropertyDefinition is called when production objectPropertyDefinition is entered.
func (s *semanticTokenProviderImpl) EnterObjectPropertyDefinition(ctx *parser_cmsl.ObjectPropertyDefinitionContext) {
	// Add the object property type name
	objectPropertyTypeName := ctx.SimpleName()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      objectPropertyTypeName.GetStart().GetLine() - 1,
		Column:    objectPropertyTypeName.GetStart().GetColumn(),
		Length:    len(objectPropertyTypeName.GetText()),
		TokenType: STTVariable,
	})
}

// EnterFieldName is called when production fieldName is entered.
func (s *semanticTokenProviderImpl) EnterFieldName(ctx *parser_cmsl.FieldNameContext) {
	// Add the field name token
	s.tokens = append(s.tokens, ParsedToken{
		Line:      ctx.GetStart().GetLine() - 1,
		Column:    ctx.GetStart().GetColumn(),
		Length:    len(ctx.GetText()),
		TokenType: STTVariable,
	})
}

// EnterShortMetadataExpression is called when production shortMetadataExpression is entered.
func (s *semanticTokenProviderImpl) EnterShortMetadataExpression(ctx *parser_cmsl.ShortMetadataExpressionContext) {
	if ctx.OPTIONAL_METAD_KW() != nil {
		// Add the optional keyword token
		optionalKeyword := ctx.OPTIONAL_METAD_KW()
		s.tokens = append(s.tokens, ParsedToken{
			Line:      optionalKeyword.GetSymbol().GetLine() - 1,
			Column:    optionalKeyword.GetSymbol().GetColumn(),
			Length:    len(optionalKeyword.GetText()),
			TokenType: STTKeyword,
		})
	}
}

// EnterTypeMetadata is called when production typeMetadata is entered.
func (s *semanticTokenProviderImpl) EnterTypeMetadata(ctx *parser_cmsl.TypeMetadataContext) {
	// Add the type keyword token
	typeKeyword := ctx.TYPE_METAD_KW()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      typeKeyword.GetSymbol().GetLine() - 1,
		Column:    typeKeyword.GetSymbol().GetColumn(),
		Length:    len(typeKeyword.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterNotesMetadata is called when production notesMetadata is entered.
func (s *semanticTokenProviderImpl) EnterNotesMetadata(ctx *parser_cmsl.NotesMetadataContext) {
	// Add the notes keyword token
	notesKeyword := ctx.NOTES_METAD_KW()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      notesKeyword.GetSymbol().GetLine() - 1,
		Column:    notesKeyword.GetSymbol().GetColumn(),
		Length:    len(notesKeyword.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterDefaultMetadata is called when production defaultMetadata is entered.
func (s *semanticTokenProviderImpl) EnterDefaultMetadata(ctx *parser_cmsl.DefaultMetadataContext) {
	// Add the default keyword token
	defaultKeyword := ctx.DEFAULT_METAD_KW()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      defaultKeyword.GetSymbol().GetLine() - 1,
		Column:    defaultKeyword.GetSymbol().GetColumn(),
		Length:    len(defaultKeyword.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterOptionalMetadata is called when production optionalMetadata is entered.
func (s *semanticTokenProviderImpl) EnterOptionalMetadata(ctx *parser_cmsl.OptionalMetadataContext) {
	// Add the optional keyword token
	optionalKeyword := ctx.OPTIONAL_METAD_KW()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      optionalKeyword.GetSymbol().GetLine() - 1,
		Column:    optionalKeyword.GetSymbol().GetColumn(),
		Length:    len(optionalKeyword.GetText()),
		TokenType: STTKeyword,
	})

	// Add the boolean token
	booleanKeyword := ctx.BOOL()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      booleanKeyword.GetSymbol().GetLine() - 1,
		Column:    booleanKeyword.GetSymbol().GetColumn(),
		Length:    len(booleanKeyword.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterTypeExpr is called when production typeExpr is entered.
func (s *semanticTokenProviderImpl) EnterTypeExpr(ctx *parser_cmsl.TypeExprContext) {
	// Check if the type is a primitive
	if ctx.IDENTIFIER() != nil {
		// Add the type token
		s.tokens = append(s.tokens, ParsedToken{
			Line:      ctx.GetStart().GetLine() - 1,
			Column:    ctx.GetStart().GetColumn(),
			Length:    len(ctx.GetText()),
			TokenType: STTType,
		})
	} else if ctx.LIST_TYPE_KW() != nil {
		// Get the list type keyword
		listKeyword := ctx.LIST_TYPE_KW()
		s.tokens = append(s.tokens, ParsedToken{
			Line:      listKeyword.GetSymbol().GetLine() - 1,
			Column:    listKeyword.GetSymbol().GetColumn(),
			Length:    len(listKeyword.GetText()),
			TokenType: STTType,
		})
	}
}

// EnterString is called when production string is entered.
func (s *semanticTokenProviderImpl) EnterString(ctx *parser_cmsl.StringContext) {
	// Add the string token
	s.tokens = append(s.tokens, ParsedToken{
		Line:      ctx.GetStart().GetLine() - 1,
		Column:    ctx.GetStart().GetColumn(),
		Length:    len(ctx.GetText()),
		TokenType: STTString,
	})
}

// EnterInt is called when production int is entered.
func (s *semanticTokenProviderImpl) EnterInt(ctx *parser_cmsl.IntContext) {
	// Add the int token
	s.tokens = append(s.tokens, ParsedToken{
		Line:      ctx.GetStart().GetLine() - 1,
		Column:    ctx.GetStart().GetColumn(),
		Length:    len(ctx.GetText()),
		TokenType: STTNumber,
	})
}

// EnterFloat is called when production float is entered.
func (s *semanticTokenProviderImpl) EnterFloat(ctx *parser_cmsl.FloatContext) {
	// Add the float token
	s.tokens = append(s.tokens, ParsedToken{
		Line:      ctx.GetStart().GetLine() - 1,
		Column:    ctx.GetStart().GetColumn(),
		Length:    len(ctx.GetText()),
		TokenType: STTNumber,
	})
}

// EnterBoolean is called when production boolean is entered.
func (s *semanticTokenProviderImpl) EnterBoolean(ctx *parser_cmsl.BooleanContext) {
	// Add the boolean token
	s.tokens = append(s.tokens, ParsedToken{
		Line:      ctx.GetStart().GetLine() - 1,
		Column:    ctx.GetStart().GetColumn(),
		Length:    len(ctx.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterStringExpr is called when production stringExpr is entered.
func (s *semanticTokenProviderImpl) EnterStringExpr(ctx *parser_cmsl.StringExprContext) {
	// Add the string token
	s.tokens = append(s.tokens, ParsedToken{
		Line:      ctx.GetStart().GetLine() - 1,
		Column:    ctx.GetStart().GetColumn(),
		Length:    len(ctx.GetText()),
		TokenType: STTString,
	})
}

// EnterOrExpr is called when production orExpr is entered.
func (s *semanticTokenProviderImpl) EnterOrExpr(ctx *parser_cmsl.OrExprContext) {
	// Get the token of each OR operator
	for _, orOperator := range ctx.AllOR_SYM() {
		s.tokens = append(s.tokens, ParsedToken{
			Line:      orOperator.GetSymbol().GetLine() - 1,
			Column:    orOperator.GetSymbol().GetColumn(),
			Length:    len(orOperator.GetText()),
			TokenType: STTOperator,
		})
	}
}

// EnterAndExpr is called when production andExpr is entered.
func (s *semanticTokenProviderImpl) EnterAndExpr(ctx *parser_cmsl.AndExprContext) {
	// Get the token of each AND operator
	for _, andOperator := range ctx.AllAND_SYM() {
		s.tokens = append(s.tokens, ParsedToken{
			Line:      andOperator.GetSymbol().GetLine() - 1,
			Column:    andOperator.GetSymbol().GetColumn(),
			Length:    len(andOperator.GetText()),
			TokenType: STTOperator,
		})
	}
}

// EnterIf is called when production if is entered.
func (s *semanticTokenProviderImpl) EnterIf(ctx *parser_cmsl.IfContext) {
	// Add the if keyword token
	ifKeyword := ctx.IF_SYM()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      ifKeyword.GetSymbol().GetLine() - 1,
		Column:    ifKeyword.GetSymbol().GetColumn(),
		Length:    len(ifKeyword.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterElseif is called when production elseif is entered.
func (s *semanticTokenProviderImpl) EnterElseif(ctx *parser_cmsl.ElseifContext) {
	// Add the elseif keyword token
	elseifKeyword := ctx.ELSEIF_SYM()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      elseifKeyword.GetSymbol().GetLine() - 1,
		Column:    elseifKeyword.GetSymbol().GetColumn(),
		Length:    len(elseifKeyword.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterElse is called when production else is entered.
func (s *semanticTokenProviderImpl) EnterElse(ctx *parser_cmsl.ElseContext) {
	// Add the else keyword token
	elseKeyword := ctx.ELSE_SYM()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      elseKeyword.GetSymbol().GetLine() - 1,
		Column:    elseKeyword.GetSymbol().GetColumn(),
		Length:    len(elseKeyword.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterForeach is called when production foreach is entered.
func (s *semanticTokenProviderImpl) EnterForeach(ctx *parser_cmsl.ForeachContext) {
	// Add the foreach keyword token
	foreachKeyword := ctx.FOREACH_SYM()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      foreachKeyword.GetSymbol().GetLine() - 1,
		Column:    foreachKeyword.GetSymbol().GetColumn(),
		Length:    len(foreachKeyword.GetText()),
		TokenType: STTKeyword,
	})

	// Add the in loop identifier
	inKeyword := ctx.IDENTIFIER()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      inKeyword.GetSymbol().GetLine() - 1,
		Column:    inKeyword.GetSymbol().GetColumn(),
		Length:    len(inKeyword.GetText()),
		TokenType: STTVariable,
	})
}

// EnterNot is called when production not is entered.
func (s *semanticTokenProviderImpl) EnterNot(ctx *parser_cmsl.NotContext) {
	// Add the not keyword token
	notKeyword := ctx.NOT_SYM()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      notKeyword.GetSymbol().GetLine() - 1,
		Column:    notKeyword.GetSymbol().GetColumn(),
		Length:    len(notKeyword.GetText()),
		TokenType: STTKeyword,
	})
}

// EnterFunction is called when production function is entered.
func (s *semanticTokenProviderImpl) EnterFunction(ctx *parser_cmsl.FunctionContext) {
	// Get function name
	functionName := ctx.IDENTIFIER()
	s.tokens = append(s.tokens, ParsedToken{
		Line:      functionName.GetSymbol().GetLine() - 1,
		Column:    functionName.GetSymbol().GetColumn(),
		Length:    len(functionName.GetText()),
		TokenType: STTMethod,
	})
}
