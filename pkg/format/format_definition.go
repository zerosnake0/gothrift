package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeDefinitions() {
	for _, def := range f.Doc.Definitions {
		f.encodeDefinition(def)
	}
}

func (f *Formatter) encodeDefinition(def ast.Definition) {
	f.forward(false, def.StartPos())
	f.newScope()
	switch x := def.(type) {
	case ast.Const:
		f.encodeConst(x)
	case ast.TypeDef:
		f.encodeTypeDef(x)
	case ast.Enum:
		f.encodeEnum(x)
	case ast.Senum:
		f.encodeSenum(x)
	case ast.Struct:
		f.encodeStruct(x)
	case ast.Exception:
		f.encodeException(x)
	case ast.Service:
		f.encodeService(x)
	default:
		panic("should not reach")
	}
	f.endScope()
}
