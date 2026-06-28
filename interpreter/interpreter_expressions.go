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
	case *ast.PrefixExpression:
		return i.evalPrefixExpression(e)
	case *ast.TypeLiteral:
		return &typeCheck{TypeName: e.TypeName}, nil
	case *ast.MethodCallExpression:
		return i.evalMethodCallExpression(e)
	}
	return nil, fmt.Errorf("unknown expression type")
}

func (i *Interpreter) evalPrefixExpression(e *ast.PrefixExpression) (interface{}, error) {
	right, err := i.evalExpression(e.Right)
	if err != nil {
		return nil, err
	}
	if e.Operator == "not" {
		return &notValue{Value: right}, nil
	}
	return nil, fmt.Errorf("unknown prefix operator '%s' at line %d", e.Operator, e.Line)
}

func (i *Interpreter) evalBinaryExpression(e *ast.BinaryExpression) (interface{}, error) {
	left, err := i.evalExpression(e.Left)
	if err != nil {
		return nil, err
	}

	if e.Operator == "or" {
		right, err := i.evalExpression(e.Right)
		if err != nil {
			return nil, err
		}
		return mergeOrValues(left, right), nil
	}

	if e.Operator == "and" {
		right, err := i.evalExpression(e.Right)
		if err != nil {
			return nil, err
		}
		return mergeAndValues(left, right), nil
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

func (i *Interpreter) evalMethodCallExpression(e *ast.MethodCallExpression) (interface{}, error) {
	obj, err := i.evalExpression(e.Object)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, len(e.Args))
	for j, arg := range e.Args {
		args[j], err = i.evalExpression(arg)
		if err != nil {
			return nil, err
		}
	}

	switch e.Method {
	case "slice":
		if len(args) != 2 {
			return nil, fmt.Errorf("slice requires 2 arguments (start, length) at line %d", e.Line)
		}
		start, ok := args[0].(float64)
		if !ok {
			return nil, fmt.Errorf("slice start must be a number at line %d", e.Line)
		}
		length, ok := args[1].(float64)
		if !ok {
			return nil, fmt.Errorf("slice length must be a number at line %d", e.Line)
		}
		switch v := obj.(type) {
		case string:
			runes := []rune(v)
			s := int(start) - 1
			l := int(length)
			if s < 0 {
				s = 0
			}
			if s > len(runes) {
				s = len(runes)
			}
			if s+l > len(runes) {
				l = len(runes) - s
			}
			return string(runes[s : s+l]), nil
		case []interface{}:
			s := int(start) - 1
			l := int(length)
			if s < 0 {
				s = 0
			}
			if s > len(v) {
				s = len(v)
			}
			if s+l > len(v) {
				l = len(v) - s
			}
			return v[s : s+l], nil
		default:
			return nil, fmt.Errorf("slice requires a string or list at line %d", e.Line)
		}

	case "upper":
		str, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("upper requires a string at line %d", e.Line)
		}
		return strings.ToUpper(str), nil

	case "lower":
		str, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("lower requires a string at line %d", e.Line)
		}
		return strings.ToLower(str), nil

	case "trim":
		str, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("trim requires a string at line %d", e.Line)
		}
		return strings.TrimSpace(str), nil

	case "split":
		str, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("split requires a string at line %d", e.Line)
		}
		if len(args) != 1 {
			return nil, fmt.Errorf("split requires 1 argument (delimiter) at line %d", e.Line)
		}
		delim, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("split delimiter must be a string at line %d", e.Line)
		}
		parts := strings.Split(str, delim)
		result := make([]interface{}, len(parts))
		for i, p := range parts {
			result[i] = p
		}
		return result, nil

	case "join":
		list, ok := obj.([]interface{})
		if !ok {
			return nil, fmt.Errorf("join requires a list at line %d", e.Line)
		}
		if len(args) != 1 {
			return nil, fmt.Errorf("join requires 1 argument (delimiter) at line %d", e.Line)
		}
		delim, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("join delimiter must be a string at line %d", e.Line)
		}
		parts := make([]string, len(list))
		for i, v := range list {
			s, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("join requires a list of strings at line %d", e.Line)
			}
			parts[i] = s
		}
		return strings.Join(parts, delim), nil

	case "replace":
		str, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("replace requires a string at line %d", e.Line)
		}
		if len(args) != 2 {
			return nil, fmt.Errorf("replace requires 2 arguments (old, new) at line %d", e.Line)
		}
		old, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("replace old must be a string at line %d", e.Line)
		}
		new, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("replace new must be a string at line %d", e.Line)
		}
		return strings.Replace(str, old, new, -1), nil

	case "contains":
		str, ok := obj.(string)
		if !ok {
			return nil, fmt.Errorf("contains requires a string at line %d", e.Line)
		}
		if len(args) != 1 {
			return nil, fmt.Errorf("contains requires 1 argument (substring) at line %d", e.Line)
		}
		sub, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("contains substring must be a string at line %d", e.Line)
		}
		return strings.Contains(str, sub), nil

	case "len":
		switch v := obj.(type) {
		case string:
			return float64(len([]rune(v))), nil
		case []interface{}:
			return float64(len(v)), nil
		default:
			return nil, fmt.Errorf("len requires a string or list at line %d", e.Line)
		}

	default:
		return nil, fmt.Errorf("unknown method '%s' at line %d", e.Method, e.Line)
	}
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
