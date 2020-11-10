package ast

type Identifier struct {
	TextNode
}

var _ FieldType = Identifier{}
var _ EnumValue = Identifier{}
var _ FunctionType = Identifier{}
