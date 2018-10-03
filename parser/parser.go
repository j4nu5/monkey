package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type parselet func(*Parser) ast.Statement
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

var parselets map[token.Type]parselet

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NOTEQ:    EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// Parser parses Monkey.
type Parser struct {
	l *lexer.Lexer

	currentToken token.Token
	nextToken    token.Token

	errors []string

	prefixParsers map[token.Type]prefixParseFn
	infixParsers  map[token.Type]infixParseFn
}

// New creates and initializes a new Parser for Monkey.
func New(l *lexer.Lexer) *Parser {
	parselets = map[token.Type]parselet{
		token.LET:    letStatementParselet,
		token.RETURN: returnStatementParselet,
		token.LBRACE: blockStatementParselet,
	}

	p := &Parser{l: l}
	p.prefixParsers = make(map[token.Type]prefixParseFn)
	p.infixParsers = make(map[token.Type]infixParseFn)
	p.consumeNextToken()
	p.consumeNextToken()

	p.registerPrefixParser(token.IDENT, p.parseIdentifier)
	p.registerPrefixParser(token.INT, p.parseIntegerLiteral)
	p.registerPrefixParser(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefixParser(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefixParser(token.BANG, p.parsePrefixOperator)
	p.registerPrefixParser(token.MINUS, p.parsePrefixOperator)
	p.registerPrefixParser(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixParser(token.IF, p.parseIfExpression)
	p.registerPrefixParser(token.FUNCTION, p.parseFunctionExpression)

	p.registerInfixParser(token.PLUS, p.parseInfixExpression)
	p.registerInfixParser(token.MINUS, p.parseInfixExpression)
	p.registerInfixParser(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixParser(token.SLASH, p.parseInfixExpression)
	p.registerInfixParser(token.LT, p.parseInfixExpression)
	p.registerInfixParser(token.GT, p.parseInfixExpression)
	p.registerInfixParser(token.EQ, p.parseInfixExpression)
	p.registerInfixParser(token.NOTEQ, p.parseInfixExpression)

	return p
}

func (p *Parser) consumeNextToken() {
	if p.currentToken.Type == token.EOF {
		return
	}

	p.currentToken = p.nextToken
	p.nextToken = p.l.NextToken()
}

func (p *Parser) registerPrefixParser(t token.Type, fn prefixParseFn) {
	p.prefixParsers[t] = fn
}

func (p *Parser) registerInfixParser(t token.Type, fn infixParseFn) {
	p.infixParsers[t] = fn
}

// ParseProgram returns a parsed Monkey program.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.currentToken.Type != token.EOF {
		program.Statements = append(program.Statements, p.parseStatement())
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	parselet, ok := parselets[p.currentToken.Type]
	if !ok {
		return expressionStatementParselet(p)
	}
	return parselet(p)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixParser, ok := p.prefixParsers[p.currentToken.Type]
	if !ok {
		p.addError(fmt.Sprintf("Unable to consume prefix: %v", p.currentToken))
		p.consumeNextToken()
		return nil
	}
	left := prefixParser()

	for p.currentToken.Type != token.SEMICOLON && precedence < p.getCurrentPrecedence() {
		infixParser, ok := p.infixParsers[p.currentToken.Type]
		if !ok {
			return left
		}
		left = infixParser(left)
	}

	return left
}

func (p *Parser) getCurrentPrecedence() int {
	precedence, ok := precedences[p.currentToken.Type]
	if !ok {
		return LOWEST
	}

	return precedence
}

func (p *Parser) parseIdentifier() ast.Expression {
	expr := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
	p.consumeNextToken()

	return expr
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	expr := &ast.IntegerLiteral{Token: p.currentToken}
	val, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	p.consumeNextToken()
	if err != nil {
		p.addError(fmt.Sprintf("Could not parse %q as integer", p.currentToken.Literal))
		return nil
	}
	expr.Value = val
	return expr
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	expr := &ast.Boolean{Token: p.currentToken}
	expr.Value = (p.currentToken.Type == token.TRUE)
	p.consumeNextToken()
	return expr
}

func (p *Parser) parsePrefixOperator() ast.Expression {
	expr := &ast.PrefixExpression{Token: p.currentToken, Operator: p.currentToken.Literal}
	p.consumeNextToken()
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.consumeToken(token.LPAREN)
	expr := p.parseExpression(LOWEST)
	p.consumeToken(token.RPAREN)
	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal,
	}
	precedence := p.getCurrentPrecedence()
	p.consumeNextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.currentToken}
	p.consumeToken(token.IF)
	p.consumeToken(token.LPAREN)
	expr.Condition = p.parseExpression(LOWEST)
	p.consumeToken(token.RPAREN)

	expr.Consequence = blockStatementParselet(p).(*ast.BlockStatement)
	if p.currentToken.Type != token.ELSE {
		return expr
	}

	p.consumeToken(token.ELSE)
	expr.Alternative = blockStatementParselet(p).(*ast.BlockStatement)

	return expr
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	expr := &ast.FunctionLiteral{Token: p.currentToken}
	expr.Parameters = []*ast.Identifier{}

	p.consumeToken(token.FUNCTION)
	p.consumeToken(token.LPAREN)
	for p.currentToken.Type != token.EOF && p.currentToken.Type != token.RPAREN {
		expr.Parameters = append(expr.Parameters, p.parseIdentifier().(*ast.Identifier))
		if p.currentToken.Type == token.COMMA {
			p.consumeToken(token.COMMA)
		}
	}
	p.consumeToken(token.RPAREN)

	expr.Body = blockStatementParselet(p).(*ast.BlockStatement)

	return expr
}

func (p *Parser) consumeToken(t token.Type) {
	if p.currentToken.Type != t {
		p.addError(
			fmt.Sprintf("Could not properly consume token: %v expected: %v", p.currentToken, t))
	}
	p.consumeNextToken()
}

// Errors returns all parsing errors seen so far.
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(err string) {
	p.errors = append(p.errors, err)
}

// Parselets

func letStatementParselet(p *Parser) ast.Statement {
	if p.currentToken.Type != token.LET {
		return nil
	}

	stmt := &ast.LetStatement{Token: p.currentToken}
	p.consumeToken(token.LET)

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	p.consumeToken(token.IDENT)
	p.consumeToken(token.ASSIGN)

	stmt.Value = p.parseExpression(LOWEST)
	p.consumeToken(token.SEMICOLON)

	return stmt
}

func returnStatementParselet(p *Parser) ast.Statement {
	if p.currentToken.Type != token.RETURN {
		return nil
	}

	stmt := &ast.ReturnStatement{Token: p.currentToken}
	p.consumeToken(token.RETURN)

	stmt.ReturnValue = p.parseExpression(LOWEST)
	p.consumeToken(token.SEMICOLON)

	return stmt
}

func expressionStatementParselet(p *Parser) ast.Statement {
	stmt := &ast.ExpressionStatement{
		Token:      p.currentToken,
		Expression: p.parseExpression(LOWEST),
	}
	if p.currentToken.Type == token.SEMICOLON {
		p.consumeToken(token.SEMICOLON)
	}

	return stmt
}

func blockStatementParselet(p *Parser) ast.Statement {
	if p.currentToken.Type != token.LBRACE {
		return nil
	}

	stmt := &ast.BlockStatement{Token: p.currentToken}
	stmt.Statements = []ast.Statement{}
	p.consumeToken(token.LBRACE)
	for p.currentToken.Type != token.EOF && p.currentToken.Type != token.RBRACE {
		stmt.Statements = append(stmt.Statements, p.parseStatement())
	}
	p.consumeToken(token.RBRACE)

	return stmt
}
