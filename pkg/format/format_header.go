package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeHeaders() {
	for _, header := range f.Doc.Headers {
		f.encodeHeader(header)
	}
}

func (f *Formatter) encodeHeader(header ast.Header) {
	start := header.StartPos()
	f.newScope(start)
	f.forward(false, start)
	switch x := header.(type) {
	case ast.Include:
		f.encodeInclude(x)
	case ast.Namespace:
		f.encodeNamespace(x)
	default:
		shouldNotReach()
	}
	f.endScope()
}

func (f *Formatter) encodeInclude(inc ast.Include) {
	f.encodeIdentifier(&inc.Keyword)
	f.forward(true, inc.Literal.Start)
	f.encodeLiteral(inc.Literal)
}

func (f *Formatter) encodeNamespace(ns ast.Namespace) {
	f.encodeKeyword(ns.Start, "namespace")
	f.forward(true, ns.Language.Start)
	f.encodeIdentifier(&ns.Language)
	f.forward(true, ns.Namespace.Start)
	f.encodeIdentifier(&ns.Namespace)
	if ns.Annotations != nil {
		f.encodeAnnotations(ns.Annotations)
	}
}
