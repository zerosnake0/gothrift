package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeSenum(senum ast.Senum) {
	f.encodeKeyword(senum.Start, "senum")
	f.forward(true, senum.Identifier.Start)
	f.encodeIdentifier(&senum.Identifier)
	f.forward(true, senum.LBrace)
	f.startBrace(senum.LBrace)
	for i := range senum.List {
		item := &senum.List[i]
		if i > 0 {
			f.print(",")
		}
		f.forward(false, item.StartPos())
		f.encodeSenumDef(item)
	}
	f.endBrace(senum.RBrace)
	if senum.Annotations != nil {
		f.encodeAnnotations(senum.Annotations)
	}
}

func (f *Formatter) encodeSenumDef(item *ast.SenumDef) {
	f.encodeLiteral(item.Literal)
}
