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

type InterpolationPart struct {
	Text string
	Expr Expression
}

type InterpolatedStringLiteral struct {
	Parts []InterpolationPart
}

func (isl *InterpolatedStringLiteral) expressionNode() {}
func (isl *InterpolatedStringLiteral) String() string {
	var out string
	for _, p := range isl.Parts {
		if p.Expr != nil {
			out += "{" + p.Expr.String() + "}"
		} else {
			out += p.Text
		}
	}
	return "$" + out
}

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

type PrefixExpression struct {
	Operator string
	Right    Expression
	Line     int
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s %s)", pe.Operator, pe.Right.String())
}

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

type ListLiteral struct {
	Elements []Expression
}

func (ll *ListLiteral) expressionNode() {}
func (ll *ListLiteral) String() string {
	out := "("
	for i, e := range ll.Elements {
		if i > 0 {
			out += ", "
		}
		out += e.String()
	}
	return out + ")"
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
func (fs *FuncStatement) String() string { return fmt.Sprintf("fn.%s { ... }", fs.Name) }

type RunStatement struct {
	Name string
}

func (rs *RunStatement) statementNode() {}
func (rs *RunStatement) String() string { return fmt.Sprintf("run %s", rs.Name) }

type AddStatement struct {
	Name  string
	Value Expression
}

func (as *AddStatement) statementNode() {}
func (as *AddStatement) String() string { return fmt.Sprintf("add %s %s", as.Name, as.Value.String()) }

type InputExpression struct {
	Prompt Expression
}

func (ie *InputExpression) expressionNode() {}
func (ie *InputExpression) String() string { return fmt.Sprintf("input %s", ie.Prompt.String()) }

type RangeExpression struct {
	Start Expression
	End   Expression
}

func (re *RangeExpression) expressionNode() {}
func (re *RangeExpression) String() string {
	if re.Start == nil {
		return fmt.Sprintf("range:%s", re.End.String())
	}
	return fmt.Sprintf("range:(%s, %s)", re.Start.String(), re.End.String())
}

type ListIndexExpression struct {
	Name  string
	Index Expression
	Line  int
}

func (lie *ListIndexExpression) expressionNode() {}
func (lie *ListIndexExpression) String() string {
	return fmt.Sprintf("%s:%s", lie.Name, lie.Index.String())
}

type ListSliceExpression struct {
	Name  string
	Start Expression
	End   Expression
	Line  int
}

func (lse *ListSliceExpression) expressionNode() {}
func (lse *ListSliceExpression) String() string {
	return fmt.Sprintf("%s:(%s, %s)", lse.Name, lse.Start.String(), lse.End.String())
}

type ForStatement struct {
	Variable string
	Iterable Expression
	Body     *BlockStatement
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	return fmt.Sprintf("for %s in %s { ... }", fs.Variable, fs.Iterable.String())
}

type TypeLiteral struct {
	TypeName string
}

func (tl *TypeLiteral) expressionNode() {}
func (tl *TypeLiteral) String() string  { return tl.TypeName }

type WhileStatement struct {
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) String() string {
	return fmt.Sprintf("while %s { ... }", ws.Condition.String())
}

type BreakStatement struct{}

func (bs *BreakStatement) statementNode() {}
func (bs *BreakStatement) String() string { return "break" }

type ContinueStatement struct{}

func (cs *ContinueStatement) statementNode() {}
func (cs *ContinueStatement) String() string { return "continue" }

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
