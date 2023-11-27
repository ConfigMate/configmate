// This grammar describes the syntax of the language
// of specifications in ConfigMate. CMSL stands for ConfigMate Specification Language.
grammar CMSL;
import CMCL;

// The top-level rule of the grammar.
cmsl: specification EOF;

// A CMSL specification contains a file declaration, a list of imports,
// a specification body, and an optional list of custom object types.
specification: configDeclaration importStatement? specificationBody objectDefinitions?;

// A file declaration contains the path and format of the file.
configDeclaration: CONFIG_DCLR_KW COLON SHORT_STRING IDENTIFIER;

// An import contains the name of the file to import.
importStatement: IMPORT_KW LPAREN importItem (COMMA importItem)* RPAREN;

// An import item contains the name of the file to import.
importItem: IDENTIFIER COLON SHORT_STRING;

// A specification body contains a list of declarations.
specificationBody: SPEC_ROOT_KW LBRACE specificationItem* RBRACE;

// A collection of custom object types.
objectDefinitions: OBJ_DEF_KW LBRACE objectDefinition* RBRACE;

// A specification item starts with the field name, followed by the
// metadata inside angled brackets, optionally followed by a list of semicolon separated
// checks (CMCL expressions), and optionally followed with the specification of underlying
// fields insided curly braces.
specificationItem: fieldName (longMetadataExpression | shortMetadataExpression) ( LPAREN (check SEMICOLON)+ RPAREN )? (LBRACE specificationItem* RBRACE)?;

// A long metadata expression is a list of metadata items inside angled brackets.
longMetadataExpression : LANGLE metadataItem (COMMA metadataItem)* RANGLE ; 

// A short metadata expression is a type expression inside angled brackets.
shortMetadataExpression : LANGLE typeExpr RANGLE ;

// A metadata item is a key-value pair of strings.
metadataItem
    : TYPE_METAD_KW COLON typeExpr # typeMetadata
    | NOTES_METAD_KW COLON stringExpr  # notesMetadata
    | DEFAULT_METAD_KW COLON primitive # defaultMetadata
    | OPTIONAL_METAD_KW COLON BOOL # optionalMetadata
    ;

// A type expression denotes the type.
typeExpr
    : IDENTIFIER
    | LIST_TYPE_KW LANGLE typeExpr RANGLE
    ;

// A definition of a custom object type.
objectDefinition
    : IDENTIFIER LBRACE objectPropertyDefinition* RBRACE
    ;

// A definition of a property of a custom object type.
objectPropertyDefinition
    : simpleName shortMetadataExpression OPTIONAL_METAD_KW?
    ;

// A primitive is a string, an integer, a float, or a boolean.
primitive
    : SHORT_STRING # string
    | INT # int
    | FLOAT # float
    | BOOL # boolean
    ;

fieldName: simpleName | dottedName;

simpleName: LITERAL_STRING | IDENTIFIER;

dottedName: simpleName (DOT simpleName)+;

// A string expression is either a short string or a long string.
stringExpr
    : SHORT_STRING 
    | DOUBLE_QUOTES LONG_STRING DOUBLE_QUOTES
    ;

// Keywords
CONFIG_DCLR_KW : 'config' ;     // Config declaration keyword
IMPORT_KW : 'import' ;     // Import keyword
SPEC_ROOT_KW : 'spec' ;     // Specification keyword
OBJ_DEF_KW : 'objects' ;     // Object definition keyword
TYPE_METAD_KW : 'type' ;         // Type keyword
OPTIONAL_METAD_KW : 'optional' ; // Optional keyword
DEFAULT_METAD_KW : 'default' ;   // Default keyword
NOTES_METAD_KW : 'notes' ;       // Notes keyword
LIST_TYPE_KW : 'list' ;         // List keyword

// Common Tokens
LPAREN : '(' ;            // Left parenthesis
RPAREN : ')' ;            // Right parenthesis
LBRACE : '{' ;            // Left curly brace
RBRACE : '}' ;            // Right curly brace
LANGLE : '<' ;            // Less than symbol, used as left angle bracket
RANGLE : '>' ;            // Greater than symbol, used as right angle bracket
SEMICOLON : ';' ;         // Semicolon
COMMA : ',' ;             // Comma
COLON : ':' ;             // Colon
DOT : '.' ;               // Dot
DOUBLE_QUOTES : '""' ;      // Double quote

SHORT_STRING: '"'  ('\\' (RN | .) | ~[\\\r\n"])* '"';
LITERAL_STRING : '\'' (~['\n])*? '\'' ;
LONG_STRING: '"' LONG_STRING_ITEM*? '"';
INT : DIGIT+ ;               // Integer numbers
FLOAT : DIGIT+ '.' DIGIT+ ;  // Floating point numbers
BOOL : 'true' | 'false' ;    // Boolean values

// CMCL Tokens
IF_SYM: 'if';
ELSEIF_SYM: 'elseif';
ELSE_SYM: 'else';
FOREACH_SYM: 'foreach';
AND_SYM: '&&';
OR_SYM: '||';
NOT_SYM: '!';

IDENTIFIER : (CHARACTER)+ ;    // Typical definition of an identifier

WS : [ \t\r\n]+ -> skip ;    // Skip whitespace

// Auxiliary lexer rules
fragment LETTER : [a-zA-Z] ;
fragment DIGIT : [0-9] ;
fragment CHARACTER : [a-zA-Z0-9_-] ;
fragment LONG_STRING_ITEM
    : ~'\\'
    | '\\' (RN | .)
    ;
fragment RN : '\r'? '\n';