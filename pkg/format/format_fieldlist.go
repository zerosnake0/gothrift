package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeChevronFieldList(l, r ast.Pos, fields []ast.Field) {
	f.encodeChevron(l, r, ",", fields, func(span ast.Span) {
		f.encodeField(span.(*ast.Field))
	})
}

func (f *Formatter) encodeBraceFieldList(l, r ast.Pos, fields []ast.Field) {
	f.encodeBrace(l, r, ",", fields, func(span ast.Span) {
		f.encodeField(span.(*ast.Field))
	})
}
