package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeEnum(enum ast.Enum) {
	f.encodeKeyword(enum.Start, "enum")
	f.forward(true, enum.Identifier.Start)
	f.encodeIdentifier(&enum.Identifier)
	f.forward(true, enum.LBrace)
	f.startBrace(enum.LBrace)
	for i := range enum.List {
		item := &enum.List[i]
		if i > 0 {
			f.print(",")
		}
		f.forward(false, item.StartPos())
		f.encodeEnumDef(item)
	}
	f.endBrace(enum.RBrace)
	if enum.Annotations != nil {
		f.encodeAnnotations(enum.Annotations)
	}
}

func (f *Formatter) encodeEnumDef(item *ast.EnumDef) {
	f.encodeEnumValue(item.Value)
	if item.Annotations != nil {
		f.encodeAnnotations(item.Annotations)
	}
}

func (f *Formatter) encodeEnumValue(value ast.EnumValue) {
	switch x := value.(type) {
	case ast.EnumValueWithNumber:
		f.encodeEnumValueWithNumber(x)
	case ast.Identifier:
		f.encodeIdentifier(&x)
	default:
		panic("should not reach")
	}
}

func (f *Formatter) encodeEnumValueWithNumber(value ast.EnumValueWithNumber) {
	f.encodeIdentifier(&value.Identifier)
	f.forward(true, value.Eq)
	f.encodeKeyword(value.Eq, "=")
	f.forward(true, value.Number.Start)
	f.encodeNumber(value.Number)
}
