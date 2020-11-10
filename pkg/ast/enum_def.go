package ast

type EnumDef struct {
	Value       EnumValue
	Annotations *Annotations
	End         *Identifier
}

var _ Span = EnumDef{}

func (e EnumDef) StartPos() Pos {
	return e.Value.StartPos()
}

func (e EnumDef) EndPos() Pos {
	if e.End != nil {
		return e.End.End
	}
	if e.Annotations != nil {
		return e.Annotations.End
	}
	return e.Value.EndPos()
}
