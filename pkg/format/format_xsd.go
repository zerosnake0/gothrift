package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeXsdAttrs(attrs *ast.XsdAttrs) {
	f.encodeKeyword(attrs.Start, "xsd_attrs")
	f.forward(true, attrs.LBrace)
	f.encodeBraceFieldList(attrs.LBrace, attrs.RBrace, attrs.FieldList)
}
