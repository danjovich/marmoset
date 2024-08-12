package ast

type Node interface {
	TokenLiteral() string
	// for debugging
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Type int

const (
	INT = iota
	CHAR
	BOOL
	VOID
	POINTER
)

type Expression interface {
	Node
	expressionNode()
	GetType() Type
}
