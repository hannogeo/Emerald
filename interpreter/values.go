package interpreter

import "fmt"

type notValue struct {
	Value interface{}
}

type orValue struct {
	Values []interface{}
}

type andValue struct {
	Values []interface{}
}

type typeCheck struct {
	TypeName string
}

func mergeOrValues(left, right interface{}) *orValue {
	var vals []interface{}
	if ov, ok := left.(*orValue); ok {
		vals = append(vals, ov.Values...)
	} else {
		vals = append(vals, left)
	}
	if ov, ok := right.(*orValue); ok {
		vals = append(vals, ov.Values...)
	} else {
		vals = append(vals, right)
	}
	return &orValue{Values: vals}
}

func mergeAndValues(left, right interface{}) *andValue {
	var vals []interface{}
	if av, ok := left.(*andValue); ok {
		vals = append(vals, av.Values...)
	} else {
		vals = append(vals, left)
	}
	if av, ok := right.(*andValue); ok {
		vals = append(vals, av.Values...)
	} else {
		vals = append(vals, right)
	}
	return &andValue{Values: vals}
}

func isTruthy(val interface{}) bool {
	switch v := val.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case string:
		return v != ""
	case nil:
		return false
	default:
		return true
	}
}

func (nv *notValue) String() string {
	return fmt.Sprintf("not(%v)", nv.Value)
}

func (ov *orValue) String() string {
	return fmt.Sprintf("or(%v)", ov.Values)
}
