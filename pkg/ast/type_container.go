package ast

import (
	"fmt"
)

type ContainerType struct {
	SimpleContainerType SimpleContainerType
	Annotations         *Annotations
}

var _ FieldType = ContainerType{}

func (c ContainerType) StartPos() Pos {
	return c.SimpleContainerType.StartPos()
}

func (c ContainerType) EndPos() Pos {
	if c.Annotations == nil {
		return c.SimpleContainerType.EndPos()
	}
	return c.Annotations.End
}

func (c ContainerType) Format(f fmt.State, r rune) {
	panic("implement me")
}

type SimpleContainerType interface {
	TextSpan
}
