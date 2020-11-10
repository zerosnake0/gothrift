package ast

import (
	"fmt"
)

type Service struct {
	Start        Pos
	Identifier   Identifier
	Extends      *Extends
	LBrace       Pos
	FunctionList []Function
	RBrace       Pos
	Annotations  *Annotations
}

var _ Definition = Service{}

func (s Service) StartPos() Pos {
	return s.Start
}

func (s Service) EndPos() Pos {
	if s.Annotations != nil {
		return s.Annotations.End
	}
	end := s.RBrace
	end.Col++
	return end
}

func (s Service) Format(f fmt.State, c rune) {
	panic("implement me")
}
