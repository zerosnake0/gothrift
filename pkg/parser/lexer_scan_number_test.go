package parser

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func testReadNumber(t *testing.T, s string, err error, num string) {
	testLexer(t, s, err, func(must *require.Assertions, lval *exprSymType, lx *exprLexerImpl) {
		must.Equal(lval.number.Text, num, "%s expecting %v %s", s, err, num)
		must.Equal(s[len(num):], string(lx.Buffer()), "%s expecting %v %s", s, err, num)
	})
}

func TestExprLexerImpl_readNumber(t *testing.T) {
	// test cases
	// - whitespace
	// - sig
	// - digit
	// - dot
	// - exp
	// - other
	t.Run("sig part", func(t *testing.T) {
		testReadNumber(t, "+", nil, "+")
		testReadNumber(t, "-", nil, "-")
		for _, c := range " +-aeE" {
			testReadNumber(t, "+"+string(c), nil, "+")
		}
		testReadNumber(t, "+1", nil, "+1")
		testReadNumber(t, "+.", io.EOF, "")
	})
	t.Run("integer part", func(t *testing.T) {
		testReadNumber(t, "12", nil, "12")
		testReadNumber(t, "+12", nil, "+12")
		for _, c := range " +-eEa" {
			testReadNumber(t, "-1"+string(c), nil, "-1")
		}
		testReadNumber(t, "1.", io.EOF, "")
	})
	t.Run("fraction part", func(t *testing.T) {
		for _, prefix := range []string{
			"", "+", "-1",
		} {
			prefix += "." // .
			testReadNumber(t, prefix, io.EOF, "")
			for _, c := range " +-.eEa" {
				testReadNumber(t, prefix+string(c), UnexpectedByteError{got: byte(c)}, "")
			}

			prefix += "1" // .1
			testReadNumber(t, prefix, nil, prefix)
			for _, c := range " +-.eEa" {
				testReadNumber(t, prefix+string(c), nil, prefix)
			}
		}
	})
	t.Run("exp part", func(t *testing.T) {
		for _, prefix := range []string{
			"+", "-", // only sign
			"-1",        // no fraction
			".1", "1.2", // with fraction
		} {
			testReadNumber(t, prefix+"e", nil, prefix)
			testReadNumber(t, prefix+"e+", nil, prefix)
			testReadNumber(t, prefix+"e+1", nil, prefix+"e+1")
		}

		//for _, prefix := range []string{
		//	"+", "-1", ".1", "1.2",
		//} {
		//
		//}
	})
}
