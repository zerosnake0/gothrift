package ast

import (
	"fmt"
)

type Exception struct {
	Start       Pos
	Identifier  Identifier
	LBrace      Pos
	FieldList   []Field
	RBrace      Pos
	Annotations *Annotations
}

var _ Definition = Exception{}

func (e Exception) StartPos() Pos {
	return e.Start
}

func (e Exception) EndPos() Pos {
	if e.Annotations != nil {
		return e.Annotations.End
	}
	end := e.RBrace
	end.Col++
	return end
}

func (e Exception) Format(f fmt.State, c rune) {
	panic("implement me")
}
