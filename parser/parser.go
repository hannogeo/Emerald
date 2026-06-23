package parser

import (
	"fmt"

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
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	CALL
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
	p.registerPrefix(lexer.BOOLEAN, p.parseBooleanLiteral)
	p.registerPrefix(lexer.NULL, p.parseNullLiteral)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)

	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.ASTERISK, p.parseInfixExpression)
	p.registerInfix(lexer.SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LE, p.parseInfixExpression)
	p.registerInfix(lexer.GE, p.parseInfixExpression)

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

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.VAR:
		return p.parseVarStatement()
	case lexer.PRINT:
		return p.parsePrintStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FUNC:
		return p.parseFuncStatement()
	case lexer.RUN:
		return p.parseRunStatement()
	case lexer.ADD:
		return p.parseAddStatement()
	default:
		if p.curToken.Type == lexer.NEWLINE {
			return nil
		}
		if p.curToken.Type == lexer.IDENTIFIER && p.peekToken.Type == lexer.INPUT {
			return p.parseInputStatement()
		}
		p.error(fmt.Sprintf("unexpected token '%s' at line %d, col %d", p.curToken.Literal, p.curToken.Line, p.curToken.Col))
		p.nextToken()
		return nil
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{}
	p.nextToken()

	if p.curToken.Type != lexer.IDENTIFIER {
		p.error(fmt.Sprintf("expected variable name after 'var.' at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
		return nil
	}

	stmt.Name = p.curToken.Literal
	p.nextToken()

	expr := p.parseExpression(LOWEST)
	if expr == nil {
		return nil
	}
	stmt.Value = expr

	for p.curToken.Type != lexer.NEWLINE && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	stmt := &ast.PrintStatement{}
	p.nextToken()

	expr := p.parseExpression(LOWEST)
	if expr == nil {
		return nil
	}
	stmt.Value = expr

	for p.curToken.Type != lexer.NEWLINE && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	for p.curToken.Type != lexer.LBRACE && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}

	stmt.Consequence = p.parseBlockStatement()

	for p.curToken.Type == lexer.NEWLINE {
		p.nextToken()
	}

	if p.curToken.Type == lexer.ELIF {
		stmt.Alternative = p.parseIfStatement()
	} else if p.curToken.Type == lexer.ELSE {
		p.nextToken()
		for p.curToken.Type != lexer.LBRACE && p.curToken.Type != lexer.EOF {
			p.nextToken()
		}
		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseFuncStatement() *ast.FuncStatement {
	stmt := &ast.FuncStatement{}
	p.nextToken()

	if p.curToken.Type != lexer.IDENTIFIER {
		p.error(fmt.Sprintf("expected function name after 'func.' at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
		return nil
	}

	stmt.Name = p.curToken.Literal
	p.nextToken()

	for p.curToken.Type != lexer.LBRACE && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseRunStatement() *ast.RunStatement {
	stmt := &ast.RunStatement{}
	p.nextToken()

	if p.curToken.Type != lexer.IDENTIFIER {
		p.error(fmt.Sprintf("expected function name after 'run.' at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
		return nil
	}

	stmt.Name = p.curToken.Literal
	p.nextToken()

	for p.curToken.Type != lexer.NEWLINE && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAddStatement() *ast.AddStatement {
	stmt := &ast.AddStatement{}
	p.nextToken()

	if p.curToken.Type != lexer.IDENTIFIER {
		p.error(fmt.Sprintf("expected list name after 'add.' at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
		return nil
	}

	stmt.Name = p.curToken.Literal
	p.nextToken()

	expr := p.parseExpression(LOWEST)
	if expr == nil {
		return nil
	}
	stmt.Value = expr

	for p.curToken.Type != lexer.NEWLINE && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseInputStatement() *ast.InputStatement {
	stmt := &ast.InputStatement{}
	stmt.Name = p.curToken.Literal
	p.nextToken()
	p.nextToken()

	if p.curToken.Type != lexer.STRING {
		p.error(fmt.Sprintf("expected prompt string after 'input.' at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
		return nil
	}
	stmt.Prompt = p.curToken.Literal
	p.nextToken()

	for p.curToken.Type != lexer.NEWLINE && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	p.nextToken()

	for p.curToken.Type == lexer.NEWLINE {
		p.nextToken()
	}

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		for p.curToken.Type == lexer.NEWLINE {
			p.nextToken()
		}
	}

	if p.curToken.Type == lexer.RBRACE {
		p.nextToken()
	}

	return block
}

func (p *Parser) peekPrecedence() int {
	if fn, ok := p.infixFns[p.peekToken.Type]; ok {
		_ = fn
		switch p.peekToken.Type {
		case lexer.EQ:
			return EQUALS
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
