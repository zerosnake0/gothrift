package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

func TestExprLexerImpl_scanLineComment(t *testing.T) {
	t.Run("line", func(t *testing.T) {
		must, doc := parseAsReader(t, `  // line comment`)

		must.Equal(1, len(doc.Comments))

		cmt := doc.Comments[0]
		lcmt, ok := cmt.(ast.LineComment)
		must.True(ok)

		checkSpan(must, lcmt, 0, 2, 0, 17)
		checkTextSpan(must, lcmt, 0, 2, 0, 17, "// line comment")
	})
	t.Run("line", func(t *testing.T) {
		must, doc := parseAsReader(t, `  // line comment
`)
		must.Equal(1, len(doc.Comments))

		cmt := doc.Comments[0]
		lcmt, ok := cmt.(ast.LineComment)
		must.True(ok)

		checkSpan(must, lcmt, 0, 2, 0, 17)
		must.Equal(" line comment", lcmt.Text)
		checkTextSpan(must, lcmt, 0, 2, 0, 17, "// line comment")
	})
}

func TestExprLexerImpl_scanBlockComment(t *testing.T) {
	t.Run("start with slash", func(t *testing.T) {
		must, doc := parseAsReader(t, ` /*/*/ `)

		must.Equal(1, len(doc.Comments))

		cmt := doc.Comments[0]
		bcmt, ok := cmt.(ast.BlockComment)
		must.True(ok)

		checkSpan(must, bcmt, 0, 1, 0, 6)
		must.Equal("/", bcmt.Text)
		checkTextSpan(must, bcmt, 0, 1, 0, 6, "/*/*/")
	})
	t.Run("change line", func(t *testing.T) {
		must, doc := parseAsReader(t, ` /*
*/ `)
		must.Equal(1, len(doc.Comments))

		cmt := doc.Comments[0]
		bcmt, ok := cmt.(ast.BlockComment)
		must.True(ok)

		checkSpan(must, bcmt, 0, 1, 1, 2)
		must.Equal("\n", bcmt.Text)
		checkTextSpan(must, bcmt, 0, 1, 1, 2, "/*\n*/")
	})
	t.Run("error", func(t *testing.T) {
		must := require.New(t)
		_, err := ParseString(`/*/`)
		must.Error(err)
	})
}
