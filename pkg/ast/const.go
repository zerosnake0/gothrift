package ast

import (
	"fmt"
)

type Const struct {
	Start     Pos
	FieldType FieldType
	Key       Identifier
	Eq        Pos
	Value     ConstValue
	End       *Identifier
}

var _ Definition = Const{}

func (c Const) StartPos() Pos {
	return c.Start
}

func (c Const) EndPos() Pos {
	if c.End == nil {
		return c.Value.EndPos()
	}
	return c.End.End
}

func (c Const) Format(f fmt.State, r rune) {
	panic("implement me")
}

type ConstValue interface {
	TextSpan
}
