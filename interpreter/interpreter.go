package interpreter

import (
	"bufio"
	"fmt"
	"os"

	"emerald/ast"
)

type Interpreter struct {
	env    map[string]interface{}
	reader *bufio.Reader
}

func NewInterpreter() *Interpreter {
	return &Interpreter{env: make(map[string]interface{}), reader: bufio.NewReader(os.Stdin)}
}

func (i *Interpreter) Eval(program *ast.Program) error {
	for _, stmt := range program.Statements {
		err := i.evalStatement(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) evalStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.VarStatement:
		return i.evalVarStatement(s)
	case *ast.PrintStatement:
		return i.evalPrintStatement(s)
	case *ast.IfStatement:
		return i.evalIfStatement(s)
	case *ast.FuncStatement:
		return i.evalFuncStatement(s)
	case *ast.RunStatement:
		return i.evalRunStatement(s)
	case *ast.AddStatement:
		return i.evalAddStatement(s)
	case *ast.ForStatement:
		return i.evalForStatement(s)
	case *ast.WhileStatement:
		return i.evalWhileStatement(s)
	case *ast.BlockStatement:
		return i.evalBlockStatement(s)
	}
	return nil
}

func (i *Interpreter) evalVarStatement(stmt *ast.VarStatement) error {
	val, err := i.evalExpression(stmt.Value)
	if err != nil {
		return err
	}
	i.env[stmt.Name] = val
	return nil
}

func (i *Interpreter) evalPrintStatement(stmt *ast.PrintStatement) error {
	val, err := i.evalExpression(stmt.Value)
	if err != nil {
		return err
	}
	fmt.Println(formatValue(val))
	return nil
}
