package parser

import "fmt"

type UnexpectedByteError struct {
	got  byte
	exp  byte
	exp2 byte
}

func (e UnexpectedByteError) Error() string {
	if e.exp == 0 {
		return fmt.Sprintf("unexpected character %q", e.got)
	}
	if e.exp2 == 0 {
		return fmt.Sprintf("expecting %q but got %q", e.exp, e.got)
	}
	return fmt.Sprintf("expecting %q or %q but got %q", e.exp, e.exp2, e.got)
}

type InvalidStringCharError struct {
	c byte
}

func (e InvalidStringCharError) Error() string {
	return fmt.Sprintf("invalid character %x found", e.c)
}

type UnknownIdentifierError struct {
	identifier string
}

func (e UnknownIdentifierError) Error() string {
	return fmt.Sprintf("unknown identifier %q", e.identifier)
}
