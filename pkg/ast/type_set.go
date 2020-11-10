package ast

import (
	"fmt"
)

type SetType struct {
	Start    Pos
	CppType  *CppType
	LChevron Pos
	Elem     FieldType
	RChevron Pos
}

var _ SimpleContainerType = SetType{}

func (s SetType) StartPos() Pos {
	return s.Start
}

func (s SetType) EndPos() Pos {
	end := s.RChevron
	end.Col++
	return end
}

func (s SetType) Format(f fmt.State, c rune) {
	panic("implement me")
}
