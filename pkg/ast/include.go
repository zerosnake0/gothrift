package ast

import (
	"fmt"
)

type Include struct {
	Keyword Identifier
	Literal Literal
}

var _ Header = Include{}

func (inc Include) StartPos() Pos {
	return inc.Keyword.Start
}

func (inc Include) EndPos() Pos {
	return inc.Literal.End
}

func (inc Include) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "%s %s", inc.Keyword, inc.Literal)
}
