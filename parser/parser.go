package parser

import (
	"emerald/ast"
	"emerald/lexer"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
	prefixFns map[lexer.TokenType]prefixParseFn
	infixFns  map[lexer.TokenType]infixParseFn
}

const (
	_ int = iota
	LOWEST
	EQUALS      // =
	OR          // or
	AND         // and
	LESSGREATER // < > <= >=
	SUM         // + -
	PRODUCT     // * /
	CALL        // ()
	PREFIX
)

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:         l,
		prefixFns: make(map[lexer.TokenType]prefixParseFn),
		infixFns:  make(map[lexer.TokenType]infixParseFn),
	}

	p.registerPrefix(lexer.IDENTIFIER, p.parseIdentifierOrCall)
	p.registerPrefix(lexer.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.DOLLAR_STRING, p.parseInterpolatedStringLiteral)
	p.registerPrefix(lexer.INPUT, p.parseInputExpression)
	p.registerPrefix(lexer.BOOLEAN, p.parseBooleanLiteral)
	p.registerPrefix(lexer.NULL, p.parseNullLiteral)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.NOT, p.parsePrefixExpression)
	p.registerPrefix(lexer.NUM, p.parseTypeOrCall)
	p.registerPrefix(lexer.STR, p.parseTypeOrCall)
	p.registerPrefix(lexer.BOOL, p.parseTypeOrCall)

	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.ASTERISK, p.parseInfixExpression)
	p.registerInfix(lexer.SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LE, p.parseInfixExpression)
	p.registerInfix(lexer.GE, p.parseInfixExpression)
	p.registerInfix(lexer.OR, p.parseInfixExpression)
	p.registerInfix(lexer.AND, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) registerPrefix(t lexer.TokenType, fn prefixParseFn) {
	p.prefixFns[t] = fn
}

func (p *Parser) registerInfix(t lexer.TokenType, fn infixParseFn) {
	p.infixFns[t] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		for p.curToken.Type == lexer.NEWLINE {
			p.nextToken()
		}
	}

	return program
}

func (p *Parser) peekPrecedence() int {
	if fn, ok := p.infixFns[p.peekToken.Type]; ok {
		_ = fn
		switch p.peekToken.Type {
		case lexer.EQ:
			return EQUALS
		case lexer.OR:
			return OR
		case lexer.AND:
			return AND
		case lexer.LT, lexer.GT, lexer.LE, lexer.GE:
			return LESSGREATER
		case lexer.PLUS, lexer.MINUS:
			return SUM
		case lexer.ASTERISK, lexer.SLASH:
			return PRODUCT
		default:
			return LOWEST
		}
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	switch p.curToken.Type {
	case lexer.EQ:
		return EQUALS
	case lexer.OR:
		return OR
	case lexer.AND:
		return AND
	case lexer.LT, lexer.GT, lexer.LE, lexer.GE:
		return LESSGREATER
	case lexer.PLUS, lexer.MINUS:
		return SUM
	case lexer.ASTERISK, lexer.SLASH:
		return PRODUCT
	default:
		return LOWEST
	}
}

func (p *Parser) error(msg string) {
	p.errors = append(p.errors, msg)
}
