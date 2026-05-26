package parser

import (
	"fmt"
	"strconv"

	"emerald/ast"
	"emerald/lexer"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn, ok := p.prefixFns[p.curToken.Type]
	if !ok {
		p.error(fmt.Sprintf("unexpected token '%s' at line %d, col %d", p.curToken.Literal, p.curToken.Line, p.curToken.Col))
		return nil
	}
	left := prefixFn()

	for p.peekToken.Type != lexer.NEWLINE &&
		p.peekToken.Type != lexer.EOF &&
		precedence < p.peekPrecedence() {
		infixFn, ok := p.infixFns[p.peekToken.Type]
		if !ok {
			return left
		}
		p.nextToken()
		left = infixFn(left)
	}

	return left
}

func (p *Parser) parseIdentifierOrCall() ast.Expression {
	if p.peekToken.Type == lexer.LPAREN {
		return p.parseCallExpression()
	}
	return p.parseIdentifier()
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Value: p.curToken.Literal}
}

func (p *Parser) parseCallExpression() ast.Expression {
	name := p.curToken.Literal
	line := p.curToken.Line
	p.nextToken()
	p.nextToken()
	arg := p.parseExpression(LOWEST)
	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken()
	}
	return &ast.CallExpression{Function: name, Argument: arg, Line: line}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.error(fmt.Sprintf("invalid number '%s' at line %d", p.curToken.Literal, p.curToken.Line))
		return nil
	}
	return &ast.NumberLiteral{Value: value}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	value := p.curToken.Literal == "True"
	return &ast.BooleanLiteral{Value: value}
}

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken()
	}
	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	op := p.curToken.Literal
	line := p.curToken.Line
	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)
	return &ast.BinaryExpression{Left: left, Operator: op, Right: right, Line: line}
}
