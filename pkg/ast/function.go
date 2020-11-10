package ast

type Function struct {
	Oneway       *Identifier
	FunctionType FunctionType
	Identifier   Identifier
	LChevron     Pos
	FieldList    []Field
	RChevron     Pos
	Throws       *Throws
	Annotations  *Annotations
	End          *Identifier
}

var _ Span = Function{}

func (f Function) StartPos() Pos {
	if f.Oneway != nil {
		return f.Oneway.Start
	}
	return f.FunctionType.StartPos()
}

func (f Function) EndPos() Pos {
	if f.End != nil {
		return f.End.End
	}
	if f.Annotations != nil {
		return f.Annotations.End
	}
	if f.Throws != nil {
		return f.Throws.EndPos()
	}
	end := f.RChevron
	end.Col++
	return end
}
