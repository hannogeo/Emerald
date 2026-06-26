package lexer

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

	if word == "fn" && l.ch == '.' {
		l.readChar()
		tok.Type = FUNC
		tok.Literal = "fn."
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
	case "input":
		tok.Type = INPUT
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
