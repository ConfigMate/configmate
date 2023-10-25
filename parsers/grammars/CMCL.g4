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

foreach: FOREACH_SYM LPAREN NAME COLON field RPAREN LBRACE check RBRACE;

not: NOT_SYM atom;

functionExpression: function (DOT function)*;

fieldExpression: field (DOT functionExpression)?;

function
    : NAME LPAREN argument (COMMA argument)* RPAREN
    | NAME LPAREN RPAREN;

argument: primitive | check;

primitive
    : STRING # string
    | INT # int
    | FLOAT # float
    | BOOL # boolean
    ;

// A field is a list of comma separated unquoted strings or index expressions
// Such as "foo", "bar.xyz".
field: NAME (DOT NAME)*;

// Tokens
IF_SYM: 'if';
ELSEIF_SYM: 'elseif';
ELSE_SYM: 'else';
FOREACH_SYM: 'foreach';
AND_SYM: '&&';
OR_SYM: '||';
NOT_SYM: '!';
BOOL: 'true' | 'false';
NAME: [a-zA-Z_][a-zA-Z0-9_]*;
INT: [0-9]+;
FLOAT: [0-9]+[.][0-9]+;
STRING: '"' (ESC | ~["\\])* '"' | '\'' (ESC | ~["\\])* '\'';
LPAREN: '(';
RPAREN: ')';
LBRACKET: '[';
RBRACKET: ']';
LBRACE: '{';
RBRACE: '}';
DOT: '.';
COMMA: ',';
COLON: ':';
ESC: '\\' [btnfr"'\\];
WS: [ \t\r\n]+ -> skip;




