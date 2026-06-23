package lexer

import (
	"strings"
)

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
	line    int
	col     int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, col: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.ch == '\n' {
		l.line++
		l.col = 0
	}
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.col++
}

func (l *Lexer) NextToken() Token {
	for {
		l.skipWhitespace()

		if l.ch == '!' {
			l.skipComment()
			continue
		}

		break
	}

	var tok Token

	switch {
	case l.ch == 0:
		tok = Token{Type: EOF, Literal: "", Line: l.line, Col: l.col}

	case l.ch == '\n':
		tok = Token{Type: NEWLINE, Literal: "\\n", Line: l.line, Col: l.col}
		l.readChar()

	case l.ch == '$' && l.peekChar() == '"':
		tok = l.readInterpolatedString()

	case l.ch == '"':
		tok = l.readString()

	case isDigit(l.ch) || (l.ch == '.' && l.readPos < len(l.input) && isDigit(l.input[l.readPos])):
		tok = l.readNumber()

	case isLetter(l.ch):
		tok = l.readIdentifierOrKeyword()

	case l.ch == '+':
		tok = Token{Type: PLUS, Literal: "+", Line: l.line, Col: l.col}
		l.readChar()

	case l.ch == '-':
		tok = Token{Type: MINUS, Literal: "-", Line: l.line, Col: l.col}
		l.readChar()

	case l.ch == '*':
		tok = Token{Type: ASTERISK, Literal: "*", Line: l.line, Col: l.col}
		l.readChar()

	case l.ch == '/':
		tok = Token{Type: SLASH, Literal: "/", Line: l.line, Col: l.col}
		l.readChar()

	case l.ch == '(':
		tok = Token{Type: LPAREN, Literal: "(", Line: l.line, Col: l.col}
		l.readChar()

	case l.ch == ')':
		tok = Token{Type: RPAREN, Literal: ")", Line: l.line, Col: l.col}
		l.readChar()

	case l.ch == '=':
		tok = Token{Type: EQ, Literal: "=", Line: l.line, Col: l.col}
		l.readChar()

	case l.ch == '<':
		if l.peekChar() == '=' {
			tok = Token{Type: LE, Literal: "<=", Line: l.line, Col: l.col}
			l.readChar()
			l.readChar()
		} else {
			tok = Token{Type: LT, Literal: "<", Line: l.line, Col: l.col}
			l.readChar()
		}

	case l.ch == '>':
		if l.peekChar() == '=' {
			tok = Token{Type: GE, Literal: ">=", Line: l.line, Col: l.col}
			l.readChar()
			l.readChar()
		} else {
			tok = Token{Type: GT, Literal: ">", Line: l.line, Col: l.col}
			l.readChar()
		}

	case l.ch == ',':
		tok = Token{Type: COMMA, Literal: ",", Line: l.line, Col: l.col}
		l.readChar()

	case l.ch == '{':
		tok = Token{Type: LBRACE, Literal: "{", Line: l.line, Col: l.col}
		l.readChar()

	case l.ch == '}':
		tok = Token{Type: RBRACE, Literal: "}", Line: l.line, Col: l.col}
		l.readChar()

	default:
		tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: l.line, Col: l.col}
		l.readChar()
	}

	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	l.readChar()
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) readString() Token {
	tok := Token{Type: STRING, Line: l.line, Col: l.col}
	l.readChar()
	var s strings.Builder
	for l.ch != '"' && l.ch != '\n' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				s.WriteByte('\n')
			case 't':
				s.WriteByte('\t')
			case '\\':
				s.WriteByte('\\')
			case '"':
				s.WriteByte('"')
			default:
				s.WriteByte('\\')
				s.WriteByte(l.ch)
			}
		} else {
			s.WriteByte(l.ch)
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar()
	} else {
		return Token{Type: ILLEGAL, Literal: "unterminated string", Line: l.line, Col: l.col}
	}
	tok.Literal = s.String()
	return tok
}

func (l *Lexer) readInterpolatedString() Token {
	tok := Token{Type: DOLLAR_STRING, Line: l.line, Col: l.col}
	l.readChar()
	l.readChar()
	var s strings.Builder
	for l.ch != '"' && l.ch != '\n' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				s.WriteByte('\n')
			case 't':
				s.WriteByte('\t')
			case '\\':
				s.WriteByte('\\')
			case '"':
				s.WriteByte('"')
			default:
				s.WriteByte('\\')
				s.WriteByte(l.ch)
			}
		} else {
			s.WriteByte(l.ch)
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar()
	} else {
		return Token{Type: ILLEGAL, Literal: "unterminated string", Line: l.line, Col: l.col}
	}
	tok.Literal = s.String()
	return tok
}

func (l *Lexer) readNumber() Token {
	tok := Token{Type: NUMBER, Literal: "", Line: l.line, Col: l.col}
	var s strings.Builder
	hasDot := false
	if l.ch == '.' {
		s.WriteByte('.')
		hasDot = true
		l.readChar()
	}
	for isDigit(l.ch) {
		s.WriteByte(l.ch)
		l.readChar()
		if !hasDot && l.ch == '.' {
			s.WriteByte('.')
			hasDot = true
			l.readChar()
		}
	}
	tok.Literal = s.String()
	return tok
}

func (l *Lexer) readIdentifierOrKeyword() Token {
	tok := Token{Line: l.line, Col: l.col}
	start := l.pos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	word := l.input[start:l.pos]

	if word == "var" && l.ch == '.' {
		l.readChar()
		tok.Type = VAR
		tok.Literal = "var."
		return tok
	}

	if word == "add" && l.ch == '.' {
		l.readChar()
		tok.Type = ADD
		tok.Literal = "add."
		return tok
	}

	if word == "func" && l.ch == '.' {
		l.readChar()
		tok.Type = FUNC
		tok.Literal = "func."
		return tok
	}

	if word == "input" && l.ch == '.' {
		l.readChar()
		tok.Type = INPUT
		tok.Literal = "input."
		return tok
	}

	if word == "run" && l.ch == '.' {
		l.readChar()
		tok.Type = RUN
		tok.Literal = "run."
		return tok
	}

	switch word {
	case "print":
		tok.Type = PRINT
		tok.Literal = word
	case "True":
		tok.Type = BOOLEAN
		tok.Literal = word
	case "False":
		tok.Type = BOOLEAN
		tok.Literal = word
	case "Null":
		tok.Type = NULL
		tok.Literal = word
	case "if":
		tok.Type = IF
		tok.Literal = word
	case "elif":
		tok.Type = ELIF
		tok.Literal = word
	case "else":
		tok.Type = ELSE
		tok.Literal = word
	default:
		tok.Type = IDENTIFIER
		tok.Literal = word
	}

	return tok
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
