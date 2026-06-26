package lexer

import "strings"

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
