package interpreter

import (
	"fmt"
	"strings"

	"emerald/ast"
)

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
	case *ast.ListIndexExpression:
		return i.evalListIndexExpression(e)
	case *ast.ListSliceExpression:
		return i.evalListSliceExpression(e)
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
	line, err := i.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return strings.TrimRight(line, "\r\n"), nil
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

func (i *Interpreter) evalListIndexExpression(e *ast.ListIndexExpression) (interface{}, error) {
	val, ok := i.env[e.Name]
	if !ok {
		return nil, fmt.Errorf("undefined variable '%s' at line %d", e.Name, e.Line)
	}
	list, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' is not a list at line %d", e.Name, e.Line)
	}
	idxVal, err := i.evalExpression(e.Index)
	if err != nil {
		return nil, err
	}
	idx, ok := idxVal.(float64)
	if !ok {
		return nil, fmt.Errorf("list index must be a number at line %d", e.Line)
	}
	if idx < 1 || int(idx) > len(list) {
		return nil, fmt.Errorf("list index out of range at line %d", e.Line)
	}
	return list[int(idx)-1], nil
}

func (i *Interpreter) evalListSliceExpression(e *ast.ListSliceExpression) (interface{}, error) {
	val, ok := i.env[e.Name]
	if !ok {
		return nil, fmt.Errorf("undefined variable '%s' at line %d", e.Name, e.Line)
	}
	list, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' is not a list at line %d", e.Name, e.Line)
	}
	startVal, err := i.evalExpression(e.Start)
	if err != nil {
		return nil, err
	}
	endVal, err := i.evalExpression(e.End)
	if err != nil {
		return nil, err
	}
	start, ok := startVal.(float64)
	if !ok {
		return nil, fmt.Errorf("slice start must be a number at line %d", e.Line)
	}
	end, ok := endVal.(float64)
	if !ok {
		return nil, fmt.Errorf("slice end must be a number at line %d", e.Line)
	}
	if start < 1 || int(end) > len(list) || start > end {
		return nil, fmt.Errorf("slice out of range at line %d", e.Line)
	}
	return list[int(start)-1 : int(end)], nil
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
