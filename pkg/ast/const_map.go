package ast

import (
	"fmt"
)

type ConstMap struct {
	Node
	Content []ConstMapItem
}

var _ ConstValue = ConstMap{}

func (c2 ConstMap) Format(f fmt.State, c rune) {
	panic("implement me")
}

type ConstMapItem struct {
	Key   ConstValue
	Colon Pos
	Value ConstValue
	End   *Identifier
}

var _ Span = ConstMapItem{}

func (c ConstMapItem) StartPos() Pos {
	return c.Key.StartPos()
}

func (c ConstMapItem) EndPos() Pos {
	if c.End != nil {
		return c.End.End
	}
	return c.Value.EndPos()
}
