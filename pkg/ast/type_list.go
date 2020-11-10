package ast

import (
	"fmt"
)

type ListType struct {
	Start    Pos
	LChevron Pos
	Elem     FieldType
	RChevron Pos
	CppType  *CppType
}

var _ SimpleContainerType = ListType{}

func (l ListType) StartPos() Pos {
	return l.Start
}

func (l ListType) EndPos() Pos {
	if l.CppType == nil {
		end := l.RChevron
		end.Col++
		return end
	}
	return l.CppType.Literal.End
}

func (l ListType) Format(f fmt.State, c rune) {
	panic("implement me")
}
