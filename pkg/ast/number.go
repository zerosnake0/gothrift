package ast

type Number struct {
	TextNode
}

var _ ConstValue = Number{}
