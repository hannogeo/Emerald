package interpreter

import (
	"fmt"
	"strconv"
)

func builtinStr(arg interface{}, line int) (interface{}, error) {
	switch v := arg.(type) {
	case float64:
		return formatNumber(v), nil
	case bool:
		if v {
			return "True", nil
		}
		return "False", nil
	case string:
		return v, nil
	case nil:
		return "Null", nil
	default:
		return nil, fmt.Errorf("str() argument type not supported at line %d", line)
	}
}

func builtinNum(arg interface{}, line int) (interface{}, error) {
	switch v := arg.(type) {
	case string:
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert '%s' to number at line %d", v, line)
		}
		return val, nil
	case float64:
		return v, nil
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	case nil:
		return 0.0, nil
	default:
		return nil, fmt.Errorf("num() argument type not supported at line %d", line)
	}
}
