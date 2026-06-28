package interpreter

import (
	"fmt"
	"strings"
)

func binaryOperation(left interface{}, operator string, right interface{}, line int) (interface{}, error) {
	switch operator {
	case "=":
		if ov, ok := right.(*orValue); ok {
			for _, val := range ov.Values {
				if matchEq(left, val) {
					return true, nil
				}
			}
			return false, nil
		}
		if nv, ok := right.(*notValue); ok {
			return left != nv.Value, nil
		}
		if av, ok := right.(*andValue); ok {
			if len(av.Values) == 0 {
				return false, nil
			}
			return matchEq(left, av.Values[0]), nil
		}
		if tc, ok := right.(*typeCheck); ok {
			switch tc.TypeName {
			case "num":
				_, ok := left.(float64)
				return ok, nil
			case "str":
				_, ok := left.(string)
				return ok, nil
			case "bool":
				_, ok := left.(bool)
				return ok, nil
			default:
				return false, fmt.Errorf("unknown type '%s' at line %d", tc.TypeName, line)
			}
		}
		return left == right, nil
	case "<", ">", "<=", ">=":
		lNum, ok1 := left.(float64)
		rNum, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("comparison operator '%s' requires numbers at line %d", operator, line)
		}
		switch operator {
		case "<":
			return lNum < rNum, nil
		case ">":
			return lNum > rNum, nil
		case "<=":
			return lNum <= rNum, nil
		case ">=":
			return lNum >= rNum, nil
		}
	}

	switch l := left.(type) {
	case float64:
		return numberOp(l, operator, right, line)
	case string:
		return stringOp(l, operator, right, line)
	case bool:
		return nil, fmt.Errorf("operator '%s' not supported for booleans at line %d", operator, line)
	case nil:
		return nil, fmt.Errorf("operator '%s' not supported for Null at line %d", operator, line)
	default:
		return nil, fmt.Errorf("unknown type at line %d", line)
	}
}

func numberOp(left float64, operator string, right interface{}, line int) (interface{}, error) {
	r, ok := right.(float64)
	if !ok {
		return nil, fmt.Errorf("type mismatch: cannot use '%s' between number and %s at line %d",
			operator, typeName(right), line)
	}

	switch operator {
	case "+":
		return left + r, nil
	case "-":
		return left - r, nil
	case "*":
		return left * r, nil
	case "/":
		if r == 0 {
			return nil, fmt.Errorf("division by zero at line %d", line)
		}
		return left / r, nil
	default:
		return nil, fmt.Errorf("unknown operator '%s' at line %d", operator, line)
	}
}

func stringOp(left string, operator string, right interface{}, line int) (interface{}, error) {
	switch operator {
	case "+":
		r, ok := right.(string)
		if !ok {
			return nil, fmt.Errorf("type mismatch: cannot concatenate string with %s at line %d",
				typeName(right), line)
		}
		return left + r, nil

	case "-":
		r, ok := right.(string)
		if !ok {
			return nil, fmt.Errorf("type mismatch: cannot use '-' between string and %s at line %d",
				typeName(right), line)
		}
		idx := strings.Index(left, r)
		if idx == -1 {
			return nil, fmt.Errorf("substring '%s' not found in '%s' at line %d", r, left, line)
		}
		result := left[:idx] + left[idx+len(r):]
		if result == "" {
			return nil, nil
		}
		return result, nil

	case "*":
		r, ok := right.(float64)
		if !ok {
			return nil, fmt.Errorf("type mismatch: cannot use '*' between string and %s at line %d",
				typeName(right), line)
		}
		return strings.Repeat(left, int(r)), nil

	case "/":
		return nil, fmt.Errorf("operator '/' not supported for strings at line %d", line)

	default:
		return nil, fmt.Errorf("unknown operator '%s' at line %d", operator, line)
	}
}

func matchEq(left interface{}, val interface{}) bool {
	if nv, ok := val.(*notValue); ok {
		return left != nv.Value
	}
	return left == val
}

func typeName(val interface{}) string {
	switch val.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case nil:
		return "Null"
	case *notValue:
		return "not"
	case *orValue:
		return "or"
	case *andValue:
		return "and"
	case *typeCheck:
		return "type"
	default:
		return "unknown"
	}
}
