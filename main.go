package main

import (
	"fmt"
	"os"

	"emerald/interpreter"
	"emerald/lexer"
	"emerald/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: emerald <filename.emld>")
		os.Exit(1)
	}

	filename := os.Args[1]

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lexer.NewLexer(string(content))
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, errMsg := range p.Errors() {
			fmt.Fprintln(os.Stderr, "Parse error:", errMsg)
		}
		os.Exit(1)
	}

	interp := interpreter.NewInterpreter()
	err = interp.Eval(program)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Runtime error:", err)
		os.Exit(1)
	}
}
