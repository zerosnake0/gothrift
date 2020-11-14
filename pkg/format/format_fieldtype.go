package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeFieldType(typ ast.FieldType) {
	switch x := typ.(type) {
	case ast.Identifier:
		f.encodeIdentifier(&x)
	case ast.BaseType:
		f.encodeBaseType(x)
	case ast.ContainerType:
		f.encodeContainerType(x)
	default:
		shouldNotReach()
	}
}

func (f *Formatter) encodeBaseType(typ ast.BaseType) {
	f.encodeIdentifier(&typ.Identifier)
	if typ.Annotations != nil {
		f.encodeAnnotations(typ.Annotations)
	}
}

func (f *Formatter) encodeContainerType(typ ast.ContainerType) {
	f.encodeSimpleContainerType(typ.SimpleContainerType)
	if typ.Annotations != nil {
		f.encodeAnnotations(typ.Annotations)
	}
}

func (f *Formatter) encodeSimpleContainerType(typ ast.SimpleContainerType) {
	switch x := typ.(type) {
	case ast.MapType:
		f.encodeMapType(x)
	case ast.SetType:
		f.encodeSetType(x)
	case ast.ListType:
		f.encodeListType(x)
	default:
		shouldNotReach()
	}
}

func (f *Formatter) encodeMapType(typ ast.MapType) {
	f.encodeKeyword(typ.Start, "map")
	f.forwardAndEmptySep(true, typ.LChevron)
	f.encodeKeyword(typ.LChevron, "<")
	f.forwardAndEmptySep(true, typ.Key.StartPos())
	f.encodeFieldType(typ.Key)
	f.forwardAndEmptySep(true, typ.Comma)
	f.encodeKeyword(typ.Comma, ",")
	f.forward(true, typ.Value.StartPos())
	f.encodeFieldType(typ.Value)
	f.forwardAndEmptySep(true, typ.RChevron)
	f.encodeKeyword(typ.RChevron, ">")
}

func (f *Formatter) encodeSetType(typ ast.SetType) {
	f.encodeKeyword(typ.Start, "set")
	f.forwardAndEmptySep(true, typ.LChevron)
	f.encodeKeyword(typ.LChevron, "<")
	f.forwardAndEmptySep(true, typ.Elem.StartPos())
	f.encodeFieldType(typ.Elem)
	f.forwardAndEmptySep(true, typ.RChevron)
	f.encodeKeyword(typ.RChevron, ">")
}

func (f *Formatter) encodeListType(typ ast.ListType) {
	f.encodeKeyword(typ.Start, "list")
	f.forwardAndEmptySep(true, typ.LChevron)
	f.encodeKeyword(typ.LChevron, "<")
	f.forwardAndEmptySep(true, typ.Elem.StartPos())
	f.encodeFieldType(typ.Elem)
	f.forwardAndEmptySep(true, typ.RChevron)
	f.encodeKeyword(typ.RChevron, ">")
}
