package lexer

import "monkey/token"

type lexlet func(l *Lexer) (token.Token, bool)

var lexlets = []lexlet{
	eofLexlet,
	twoCharLexlet,
	singleCharLexlet,
	literalsLexlet,
	identifiersAndKeywordsLexlet,
}

// Lexer does a lexical analysis for Monkey.
type Lexer struct {
	input        string
	position     int
	nextPosition int
	ch           byte
}

// New creates and initializes a new Lexer.
// input is the source code in ASCII.
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.consumeNextChar()

	return l
}

func (l *Lexer) consumeNextChar() {
	l.position = l.nextPosition
	if l.nextPosition < len(l.input) {
		l.nextPosition++
	}

	if l.position < len(l.input) {
		l.ch = l.input[l.position]
	} else {
		l.ch = 0
	}
}

func (l *Lexer) peekNextChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}

	return l.input[l.nextPosition]
}

func (l *Lexer) consumeWhitespaces() {
	for l.ch != 0 && isWhitespace(l.ch) {
		l.consumeNextChar()
	}
}

// NextToken returns the next token in the source.
func (l *Lexer) NextToken() token.Token {
	l.consumeWhitespaces()
	for _, lexlet := range lexlets {
		if tok, ok := lexlet(l); ok {
			return tok
		}
	}

	ch := l.ch
	l.consumeNextChar()
	return token.Token{
		Type:    token.ILLEGAL,
		Literal: string(ch),
	}
}

func eofLexlet(l *Lexer) (token.Token, bool) {
	if l.ch != 0 {
		return token.Token{}, false
	}

	return token.Token{Type: token.EOF}, true
}

func twoCharLexlet(l *Lexer) (token.Token, bool) {
	chars := string(l.ch) + string(l.peekNextChar())
	typ := token.GetTwoCharToken(chars)
	if typ == token.ILLEGAL {
		return token.Token{}, false
	}

	l.consumeNextChar()
	l.consumeNextChar()
	return token.Token{
		Type:    typ,
		Literal: chars,
	}, true
}

func singleCharLexlet(l *Lexer) (token.Token, bool) {
	typ := token.GetSingleCharToken(l.ch)
	if typ == token.ILLEGAL {
		return token.Token{}, false
	}

	ch := l.ch
	l.consumeNextChar()
	return token.Token{
		Type:    typ,
		Literal: string(ch),
	}, true
}

func literalsLexlet(l *Lexer) (token.Token, bool) {
	if l.ch == 0 || !isDigit(l.ch) {
		return token.Token{}, false
	}

	pos := l.position
	for isDigit(l.ch) {
		l.consumeNextChar()
	}

	literal := l.input[pos:l.position]
	return token.Token{
		Type:    token.INT,
		Literal: literal,
	}, true
}

func identifiersAndKeywordsLexlet(l *Lexer) (token.Token, bool) {
	if l.ch == 0 || !isCharacter(l.ch) {
		return token.Token{}, false
	}

	pos := l.position
	for isCharacter(l.ch) {
		l.consumeNextChar()
	}

	ident := l.input[pos:l.position]

	return token.Token{
		Type:    token.GetTokenType(ident),
		Literal: ident,
	}, true
}

func isCharacter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
