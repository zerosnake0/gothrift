package ast

import (
	"fmt"
)

type TypeDef struct {
	Start       Pos
	FieldType   FieldType
	Identifier  Identifier
	Annotations *Annotations
	End         *Identifier
}

func (t TypeDef) StartPos() Pos {
	return t.Start
}

func (t TypeDef) EndPos() Pos {
	if t.End != nil {
		return t.End.End
	}
	if t.Annotations != nil {
		return t.Annotations.End
	}
	return t.Identifier.End
}

func (t TypeDef) Format(f fmt.State, c rune) {
	panic("implement me")
}

var _ Definition = TypeDef{}
