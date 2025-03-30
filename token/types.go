package token

// TokenType represents a type of lexical token in the Sludge programming language.
// It is implemented as a string to provide readable token names.
type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL" // ILLEGAL represents an invalid or unexpected token

	// Single-character tokens
	LEFT_PAREN    TokenType = "(" // LEFT_PAREN represents a left parenthesis "("
	RIGHT_PAREN   TokenType = ")" // RIGHT_PAREN represents a right parenthesis ")"
	LEFT_BRACKET  TokenType = "[" // LEFT_BRACKET represents a left square bracket "["
	RIGHT_BRACKET TokenType = "]" // RIGHT_BRACKET represents a right square bracket "]"
	LEFT_BRACE    TokenType = "{" // LEFT_BRACE represents a left curly brace "{"
	RIGHT_BRACE   TokenType = "}" // RIGHT_BRACE represents a right curly brace "}"
	COMMA         TokenType = "," // COMMA represents a comma separator ","
	DOT           TokenType = "." // DOT represents a dot operator "."
	SLASH         TokenType = "/" // SLASH represents a division operator "/"
	STAR          TokenType = "*" // STAR represents a multiplication operator "*"
	PLUS          TokenType = "+" // PLUS represents an addition operator "+"
	MINUS         TokenType = "-" // MINUS represents a subtraction operator "-"
	SEMICOLON     TokenType = ";" // SEMICOLON represents a statement terminator ";"
	COLON         TokenType = ":" // COLON represents a colon operator ":"
	PERCENT       TokenType = "%" // PERCENT represents a modulo operator "%"

	// One or two character tokens
	BANG          TokenType = "!"  // BANG represents a logical NOT operator "!"
	EQUAL         TokenType = "="  // EQUAL represents an assignment operator "="
	LESS          TokenType = "<"  // LESS represents a less than operator "<"
	GREATER       TokenType = ">"  // GREATER represents a greater than operator ">"
	BANG_EQUAL    TokenType = "!=" // BANG_EQUAL represents an inequality operator "!="
	EQUAL_EQUAL   TokenType = "==" // EQUAL_EQUAL represents an equality operator "=="
	LESS_EQUAL    TokenType = "<=" // LESS_EQUAL represents a less than or equal operator "<="
	GREATER_EQUAL TokenType = ">=" // GREATER_EQUAL represents a greater than or equal operator ">="
	AND           TokenType = "&&" // AND represents a logical AND operator "&&"
	OR            TokenType = "||" // OR represents a logical OR operator "||"
	ARROW         TokenType = "=>" // ARROW represents an arrow operator "=>"
	PLUS_EQUAL    TokenType = "+=" // PLUS_EQUAL represents an addition assignment operator "+="
	MINUS_EQUAL   TokenType = "-=" // MINUS_EQUAL represents a subtraction assignment operator "-="
	STAR_EQUAL    TokenType = "*=" // STAR_EQUAL represents a multiplication assignment operator "*="
	SLASH_EQUAL   TokenType = "/=" // SLASH_EQUAL represents a division assignment operator "/="
	PERCENT_EQUAL TokenType = "%=" // PERCENT_EQUAL represents a modulo assignment operator "%="

	// Literals
	IDENTIFIER    TokenType = "IDENTIFIER"    // IDENTIFIER represents a variable or function name
	STRING        TokenType = "STRING"        // STRING represents a string literal
	RAW_STRING    TokenType = "RAW_STRING"    // RAW_STRING represents a raw string literal
	INTEGER       TokenType = "INTEGER"       // INTEGER represents an integer literal
	FLOAT         TokenType = "FLOAT"         // FLOAT represents a floating-point literal
	INTERPOLATION TokenType = "INTERPOLATION" // INTERPOLATION represents string interpolation "${}"
	TEMPLATE      TokenType = "TEMPLATE"      // TEMPLATE represents a template expression "@{}"

	// Keywords
	FUNCTION TokenType = "FUNCTION" // FUNCTION represents the "function" keyword
	LET      TokenType = "LET"      // LET represents the "let" keyword for variable declaration
	CONST    TokenType = "CONST"    // CONST represents the "const" keyword for constant declaration
	VAR      TokenType = "VAR"      // VAR represents the "var" keyword for variable declaration
	TRUE     TokenType = "TRUE"     // TRUE represents the boolean literal "true"
	FALSE    TokenType = "FALSE"    // FALSE represents the boolean literal "false"
	IF       TokenType = "IF"       // IF represents the "if" conditional keyword
	ELSE     TokenType = "ELSE"     // ELSE represents the "else" conditional keyword
	WHILE    TokenType = "WHILE"    // WHILE represents the "while" loop keyword
	FOR      TokenType = "FOR"      // FOR represents the "for" loop keyword
	RETURN   TokenType = "RETURN"   // RETURN represents the "return" statement keyword
	BREAK    TokenType = "BREAK"    // BREAK represents the "break" statement keyword
	CONTINUE TokenType = "CONTINUE" // CONTINUE represents the "continue" statement keyword
	NULL     TokenType = "NULL"     // NULL represents the "null" literal
	IMPORT   TokenType = "IMPORT"   // IMPORT represents the "import" keyword

	EOF TokenType = "EOF" // EOF represents the end of file token
)
