package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeException(exc ast.Exception) {
	f.encodeKeyword(exc.Start, "exception")
	f.forward(true, exc.Identifier.Start)
	f.encodeIdentifier(&exc.Identifier)
	f.forward(true, exc.LBrace)
	f.encodeBraceFieldList(exc.LBrace, exc.RBrace, exc.FieldList)
	if exc.Annotations != nil {
		f.encodeAnnotations(exc.Annotations)
	}
}
