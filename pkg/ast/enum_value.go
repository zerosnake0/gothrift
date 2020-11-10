package ast

type EnumValue interface {
	Span
}

type EnumValueWithNumber struct {
	Identifier Identifier
	Eq         Pos
	Number     Number
}

var _ EnumValue = EnumValueWithNumber{}

func (e EnumValueWithNumber) StartPos() Pos {
	return e.Identifier.Start
}

func (e EnumValueWithNumber) EndPos() Pos {
	return e.Number.End
}
