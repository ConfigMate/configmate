// This grammar describes the syntax of the language
// of specifications in ConfigMate. CMSL stands for ConfigMate Specification Language.
grammar CMSL;

import CMCL;

// The top-level rule of the grammar.
// A CMSL specification contains a file declaration, a list of imports,
// a specification body, and an optional list of custom object types.
specification: fileDeclaration importStatement? specificationBody;

// A file declaration contains the path and format of the file.
fileDeclaration: 'file' COLON FILEPATH COLON IDENTIFIER;

// An import contains the name of the file to import.
importStatement: 'import' LPAREN importItem (COMMA importItem)* RPAREN;

// An import item contains the name of the file to import.
importItem: IDENTIFIER COLON FILEPATH;

// A specification body contains a list of declarations.
specificationBody: 'spec' LBRACE specificationItem* RBRACE;

// A specification item starts with the field name, followed by the
// metadata inside angled brackets, optionally followed by a list of semicolon separated
// checks (CMCL expressions), and optionally followed with the specification of underlying
// fields insided curly braces.
specificationItem: fieldName metadataExpression ((checkStatement SEMICOLON)+)? (LBRACE specificationItem* RBRACE)?;

// A check statement is a CMCL expression.
checkStatement: check;

// A field name is a string literal.
fieldName: IDENTIFIER (DOT IDENTIFIER)*;

// A metadata expression is a list of metadata items inside angled brackets.
metadataExpression: LANGLE metadataItem (COMMA metadataItem)* RANGLE;

// A metadata item is a key-value pair of strings.
metadataItem: IDENTIFIER COLON primitive;

// A primitive is a string, an integer, a float, or a boolean.
primitive
    : STRING # string
    | INT # int
    | FLOAT # float
    | BOOL # boolean
    ;

// Tokens
// Define your lexer tokens here

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

FILEPATH : ('/' | '\\')? (('.' | '..') '/')* (CHARACTER | '.' | '\\')+
    (('/' | '\\') (CHARACTER | '.' | '\\')+)* ;  // Define this as per your file path requirements

IDENTIFIER : LETTER (LETTER | DIGIT | '_')* ;    // Typical definition of an identifier

STRING : '"' (ESC_SEQ | ~["\\])* '"' ;           // String literals

INT : DIGIT+ ;               // Integer numbers
FLOAT : DIGIT+ '.' DIGIT+ ;  // Floating point numbers
BOOL : 'true' | 'false' ;    // Boolean values

// Auxiliary lexer rules
fragment LETTER : [a-zA-Z] ;
fragment DIGIT : [0-9] ;
fragment CHARACTER : [a-zA-Z0-9_-] ;
fragment ESC_SEQ : '\\' ('b' | 't' | 'n' | 'f' | 'r' | 'u' | '"' | '\'' | '\\') ;  // Escape sequences




