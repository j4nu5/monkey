package token

// Type is the type of a token.
type Type string

// Token represents a token in Monkey.
type Token struct {
	Type    Type
	Literal string
}

const (
	// ILLEGAL represents an illegal token.
	ILLEGAL = "ILLEGAL"
	// EOF represents the end of source.
	EOF = "EOF"

	// IDENT represents an identifier.
	IDENT = "IDENT"

	// INT represents an integer literal.
	INT = "INT"

	// ASSIGN is the assignment operator.
	ASSIGN = "="
	// PLUS is the '+' operator.
	PLUS = "+"
	// MINUS is the '-' operator.
	MINUS = "-"
	// SLASH is the '/' operator.
	SLASH = "/"
	// ASTERISK is the '*' operator.
	ASTERISK = "*"
	// BANG is the '!' operator.
	BANG = "!"
	// LT is the '<' operator.
	LT = "<"
	// GT is the '>' operator.
	GT = ">"
	// EQ is the '==' operator.
	EQ = "=="
	// NOTEQ is the '!=' operator.
	NOTEQ = "!="

	// COMMA is the ',' delimiter.
	COMMA = ","
	// SEMICOLON is the ';' delimiter.
	SEMICOLON = ";"

	// LPAREN is the left parenthesis.
	LPAREN = "("
	// RPAREN is the right parenthesis.
	RPAREN = ")"
	// LBRACE is '{'
	LBRACE = "{"
	// RBRACE is '}'.
	RBRACE = "}"

	// TRUE is the 'true' keyword
	TRUE = "TRUE"
	// FALSE is the 'false' keyword
	FALSE = "FALSE"
	// FUNCTION is the 'fn' keyword.
	FUNCTION = "FUNCTION"
	// LET is the 'let' keyword.
	LET = "LET"
	// IF is the 'if' keyword.
	IF = "IF"
	// ELSE is the 'else' keyword.
	ELSE = "ELSE"
	// RETURN is the 'return' keyword.
	RETURN = "RETURN"
)

var keywords = map[string]Type{
	"let":    LET,
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

var singleCharTokens = map[byte]Type{
	',': COMMA,
	';': SEMICOLON,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	'=': ASSIGN,
	'+': PLUS,
	'-': MINUS,
	'*': ASTERISK,
	'/': SLASH,
	'!': BANG,
	'<': LT,
	'>': GT,
}

var twoCharTokens = map[string]Type{
	"==": EQ,
	"!=": NOTEQ,
}

// GetTokenType returns the keyword type for |identifier|, if it exists.
// Otherwise, it returns the |IDENT| type.
func GetTokenType(identifier string) Type {
	t, ok := keywords[identifier]

	if !ok {
		return IDENT
	}
	return t
}

// GetSingleCharToken returns the token type for ch, if it exists. Returns ILLEGAL otherwise.
func GetSingleCharToken(ch byte) Type {
	t, ok := singleCharTokens[ch]

	if !ok {
		return ILLEGAL
	}
	return t
}

// GetTwoCharToken returns the token type for |s|, if it exists. Returns ILLEGAL otherwise.
func GetTwoCharToken(s string) Type {
	t, ok := twoCharTokens[s]

	if !ok {
		return ILLEGAL
	}
	return t
}
