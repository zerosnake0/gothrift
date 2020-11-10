package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) startBrace(lbrace ast.Pos) {
	f.encodeKeyword(lbrace, "{")
	f.newScope()
}

func (f *Formatter) endBrace(rbrace ast.Pos) {
	f.forward(false, rbrace)
	f.endScope()
	f.encodeKeyword(rbrace, "}")
}

func (f *Formatter) startChevron(lchevron ast.Pos) {
	f.encodeKeyword(lchevron, "(")
	f.newScope()
}

func (f *Formatter) endChevron(rchevron ast.Pos) {
	f.forwardAndEmptySep(false, rchevron)
	f.endScope()
	f.encodeKeyword(rchevron, ")")
}
