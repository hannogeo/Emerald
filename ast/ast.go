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
func (cs *CommentStatement) String() string { return fmt.Sprintf("// %s", cs.Text) }

type VarStatement struct {
	Name  string
	Value Expression
}

func (vs *VarStatement) statementNode() {}
func (vs *VarStatement) String() string { return fmt.Sprintf("var.%s %s", vs.Name, vs.Value.String()) }

type PrintStatement struct {
	Value Expression
}

func (ps *PrintStatement) statementNode() {}
func (ps *PrintStatement) String() string { return fmt.Sprintf("print %s", ps.Value.String()) }

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return fmt.Sprintf("%q", sl.Value) }

type NumberLiteral struct {
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return fmt.Sprintf("%v", nl.Value) }

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

type NullLiteral struct{}

func (nl *NullLiteral) expressionNode() {}
func (nl *NullLiteral) String() string  { return "Null" }

type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
	Line     int
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Left.String(), be.Operator, be.Right.String())
}

type CallExpression struct {
	Function string
	Argument Expression
	Line     int
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	return fmt.Sprintf("%s(%s)", ce.Function, ce.Argument.String())
}

type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	var out string
	for _, s := range bs.Statements {
		out += s.String() + "\n"
	}
	return out
}

type FuncStatement struct {
	Name string
	Body *BlockStatement
}

func (fs *FuncStatement) statementNode() {}
func (fs *FuncStatement) String() string { return fmt.Sprintf("func.%s { ... }", fs.Name) }

type RunStatement struct {
	Name string
}

func (rs *RunStatement) statementNode() {}
func (rs *RunStatement) String() string { return fmt.Sprintf("run.%s", rs.Name) }

type IfStatement struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative interface{} // *IfStatement (elif) or *BlockStatement (else) or nil
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	out := fmt.Sprintf("if %s { ... }", is.Condition.String())
	if is.Alternative != nil {
		switch alt := is.Alternative.(type) {
		case *IfStatement:
			out += " elif " + alt.String()
		case *BlockStatement:
			out += " else { ... }"
		}
	}
	return out
}
