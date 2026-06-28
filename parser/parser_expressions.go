package parser

import (
	"fmt"
	"strconv"
	"strings"

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

func (p *Parser) parseColonInfix(left ast.Expression) ast.Expression {
	line := p.curToken.Line
	p.nextToken()

	if p.curToken.Type == lexer.NUMBER {
		lit := p.parseNumberLiteral()
		ident, ok := left.(*ast.Identifier)
		if !ok {
			p.error(fmt.Sprintf("cannot index non-identifier at line %d", line))
			return nil
		}
		return &ast.ListIndexExpression{
			Name:  ident.Value,
			Index: lit,
			Line:  line,
		}
	}

	if p.curToken.Type == lexer.IDENTIFIER {
		methodName := p.curToken.Literal
		p.nextToken()
		args := []ast.Expression{}
		if p.curToken.Type == lexer.LPAREN {
			p.nextToken()
			if p.curToken.Type != lexer.RPAREN {
				args = append(args, p.parseExpression(LOWEST))
				for p.peekToken.Type == lexer.COMMA {
					p.nextToken()
					p.nextToken()
					args = append(args, p.parseExpression(LOWEST))
				}
			}
			if p.curToken.Type == lexer.RPAREN {
				p.nextToken()
			} else if p.peekToken.Type == lexer.RPAREN {
				p.nextToken()
			}
		}
		return &ast.MethodCallExpression{
			Object: left,
			Method: methodName,
			Args:   args,
			Line:   line,
		}
	}

	p.error(fmt.Sprintf("unexpected token '%s' after ':' at line %d", p.curToken.Literal, line))
	return nil
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

func (p *Parser) parseInterpolatedStringLiteral() ast.Expression {
	raw := p.curToken.Literal
	parts := []ast.InterpolationPart{}
	for {
		start := strings.Index(raw, "{")
		if start == -1 {
			parts = append(parts, ast.InterpolationPart{Text: raw})
			break
		}
		if start > 0 {
			parts = append(parts, ast.InterpolationPart{Text: raw[:start]})
		}
		end := strings.Index(raw[start:], "}")
		if end == -1 {
			p.error(fmt.Sprintf("unclosed '{' in interpolated string at line %d", p.curToken.Line))
			return nil
		}
		exprStr := raw[start+1 : start+end]
		subLexer := lexer.NewLexer(exprStr)
		subParser := NewParser(subLexer)
		expr := subParser.parseExpression(LOWEST)
		if len(subParser.errors) > 0 {
			p.error(fmt.Sprintf("invalid expression in interpolated string at line %d: %s", p.curToken.Line, subParser.errors[0]))
			return nil
		}
		if expr == nil {
			p.error(fmt.Sprintf("empty expression in interpolated string at line %d", p.curToken.Line))
			return nil
		}
		parts = append(parts, ast.InterpolationPart{Expr: expr})
		raw = raw[start+end+1:]
	}
	return &ast.InterpolatedStringLiteral{Parts: parts}
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

	if p.curToken.Type == lexer.RPAREN {
		return &ast.ListLiteral{Elements: []ast.Expression{}}
	}

	exp := p.parseExpression(LOWEST)

	if p.peekToken.Type == lexer.COMMA {
		elements := []ast.Expression{exp}
		p.nextToken()
		for {
			p.nextToken()
			if p.curToken.Type == lexer.RPAREN {
				break
			}
			elements = append(elements, p.parseExpression(LOWEST))
			if p.peekToken.Type != lexer.COMMA {
				break
			}
			p.nextToken()
		}
		if p.peekToken.Type == lexer.RPAREN {
			p.nextToken()
		}
		return &ast.ListLiteral{Elements: elements}
	}

	if p.peekToken.Type == lexer.RPAREN {
		p.nextToken()
	}
	return exp
}

func (p *Parser) parseInputExpression() ast.Expression {
	p.nextToken()
	expr := p.parseExpression(LOWEST)
	return &ast.InputExpression{Prompt: expr}
}

func (p *Parser) parseTypeOrCall() ast.Expression {
	if p.peekToken.Type == lexer.LPAREN {
		return p.parseCallExpression()
	}
	return &ast.TypeLiteral{TypeName: p.curToken.Literal}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	op := p.curToken.Literal
	line := p.curToken.Line
	p.nextToken()
	right := p.parseExpression(PREFIX)
	return &ast.PrefixExpression{Operator: op, Right: right, Line: line}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	op := p.curToken.Literal
	line := p.curToken.Line
	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)
	return &ast.BinaryExpression{Left: left, Operator: op, Right: right, Line: line}
}
