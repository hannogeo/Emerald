package ast

import "fmt"

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out string
	for _, s := range p.Statements {
		out += s.String() + "\n"
	}
	return out
}

type CommentStatement struct {
	Text string
}

func (cs *CommentStatement) statementNode() {}
func (cs *CommentStatement) String() string  { return fmt.Sprintf("// %s", cs.Text) }

type VarStatement struct {
	Name  string
	Value Expression
}

func (vs *VarStatement) statementNode() {}
func (vs *VarStatement) String() string  { return fmt.Sprintf("var.%s %s", vs.Name, vs.Value.String()) }

type PrintStatement struct {
	Value Expression
}

func (ps *PrintStatement) statementNode() {}
func (ps *PrintStatement) String() string  { return fmt.Sprintf("print %s", ps.Value.String()) }

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string   { return fmt.Sprintf("%q", sl.Value) }

type NumberLiteral struct {
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string   { return fmt.Sprintf("%v", nl.Value) }

type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) String() string {
	if bl.Value {
		return "True"
	}
	return "False"
}

type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string   { return i.Value }
