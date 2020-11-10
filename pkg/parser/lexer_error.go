package parser

import (
	"fmt"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

type LexerError struct {
	pos ast.Pos
	buf string
	ex  string
	err error
}

func (e LexerError) Error() string {
	if e.err == nil {
		return fmt.Sprintf("%s near %s %s", e.ex, e.pos, e.buf)
	}
	return fmt.Sprintf("%s (%s) near %s %s", e.ex, e.err, e.pos, e.buf)
}
