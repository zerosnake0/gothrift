package parser

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
	"github.com/zerosnake0/goutil/convert"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

func parseAsReader(t *testing.T, s string) (*require.Assertions, *ast.Document) {
	must := require.New(t)
	doc, err := ParseReader(iotest.OneByteReader(strings.NewReader(s)))
	must.NoError(err)
	must.NotNil(doc)
	return must, doc
}

func checkPos(must *require.Assertions, pos ast.Pos, line, col int) {
	must.Equal(ast.Pos{
		Line: line,
		Col:  col,
	}, pos)
}

func checkSpan(must *require.Assertions, span ast.Span, line1, col1, line2, col2 int) {
	checkPos(must, span.StartPos(), line1, col1)
	checkPos(must, span.EndPos(), line2, col2)
}

func checkTextSpan(must *require.Assertions, span ast.TextSpan, line1, col1, line2, col2 int, txt string) {
	checkSpan(must, span, line1, col1, line2, col2)
	must.Equal(txt, fmt.Sprintf("%s", span))
}

func testLexer(t *testing.T, s string, err error, onNoError func(must *require.Assertions,
	lval *exprSymType, lx *exprLexerImpl)) {
	must := require.New(t)
	lx := defaultLexerPool.borrowLexer()
	defer defaultLexerPool.returnLexer(lx)
	data := convert.LocalStringToBytes(s)
	lx.ResetBytes(data)
	lval := &exprSymType{}
	res := lx.Lex(lval)
	if err != nil {
		if err == io.EOF {
			must.Equal(exprEofCode, res, "%s expecting %v", s, err)
		} else {
			must.Equal(exprErrCode, res, "%s expecting %v", s, err)
		}
	} else {
		onNoError(must, lval, lx)
	}
}
