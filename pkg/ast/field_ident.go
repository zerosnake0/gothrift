package ast

type FieldIdentifier struct {
	Number Number
	Colon  Pos
}

var _ Span = FieldIdentifier{}

func (f FieldIdentifier) StartPos() Pos {
	return f.Number.Start
}

func (f FieldIdentifier) EndPos() Pos {
	end := f.Colon
	end.Col++
	return end
}
