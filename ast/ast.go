package ast

import (
	"gosling/token"
)

type Node interface {
	Tokenliteral() string
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
	Statement []Statement
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

// Program methods
func (p *Program) Tokenliteral() string {
	if len(p.Statement) > 0 {
		return p.Statement[0].Tokenliteral()
	} else {
		return ""
	}
}

// Identifier methods
func (i *Identifier) expressionNode()      {}
func (i *Identifier) Tokenliteral() string { return i.Token.Literal }

// LetStatement methods
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) Tokenliteral() string { return ls.Token.Literal }
