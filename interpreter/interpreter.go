package interpreter

import (
	"fmt"
	"strconv"

	"emerald/ast"
)

type Interpreter struct {
	env map[string]interface{}
}

func NewInterpreter() *Interpreter {
	return &Interpreter{env: make(map[string]interface{})}
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

func (i *Interpreter) evalExpression(expr ast.Expression) (interface{}, error) {
	switch e := expr.(type) {
	case *ast.StringLiteral:
		return e.Value, nil
	case *ast.NumberLiteral:
		return e.Value, nil
	case *ast.BooleanLiteral:
		return e.Value, nil
	case *ast.Identifier:
		val, ok := i.env[e.Value]
		if !ok {
			return nil, fmt.Errorf("undefined variable '%s'", e.Value)
		}
		return val, nil
	}
	return nil, fmt.Errorf("unknown expression type")
}

func formatValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "True"
		}
		return "False"
	}
	return fmt.Sprintf("%v", val)
}
