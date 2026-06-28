package parser

import (
	"fmt"

	"emerald/ast"
	"emerald/lexer"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.VAR:
		return p.parseVarStatement()
	case lexer.PRINT:
		return p.parsePrintStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.FUNC:
		return p.parseFuncStatement()
	case lexer.RUN:
		return p.parseRunStatement()
	case lexer.ADD:
		return p.parseAddStatement()
	case lexer.DOT:
		return p.parseReassignStatement()
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

	compoundOp := ""
	switch p.curToken.Type {
	case lexer.PLUS_EQ:
		compoundOp = "+"
	case lexer.MINUS_EQ:
		compoundOp = "-"
	case lexer.ASTERISK_EQ:
		compoundOp = "*"
	case lexer.SLASH_EQ:
		compoundOp = "/"
	}

	if compoundOp != "" {
		p.nextToken()
		rhs := p.parseExpression(LOWEST)
		if rhs == nil {
			return nil
		}
		stmt.Value = &ast.BinaryExpression{
			Left:     &ast.Identifier{Value: stmt.Name},
			Operator: compoundOp,
			Right:    rhs,
		}
		for p.curToken.Type != lexer.NEWLINE && p.curToken.Type != lexer.EOF {
			p.nextToken()
		}
		return stmt
	}

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

func (p *Parser) parseReassignStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{}
	p.nextToken()

	if p.curToken.Type != lexer.IDENTIFIER {
		p.error(fmt.Sprintf("expected variable name after '.' at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
		return nil
	}

	stmt.Name = p.curToken.Literal
	p.nextToken()

	compoundOp := ""
	switch p.curToken.Type {
	case lexer.PLUS_EQ:
		compoundOp = "+"
	case lexer.MINUS_EQ:
		compoundOp = "-"
	case lexer.ASTERISK_EQ:
		compoundOp = "*"
	case lexer.SLASH_EQ:
		compoundOp = "/"
	}

	if compoundOp != "" {
		p.nextToken()
		rhs := p.parseExpression(LOWEST)
		if rhs == nil {
			return nil
		}
		stmt.Value = &ast.BinaryExpression{
			Left:     &ast.Identifier{Value: stmt.Name},
			Operator: compoundOp,
			Right:    rhs,
		}
		for p.curToken.Type != lexer.NEWLINE && p.curToken.Type != lexer.EOF {
			p.nextToken()
		}
		return stmt
	}

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
		p.error(fmt.Sprintf("expected function name after 'fn.' at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
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
		p.error(fmt.Sprintf("expected function name after 'run' at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
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
		p.error(fmt.Sprintf("expected list name after 'add' at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
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

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{}
	p.nextToken()

	if p.curToken.Type != lexer.IDENTIFIER {
		p.error(fmt.Sprintf("expected loop variable name after 'for' at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
		return nil
	}
	stmt.Variable = p.curToken.Literal
	p.nextToken()

	if p.curToken.Type != lexer.IN {
		p.error(fmt.Sprintf("expected 'in' after loop variable at line %d, got '%s'", p.curToken.Line, p.curToken.Literal))
		return nil
	}
	p.nextToken()

	if p.curToken.Type == lexer.RANGE {
		p.nextToken()
		if p.curToken.Type == lexer.LPAREN {
			p.nextToken()
			start := p.parseExpression(LOWEST)
			if p.peekToken.Type != lexer.COMMA {
				p.error(fmt.Sprintf("expected ',' in range at line %d", p.curToken.Line))
				return nil
			}
			p.nextToken()
			p.nextToken()
			end := p.parseExpression(LOWEST)
			if p.peekToken.Type == lexer.RPAREN {
				p.nextToken()
			}
			stmt.Iterable = &ast.RangeExpression{Start: start, End: end}
		} else {
			end := p.parseExpression(LOWEST)
			stmt.Iterable = &ast.RangeExpression{Start: nil, End: end}
		}
	} else {
		stmt.Iterable = p.parseExpression(LOWEST)
	}

	for p.curToken.Type != lexer.LBRACE && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{}

	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	for p.curToken.Type != lexer.LBRACE && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}

	stmt.Body = p.parseBlockStatement()
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
