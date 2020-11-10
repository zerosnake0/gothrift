//go:generate goyacc -o expr.go -p "expr" expr.y

package parser

import (
	"io"

	"github.com/zerosnake0/goutil/convert"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

func parse(lx *exprLexerImpl) (*ast.Document, error) {
	res := exprParse(lx)
	if res != 0 {
		return nil, lx.err
	}
	lx.Document.Comments = lx.comments
	return lx.Document, nil
}

func Parse(data []byte) (*ast.Document, error) {
	lx := defaultLexerPool.borrowLexer()
	defer defaultLexerPool.returnLexer(lx)
	lx.ResetBytes(data)
	return parse(lx)
}

func ParseString(s string) (*ast.Document, error) {
	return Parse(convert.LocalStringToBytes(s))
}

func ParseReader(r io.Reader) (*ast.Document, error) {
	lx := defaultLexerPool.borrowLexer()
	defer defaultLexerPool.returnLexer(lx)
	lx.Reset(r)
	return parse(lx)
}
