package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeSenum(senum ast.Senum) {
	f.encodeKeyword(senum.Start, "senum")
	f.forward(true, senum.Identifier.Start)
	f.encodeIdentifier(&senum.Identifier)
	f.forward(true, senum.LBrace)
	f.encodeBrace(senum.LBrace, senum.RBrace, ",", senum.List, func(span ast.Span) {
		f.encodeSenumDef(span.(*ast.SenumDef))
	})
	if senum.Annotations != nil {
		f.encodeAnnotations(senum.Annotations)
	}
}

func (f *Formatter) encodeSenumDef(item *ast.SenumDef) {
	f.encodeLiteral(item.Literal)
	f.encodeEndSeparator(item.End)
}
