package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeFunction(fc *ast.Function) {
	if fc.Oneway != nil {
		f.encodeIdentifier(fc.Oneway)
		f.forward(true, fc.FunctionType.StartPos())
	}
	f.encodeFunctionType(fc.FunctionType)
	f.forward(true, fc.Identifier.Start)
	f.encodeIdentifier(&fc.Identifier)
	f.forwardAndEmptySep(true, fc.LChevron)
	f.encodeChevronFieldList(fc.LChevron, fc.RChevron, fc.FieldList)
	if fc.Throws != nil {
		f.forward(true, fc.Throws.StartPos())
		f.encodeThrows(fc.Throws)
	}
	if fc.Annotations != nil {
		f.encodeAnnotations(fc.Annotations)
	}
}

func (f *Formatter) encodeFunctionType(typ ast.FunctionType) {
	switch x := typ.(type) {
	case ast.FieldType:
		f.encodeFieldType(x)
	case ast.Identifier: // VOID
		f.encodeIdentifier(&x)
	default:
		panic("should not reach")
	}
}

func (f *Formatter) encodeThrows(throws *ast.Throws) {
	f.encodeKeyword(throws.Start, "throws")
	f.forward(true, throws.LChevron)
	f.encodeChevronFieldList(throws.LChevron, throws.RChevron, throws.FieldList)
}
