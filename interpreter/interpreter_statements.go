package interpreter

import (
	"fmt"

	"emerald/ast"
)

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

func (i *Interpreter) evalIfStatement(stmt *ast.IfStatement) error {
	cond, err := i.evalExpression(stmt.Condition)
	if err != nil {
		return err
	}
	condBool, ok := cond.(bool)
	if !ok {
		switch v := cond.(type) {
		case *notValue:
			condBool = !isTruthy(v.Value)
		case *orValue:
			for _, val := range v.Values {
				if isTruthy(val) {
					condBool = true
					break
				}
			}
		default:
			return fmt.Errorf("condition must be a boolean, got %s", typeName(cond))
		}
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

func (i *Interpreter) evalForStatement(stmt *ast.ForStatement) error {
	switch iter := stmt.Iterable.(type) {
	case *ast.RangeExpression:
		var start, end float64
		if iter.Start != nil {
			val, err := i.evalExpression(iter.Start)
			if err != nil {
				return err
			}
			s, ok := val.(float64)
			if !ok {
				return fmt.Errorf("range start must be a number")
			}
			start = s
		} else {
			start = 1
		}
		val, err := i.evalExpression(iter.End)
		if err != nil {
			return err
		}
		end, ok := val.(float64)
		if !ok {
			return fmt.Errorf("range end must be a number")
		}
		for x := start; x <= end; x++ {
			i.env[stmt.Variable] = x
			err := i.evalBlockStatement(stmt.Body)
			if err != nil {
				return err
			}
		}
	case *ast.Identifier:
		val, ok := i.env[iter.Value]
		if !ok {
			return fmt.Errorf("undefined list '%s'", iter.Value)
		}
		list, ok := val.([]interface{})
		if !ok {
			return fmt.Errorf("'%s' is not a list", iter.Value)
		}
		for _, item := range list {
			i.env[stmt.Variable] = item
			err := i.evalBlockStatement(stmt.Body)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("for-in requires a range or list")
	}
	return nil
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
