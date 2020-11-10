package ast

import (
	"fmt"
)

type Senum struct {
	Start       Pos
	Identifier  Identifier
	LBrace      Pos
	List        []SenumDef
	RBrace      Pos
	Annotations *Annotations
}

var _ Definition = &Senum{}

func (e Senum) StartPos() Pos {
	return e.Start
}

func (e Senum) EndPos() Pos {
	if e.Annotations != nil {
		return e.Annotations.End
	}
	end := e.RBrace
	end.Col++
	return end
}

func (e Senum) Format(f fmt.State, c rune) {
	panic("implement me")
}
