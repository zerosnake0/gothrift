package ast

type Throws struct {
	Start     Pos
	LChevron  Pos
	FieldList []Field
	RChevron  Pos
}

var _ Span = Throws{}

func (t Throws) StartPos() Pos {
	return t.Start
}

func (t Throws) EndPos() Pos {
	end := t.RChevron
	end.Col++
	return end
}
