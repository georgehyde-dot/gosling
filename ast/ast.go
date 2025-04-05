package ast

import (
	"bytes"
	"gosling/token"
)

type Node interface {
	Tokenliteral() string
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

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

type LetStatement struct {
	Token token.Token // the token.Let token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

type ExpressionStatement struct {
	Token      token.Token // first token of expresion
	Expression Expression
}

// Program methods
func (p *Program) Tokenliteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].Tokenliteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Identifier methods
func (i *Identifier) expressionNode()      {}
func (i *Identifier) Tokenliteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// LetStatement methods
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) Tokenliteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.Tokenliteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()

}

// ReturnStatement methods
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) Tokenliteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.Tokenliteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement methods
// func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
