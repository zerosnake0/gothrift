package ast

import (
	"fmt"
)

type Literal struct {
	TextNode
}

func (l Literal) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "%s", l.Text)
}
