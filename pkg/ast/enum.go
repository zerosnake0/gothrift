package ast

import (
	"fmt"
)

type Enum struct {
	Start       Pos
	Identifier  Identifier
	LBrace      Pos
	List        []EnumDef
	RBrace      Pos
	Annotations *Annotations
}

var _ Definition = &Enum{}

func (e Enum) StartPos() Pos {
	return e.Start
}

func (e Enum) EndPos() Pos {
	if e.Annotations != nil {
		return e.Annotations.End
	}
	end := e.RBrace
	end.Col++
	return end
}

func (e Enum) Format(f fmt.State, c rune) {
	panic("implement me")
}
