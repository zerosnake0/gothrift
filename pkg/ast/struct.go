package ast

import (
	"fmt"
)

type Struct struct {
	Head        Identifier
	Identifier  Identifier
	XsdAll      *Identifier
	LBrace      Pos
	List        []Field
	RBrace      Pos
	Annotations *Annotations
}

var _ Definition = &Struct{}

func (s Struct) StartPos() Pos {
	return s.Head.Start
}

func (s Struct) EndPos() Pos {
	if s.Annotations != nil {
		return s.Annotations.End
	}
	end := s.RBrace
	end.Col++
	return end
}

func (s Struct) Format(f fmt.State, c rune) {
	panic("implement me")
}
