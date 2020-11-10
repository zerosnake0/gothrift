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
	f.startBrace(svc.LBrace)
	for i := range svc.FunctionList {
		fc := &svc.FunctionList[i]
		if i > 0 {
			f.print(",")
		}
		f.forward(false, fc.StartPos())
		f.encodeFunction(fc)
	}
	f.endBrace(svc.RBrace)
	if svc.Annotations != nil {
		f.encodeAnnotations(svc.Annotations)
	}
}

func (f *Formatter) encodeExtends(extends *ast.Extends) {
	f.encodeKeyword(extends.Start, "extends")
	f.forward(true, extends.Identifier.Start)
	f.encodeIdentifier(&extends.Identifier)
}
