package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeTypeDef(def ast.TypeDef) {
	f.encodeKeyword(def.Start, "typedef")
	f.forward(true, def.FieldType.StartPos())
	f.encodeFieldType(def.FieldType)
	f.forward(true, def.Identifier.Start)
	f.encodeIdentifier(&def.Identifier)
	if def.Annotations != nil {
		f.encodeAnnotations(def.Annotations)
	}
}
