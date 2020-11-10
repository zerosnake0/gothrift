package ast

type SenumDef struct {
	Literal Literal
	End     *Identifier
}

func (s SenumDef) StartPos() Pos {
	return s.Literal.Start
}

func (s SenumDef) EndPos() Pos {
	if s.End != nil {
		return s.End.End
	}
	return s.Literal.End
}

var _ Span = SenumDef{}
