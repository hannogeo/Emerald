package parser

import (
	"fmt"
	"strconv"

	"emerald/ast"
	"emerald/lexer"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
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
	default:
		if p.curToken.Type == lexer.NEWLINE {
			return nil
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

	expr := p.parseExpression()
	if expr == nil {
		return nil
	}
	stmt.Value = expr

	return stmt
}

func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	stmt := &ast.PrintStatement{}
	p.nextToken()

	expr := p.parseExpression()
	if expr == nil {
		return nil
	}
	stmt.Value = expr

	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	switch p.curToken.Type {
	case lexer.STRING:
		return p.parseStringLiteral()
	case lexer.NUMBER:
		return p.parseNumberLiteral()
	case lexer.BOOLEAN:
		return p.parseBooleanLiteral()
	case lexer.IDENTIFIER:
		return p.parseIdentifier()
	default:
		p.error(fmt.Sprintf("expected value at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
		p.nextToken()
		return nil
	}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{Value: p.curToken.Literal}
	p.nextToken()
	return lit
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.error(fmt.Sprintf("invalid number '%s' at line %d", p.curToken.Literal, p.curToken.Line))
		p.nextToken()
		return nil
	}
	lit := &ast.NumberLiteral{Value: value}
	p.nextToken()
	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	var value bool
	if p.curToken.Literal == "True" {
		value = true
	} else {
		value = false
	}
	lit := &ast.BooleanLiteral{Value: value}
	p.nextToken()
	return lit
}

func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Value: p.curToken.Literal}
	p.nextToken()
	return ident
}

func (p *Parser) error(msg string) {
	p.errors = append(p.errors, msg)
}
