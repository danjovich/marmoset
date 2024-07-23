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

type Expression interface {
	Node
	expressionNode()
}
