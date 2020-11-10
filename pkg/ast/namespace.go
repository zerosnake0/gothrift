package ast

import (
	"fmt"
)

type Namespace struct {
	Start       Pos
	Language    Identifier
	Namespace   Identifier
	Annotations *Annotations
}

var _ Header = Namespace{}

func (ns Namespace) StartPos() Pos {
	return ns.Start
}

func (ns Namespace) EndPos() Pos {
	if ns.Annotations == nil {
		return ns.Namespace.End
	}
	return ns.Annotations.End
}

func (ns Namespace) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "namespace %s %s", ns.Language, ns.Namespace)
	if ns.Annotations != nil && len(ns.Annotations.Annotations) > 0 {
		fmt.Fprintf(f, " %s", ns.Annotations)
	}
}
