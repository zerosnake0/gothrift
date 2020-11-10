package ast

import (
	"fmt"
)

type ConstList struct {
	Node
	Content []ConstListItem
}

var _ ConstValue = ConstList{}

func (c2 ConstList) Format(f fmt.State, c rune) {
	panic("implement me")
}

type ConstListItem struct {
	ConstValue
	End *Identifier
}
