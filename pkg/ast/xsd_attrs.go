package ast

type XsdAttrs struct {
	Start     Pos
	LBrace    Pos
	FieldList []Field
	RBrace    Pos
}

var _ Span = XsdAttrs{}

func (x XsdAttrs) StartPos() Pos {
	return x.Start
}

func (x XsdAttrs) EndPos() Pos {
	end := x.RBrace
	end.Col++
	return end
}
