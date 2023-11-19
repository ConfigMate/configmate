// This grammar describes the syntax of the language
// of checks in ConfigMate. CMCL stands for ConfigMate Check Language.

grammar CMCL;

check
    : expression          # exprCheck
    | if                  # ifCheck
    | foreach             # foreachCheck
    ;

expression
    : andExpression (OR_SYM andExpression)*  # orExpr
    ;

andExpression
    : atom (AND_SYM atom)*  # andExpr
    ;

atom
    : not                      # notExpr
    | functionExpression       # funcExpr
    | fieldExpression          # fieldExpr
    | LPAREN expression RPAREN # parenExpr
    ;

if: IF_SYM LPAREN check RPAREN LBRACE check RBRACE (elseif)* (else)?;

elseif: ELSEIF_SYM LPAREN check RPAREN LBRACE check RBRACE;

else: ELSE_SYM LBRACE check RBRACE;

foreach: FOREACH_SYM LPAREN IDENTIFIER COLON fieldName RPAREN LBRACE check RBRACE;

not: NOT_SYM atom;

functionExpression: function (DOT function)*;

fieldExpression: fieldName (DOT functionExpression)?;

function
    : IDENTIFIER LPAREN argument (COMMA argument)* RPAREN
    | IDENTIFIER LPAREN RPAREN;

argument: primitive | check;

// A primitive is a string, an integer, a float, or a boolean.
primitive
    : SHORT_STRING # string
    | INT # int
    | FLOAT # float
    | BOOL # boolean
    ;

// A field is a list of dot separated identifiers.
// Such as "foo", "bar.xyz".
fieldName: IDENTIFIER (DOT IDENTIFIER)*;

// Common Tokens
LPAREN : '(' ;            // Left parenthesis
RPAREN : ')' ;            // Right parenthesis
LBRACE : '{' ;            // Left curly brace
RBRACE : '}' ;            // Right curly brace
COMMA : ',' ;             // Comma
COLON : ':' ;             // Colon
DOT : '.' ;               // Dot

SHORT_STRING: '"'  ('\\' (RN | .) | ~[\\\r\n"])* '"';
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

IDENTIFIER : LETTER (CHARACTER)* ;    // Typical definition of an identifier

// Auxiliary lexer rules
fragment LETTER : [a-zA-Z] ;
fragment DIGIT : [0-9] ;
fragment CHARACTER : [a-zA-Z0-9_-] ;
fragment RN : '\r'? '\n';

WS: [ \t\r\n]+ -> skip;




