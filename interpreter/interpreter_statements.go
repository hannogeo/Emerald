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
