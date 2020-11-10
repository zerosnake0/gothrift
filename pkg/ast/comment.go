package ast

import (
	"fmt"
)

type Comment interface {
	TextSpan
}

type LineComment struct {
	TextNode
}

var _ Comment = LineComment{}

func (cmt LineComment) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "//%s", cmt.Text)
}

type BlockComment struct {
	TextNode
}

var _ Comment = BlockComment{}

func (cmt BlockComment) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "/*%s*/", cmt.Text)
}
