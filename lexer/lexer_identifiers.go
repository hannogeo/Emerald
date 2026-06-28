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

	if word == "fn" && l.ch == '.' {
		l.readChar()
		tok.Type = FUNC
		tok.Literal = "fn."
		return tok
	}

	if word == "range" && l.ch == ':' {
		l.readChar()
		tok.Type = RANGE
		tok.Literal = "range:"
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
	case "for":
		tok.Type = FOR
		tok.Literal = word
	case "add":
		tok.Type = ADD
		tok.Literal = word
	case "run":
		tok.Type = RUN
		tok.Literal = word
	case "in":
		tok.Type = IN
		tok.Literal = word
	case "not":
		tok.Type = NOT
		tok.Literal = word
	case "or":
		tok.Type = OR
		tok.Literal = word
	case "and":
		tok.Type = AND
		tok.Literal = word
	case "while":
		tok.Type = WHILE
		tok.Literal = word
	case "num":
		tok.Type = NUM
		tok.Literal = word
	case "str":
		tok.Type = STR
		tok.Literal = word
	case "bool":
		tok.Type = BOOL
		tok.Literal = word
	case "break":
		tok.Type = BREAK
		tok.Literal = word
	case "continue":
		tok.Type = CONTINUE
		tok.Literal = word
	default:
		tok.Type = IDENTIFIER
		tok.Literal = word
	}

	return tok
}
