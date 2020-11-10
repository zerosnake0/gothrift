package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeStruct(st ast.Struct) {
	f.encodeIdentifier(&st.Head)
	f.forward(true, st.Identifier.Start)
	f.encodeIdentifier(&st.Identifier)
	if st.XsdAll != nil {
		f.forward(true, st.XsdAll.Start)
		f.encodeIdentifier(st.XsdAll)
	}
	f.forward(true, st.LBrace)
	f.encodeBraceFieldList(st.LBrace, st.RBrace, st.List)
	if st.Annotations != nil {
		f.encodeAnnotations(st.Annotations)
	}
}

func (f *Formatter) encodeField(field *ast.Field) {
	if field.FieldIdentifier != nil {
		f.encodeFieldIdentifier(field.FieldIdentifier)
	}
	if field.Requiredness != nil {
		if field.FieldIdentifier != nil {
			f.forward(true, field.Requiredness.Start)
		}
		f.encodeIdentifier(field.Requiredness)
	}
	if field.FieldIdentifier != nil || field.Requiredness != nil {
		f.forward(true, field.FieldType.StartPos())
	}
	f.encodeFieldType(field.FieldType)
	if field.Reference != nil {
		f.forwardAndEmptySep(true, field.Reference.Start)
		f.encodeIdentifier(field.Reference)
	}
	f.forward(true, field.Identifier.Start)
	f.encodeIdentifier(&field.Identifier)
	if field.FieldValue != nil {
		f.forward(true, field.FieldValue.Start)
		f.encodeFieldValue(field.FieldValue)
	}
	if field.XsdOptional != nil {
		f.forward(true, field.XsdOptional.Start)
		f.encodeIdentifier(field.XsdOptional)
	}
	if field.XsdNillable != nil {
		f.forward(true, field.XsdNillable.Start)
		f.encodeIdentifier(field.XsdNillable)
	}
	if field.XsdAttrs != nil {
		f.forward(true, field.XsdAttrs.Start)
		f.encodeXsdAttrs(field.XsdAttrs)
	}
	if field.Annotations != nil {
		f.encodeAnnotations(field.Annotations)
	}
}

func (f *Formatter) encodeFieldIdentifier(fi *ast.FieldIdentifier) {
	f.encodeNumber(fi.Number)
	f.forwardAndEmptySep(true, fi.Colon)
	f.encodeKeyword(fi.Colon, ":")
}

func (f *Formatter) encodeFieldValue(value *ast.FieldValue) {
	f.encodeKeyword(value.Start, "=")
	f.forward(true, value.Value.StartPos())
	f.encodeConstValue(value.Value)
}
