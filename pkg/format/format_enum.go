package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeEnum(enum ast.Enum) {
	f.encodeKeyword(enum.Start, "enum")
	f.forward(true, enum.Identifier.Start)
	f.encodeIdentifier(&enum.Identifier)
	f.forward(true, enum.LBrace)
	f.encodeBrace(enum.LBrace, enum.RBrace, ",", enum.List, func(span ast.Span) {
		f.encodeEnumDef(span.(*ast.EnumDef))
	})
	if enum.Annotations != nil {
		f.encodeAnnotations(enum.Annotations)
	}
}

func (f *Formatter) encodeEnumDef(item *ast.EnumDef) {
	f.encodeEnumValue(item.Value)
	if item.Annotations != nil {
		f.encodeAnnotations(item.Annotations)
	}
	f.encodeEndSeparator(item.End)
}

func (f *Formatter) encodeEnumValue(value ast.EnumValue) {
	switch x := value.(type) {
	case ast.EnumValueWithNumber:
		f.encodeEnumValueWithNumber(x)
	case ast.Identifier:
		f.encodeIdentifier(&x)
	default:
		shouldNotReach()
	}
}

func (f *Formatter) encodeEnumValueWithNumber(value ast.EnumValueWithNumber) {
	f.encodeIdentifier(&value.Identifier)
	f.forward(true, value.Eq)
	f.encodeKeyword(value.Eq, "=")
	f.forward(true, value.Number.Start)
	f.encodeNumber(value.Number)
}
