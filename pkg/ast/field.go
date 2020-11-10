package ast

type Field struct {
	FieldIdentifier *FieldIdentifier
	Requiredness    *Identifier
	FieldType       FieldType
	Reference       *Identifier
	Identifier      Identifier
	FieldValue      *FieldValue
	XsdOptional     *Identifier
	XsdNillable     *Identifier
	XsdAttrs        *XsdAttrs
	Annotations     *Annotations
	End             *Identifier
}

var _ Span = Field{}

func (f Field) StartPos() Pos {
	if f.FieldIdentifier != nil {
		return f.FieldIdentifier.StartPos()
	}
	if f.Requiredness != nil {
		return f.Requiredness.Start
	}
	return f.FieldType.StartPos()
}

func (f Field) EndPos() Pos {
	if f.End != nil {
		return f.End.End
	}
	if f.Annotations != nil {
		return f.Annotations.End
	}
	if f.XsdAttrs != nil {
		return f.XsdAttrs.EndPos()
	}
	if f.XsdNillable != nil {
		return f.XsdNillable.End
	}
	if f.XsdOptional != nil {
		return f.XsdOptional.End
	}
	if f.FieldValue != nil {
		return f.FieldValue.EndPos()
	}
	return f.Identifier.End
}
