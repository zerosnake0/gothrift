package ast

import (
	"fmt"
)

type Pos struct {
	Line int
	Col  int
}

func (p Pos) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "(%d,%d)", p.Line, p.Col)
}

func (p Pos) Before(p2 Pos) bool {
	if p.Line != p2.Line {
		return p.Line < p2.Line
	}
	return p.Col < p2.Col
}

func (p Pos) Ptr() *Pos {
	return &p
}

type Span interface {
	StartPos() Pos
	EndPos() Pos
}

type Node struct {
	Start Pos
	End   Pos
}

func (n Node) StartPos() Pos {
	return n.Start
}

func (n Node) EndPos() Pos {
	return n.End
}

type TextSpan interface {
	Span
	fmt.Formatter
}

type TextNode struct {
	Node
	Text string
}

func (n TextNode) Content() string {
	return n.Text
}

func (n TextNode) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		fmt.Fprintf(f, "<%s,%s>%s", n.Start, n.End, n.Text)
	default:
		fmt.Fprintf(f, "%s", n.Text)
	}
}
