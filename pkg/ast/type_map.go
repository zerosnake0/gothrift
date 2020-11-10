package ast

import (
	"fmt"
)

type MapType struct {
	Start    Pos
	CppType  *CppType
	LChevron Pos
	Key      FieldType
	Comma    Pos
	Value    FieldType
	RChevron Pos
}

var _ SimpleContainerType = MapType{}

func (m MapType) StartPos() Pos {
	return m.Start
}

func (m MapType) EndPos() Pos {
	end := m.RChevron
	end.Col++
	return end
}

func (m MapType) Format(f fmt.State, c rune) {
	panic("implement me")
}
