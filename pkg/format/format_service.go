package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeService(svc ast.Service) {
	f.encodeKeyword(svc.Start, "service")
	f.forward(true, svc.Identifier.Start)
	f.encodeIdentifier(&svc.Identifier)
	if svc.Extends != nil {
		f.forward(true, svc.Extends.Start)
		f.encodeExtends(svc.Extends)
	}
	f.forward(true, svc.LBrace)
	f.encodeBrace(svc.LBrace, svc.RBrace, "", svc.FunctionList, func(span ast.Span) {
		f.encodeFunction(span.(*ast.Function))
	})
	if svc.Annotations != nil {
		f.encodeAnnotations(svc.Annotations)
	}
}

func (f *Formatter) encodeExtends(extends *ast.Extends) {
	f.encodeKeyword(extends.Start, "extends")
	f.forward(true, extends.Identifier.Start)
	f.encodeIdentifier(&extends.Identifier)
}
