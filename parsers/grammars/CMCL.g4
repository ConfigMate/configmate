// This grammar describes the syntax of the language
// of checks in ConfigMate. CMCL stands for ConfigMate Check Language.

grammar CMCL;

// A check is a list of functions
// separated with a dot
check: function (DOT function)*;

// A function is a name followed by a list of arguments
function
    : NAME LPAREN argument (COMMA argument)* RPAREN
    | NAME LPAREN RPAREN;

// An argument is either a primitive value, a function, or a field
argument: primitive | field;

// A primitive value is a string, an int, a float, or a boolean
primitive
    : STRING # string
    | INT # int
    | FLOAT # float
    | BOOL # boolean
    ;

// A field is a field name followed by an optional check
field: fieldname (DOT check)?;

// A field name is a list of comma separated unquoted strings or index expressions
// Such as "foo", "bar.xyz", "baz[0]", "[1].xyz".
fieldname: (NAME | LBRACKET INT RBRACKET) (DOT (NAME | LBRACKET INT RBRACKET))*;

// Tokens
BOOL: 'true' | 'false';
NAME: [a-zA-Z_][a-zA-Z0-9_]*;
INT: [0-9]+;
FLOAT: [0-9]+[.][0-9]+;
STRING: '"' (ESC | ~["\\])* '"';
LPAREN: '(';
RPAREN: ')';
LBRACKET: '[';
RBRACKET: ']';
DOT: '.';
COMMA: ',';
ESC: '\\' [btnfr"'\\];
WS: [ \t\r\n]+ -> skip;




