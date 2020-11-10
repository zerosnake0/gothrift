package ast

import (
	"fmt"
)

// Annotations
type Annotations struct {
	Node
	Annotations []Annotation
}

var _ TextSpan = Annotations{}

func (a Annotations) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "( ")
	if len(a.Annotations) > 0 {
		for i, anno := range a.Annotations {
			if i > 0 {
				fmt.Fprintf(f, ", ")
			}
			fmt.Fprintf(f, "%s", anno)
		}
	}
	fmt.Fprintf(f, " )")
}

// Annotation
type Annotation struct {
	Key   Identifier
	Value AnnotationValue
	End   *Identifier
}

var _ TextSpan = Annotation{}

func (a Annotation) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "%s %s", a.Key, a.Value)
}

func (a Annotation) StartPos() Pos {
	return a.Key.Start
}

func (a Annotation) EndPos() Pos {
	if a.End != nil {
		return a.End.End
	}
	return a.Value.Value.End
}

// AnnotationValue
type AnnotationValue struct {
	Start Pos
	Value Literal
}

var _ TextSpan = AnnotationValue{}

func (v AnnotationValue) StartPos() Pos {
	return v.Start
}

func (v AnnotationValue) EndPos() Pos {
	return v.Value.End
}

func (v AnnotationValue) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, `= %s`, v.Value)
}
