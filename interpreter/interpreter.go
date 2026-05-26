package interpreter

import (
	"fmt"

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
	case *ast.NullLiteral:
		return nil, nil
	case *ast.Identifier:
		val, ok := i.env[e.Value]
		if !ok {
			return nil, fmt.Errorf("undefined variable '%s'", e.Value)
		}
		return val, nil
	case *ast.BinaryExpression:
		return i.evalBinaryExpression(e)
	case *ast.CallExpression:
		return i.evalCallExpression(e)
	}
	return nil, fmt.Errorf("unknown expression type")
}

func (i *Interpreter) evalBinaryExpression(e *ast.BinaryExpression) (interface{}, error) {
	left, err := i.evalExpression(e.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evalExpression(e.Right)
	if err != nil {
		return nil, err
	}
	return binaryOperation(left, e.Operator, right, e.Line)
}

func (i *Interpreter) evalCallExpression(e *ast.CallExpression) (interface{}, error) {
	arg, err := i.evalExpression(e.Argument)
	if err != nil {
		return nil, err
	}

	switch e.Function {
	case "str":
		return builtinStr(arg, e.Line)
	case "num":
		return builtinNum(arg, e.Line)
	default:
		return nil, fmt.Errorf("undefined function '%s' at line %d", e.Function, e.Line)
	}
}

func formatValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case float64:
		return formatNumber(v)
	case bool:
		if v {
			return "True"
		}
		return "False"
	case nil:
		return "Null"
	}
	return fmt.Sprintf("%v", val)
}

func formatNumber(v float64) string {
	if v == float64(int64(v)) {
		return fmt.Sprintf("%d", int64(v))
	}
	return fmt.Sprintf("%g", v)
}
