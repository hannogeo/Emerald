package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	case *ast.IfStatement:
		return i.evalIfStatement(s)
	case *ast.FuncStatement:
		return i.evalFuncStatement(s)
	case *ast.RunStatement:
		return i.evalRunStatement(s)
	case *ast.AddStatement:
		return i.evalAddStatement(s)
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

func (i *Interpreter) evalExpression(expr ast.Expression) (interface{}, error) {
	switch e := expr.(type) {
	case *ast.StringLiteral:
		return e.Value, nil
	case *ast.InterpolatedStringLiteral:
		return i.evalInterpolatedString(e)
	case *ast.InputExpression:
		return i.evalInputExpression(e)
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
	case *ast.ListLiteral:
		return i.evalListLiteral(e)
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

func (i *Interpreter) evalListLiteral(e *ast.ListLiteral) (interface{}, error) {
	list := make([]interface{}, 0, len(e.Elements))
	for _, elem := range e.Elements {
		val, err := i.evalExpression(elem)
		if err != nil {
			return nil, err
		}
		list = append(list, val)
	}
	return list, nil
}

func (i *Interpreter) evalInputExpression(e *ast.InputExpression) (interface{}, error) {
	prompt, err := i.evalExpression(e.Prompt)
	if err != nil {
		return nil, err
	}
	promptStr, ok := prompt.(string)
	if !ok {
		return nil, fmt.Errorf("input prompt must be a string")
	}
	fmt.Print(promptStr)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func (i *Interpreter) evalAddStatement(stmt *ast.AddStatement) error {
	val, ok := i.env[stmt.Name]
	if !ok {
		return fmt.Errorf("undefined list '%s'", stmt.Name)
	}
	list, ok := val.([]interface{})
	if !ok {
		return fmt.Errorf("'%s' is not a list", stmt.Name)
	}
	item, err := i.evalExpression(stmt.Value)
	if err != nil {
		return err
	}
	i.env[stmt.Name] = append(list, item)
	return nil
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
	}

	val, ok := i.env[e.Function]
	if !ok {
		return nil, fmt.Errorf("undefined function or list '%s' at line %d", e.Function, e.Line)
	}
	list, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' is not a list at line %d", e.Function, e.Line)
	}
	idx, ok := arg.(float64)
	if !ok {
		return nil, fmt.Errorf("list index must be a number at line %d", e.Line)
	}
	if idx < 1 || int(idx) > len(list) {
		return nil, fmt.Errorf("list index out of range at line %d", e.Line)
	}
	return list[int(idx)-1], nil
}

func (i *Interpreter) evalInterpolatedString(e *ast.InterpolatedStringLiteral) (interface{}, error) {
	var result strings.Builder
	for _, part := range e.Parts {
		if part.Expr != nil {
			val, err := i.evalExpression(part.Expr)
			if err != nil {
				return nil, err
			}
			s, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("interpolated expression must be a string, use str() to convert %s", typeName(val))
			}
			result.WriteString(s)
		} else {
			result.WriteString(part.Text)
		}
	}
	return result.String(), nil
}

func (i *Interpreter) evalIfStatement(stmt *ast.IfStatement) error {
	cond, err := i.evalExpression(stmt.Condition)
	if err != nil {
		return err
	}
	condBool, ok := cond.(bool)
	if !ok {
		return fmt.Errorf("condition must be a boolean, got %s", typeName(cond))
	}

	if condBool {
		return i.evalBlockStatement(stmt.Consequence)
	}

	if stmt.Alternative != nil {
		switch alt := stmt.Alternative.(type) {
		case *ast.IfStatement:
			return i.evalIfStatement(alt)
		case *ast.BlockStatement:
			return i.evalBlockStatement(alt)
		}
	}
	return nil
}

func (i *Interpreter) evalFuncStatement(stmt *ast.FuncStatement) error {
	i.env[stmt.Name] = stmt.Body
	return nil
}

func (i *Interpreter) evalRunStatement(stmt *ast.RunStatement) error {
	val, ok := i.env[stmt.Name]
	if !ok {
		return fmt.Errorf("undefined function '%s'", stmt.Name)
	}
	block, ok := val.(*ast.BlockStatement)
	if !ok {
		return fmt.Errorf("'%s' is not a function", stmt.Name)
	}
	return i.evalBlockStatement(block)
}

func (i *Interpreter) evalBlockStatement(block *ast.BlockStatement) error {
	for _, stmt := range block.Statements {
		err := i.evalStatement(stmt)
		if err != nil {
			return err
		}
	}
	return nil
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
	case []interface{}:
		parts := make([]string, len(v))
		for i, elem := range v {
			parts[i] = formatValue(elem)
		}
		return strings.Join(parts, ", ")
	}
	return fmt.Sprintf("%v", val)
}

func formatNumber(v float64) string {
	if v == float64(int64(v)) {
		return fmt.Sprintf("%d", int64(v))
	}
	return fmt.Sprintf("%g", v)
}
