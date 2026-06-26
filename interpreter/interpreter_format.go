package interpreter

import (
	"fmt"
	"strings"
)

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
