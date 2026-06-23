package lexer

import "fmt"

type TokenType string

const (
	EOF        TokenType = "EOF"
	ILLEGAL    TokenType = "ILLEGAL"
	NEWLINE    TokenType = "NEWLINE"
	VAR        TokenType = "VAR"
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING          TokenType = "STRING"
	DOLLAR_STRING   TokenType = "DOLLAR_STRING"
	NUMBER     TokenType = "NUMBER"
	BOOLEAN    TokenType = "BOOLEAN"
	NULL       TokenType = "NULL"
	PRINT      TokenType = "PRINT"
	PLUS       TokenType = "PLUS"
	MINUS      TokenType = "MINUS"
	ASTERISK   TokenType = "ASTERISK"
	SLASH      TokenType = "SLASH"
	LPAREN     TokenType = "LPAREN"
	RPAREN     TokenType = "RPAREN"
	EQ         TokenType = "EQ"
	LT         TokenType = "LT"
	GT         TokenType = "GT"
	LE         TokenType = "LE"
	GE         TokenType = "GE"
	IF         TokenType = "IF"
	ELIF       TokenType = "ELIF"
	ELSE       TokenType = "ELSE"
	FUNC       TokenType = "FUNC"
	RUN        TokenType = "RUN"
	LBRACE     TokenType = "LBRACE"
	RBRACE     TokenType = "RBRACE"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%s, %q, Ln:%d, Col:%d}", t.Type, t.Literal, t.Line, t.Col)
}
