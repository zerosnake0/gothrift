package ast

type FieldValue struct {
	Start Pos
	Value ConstValue
}

var _ Span = FieldValue{}

func (f FieldValue) StartPos() Pos {
	return f.Start
}

func (f FieldValue) EndPos() Pos {
	return f.Value.EndPos()
}
