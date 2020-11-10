package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeXsdAttrs(attrs *ast.XsdAttrs) {
	f.encodeKeyword(attrs.Start, "xsd_attrs")
	f.forward(true, attrs.LBrace)
	f.startBrace(attrs.LBrace)
	for i := range attrs.FieldList {
		field := &attrs.FieldList[i]
		if i > 0 {
			f.print(",")
		}
		f.forward(false, field.StartPos())
		f.encodeField(field)
	}
	f.endBrace(attrs.RBrace)
}
