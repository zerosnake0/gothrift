package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeChevronFieldList(l, r ast.Pos, fields []ast.Field) {
	f.startChevron(l)
	for i := range fields {
		field := &fields[i]
		if i == 0 {
			f.forwardAndEmptySep(false, field.StartPos())
		} else {
			f.print(",")
			f.forward(false, field.StartPos())
		}
		f.encodeField(field)
	}
	f.endChevron(r)
}

func (f *Formatter) encodeBraceFieldList(l, r ast.Pos, fields []ast.Field) {
	f.startBrace(l)
	for i := range fields {
		field := &fields[i]
		if i > 0 {
			f.print(",")
		}
		f.forward(false, field.StartPos())
		f.encodeField(field)
	}
	f.endBrace(r)
}
