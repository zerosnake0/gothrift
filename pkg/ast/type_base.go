package ast

import (
	"fmt"
)

type BaseType struct {
	Identifier  Identifier
	Annotations *Annotations
}

var _ FieldType = BaseType{}

func (b BaseType) StartPos() Pos {
	return b.Identifier.Start
}

func (b BaseType) EndPos() Pos {
	if b.Annotations == nil {
		return b.Identifier.End
	}
	return b.Annotations.End
}

func (b BaseType) Format(f fmt.State, c rune) {
	panic("implement me")
}
