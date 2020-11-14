package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) encodeConst(cst ast.Const) {
	f.encodeKeyword(cst.Start, "const")
	f.forward(true, cst.FieldType.StartPos())
	f.encodeFieldType(cst.FieldType)
	f.forward(true, cst.Key.Start)
	f.encodeIdentifier(&cst.Key)
	f.forward(true, cst.Eq)
	f.encodeKeyword(cst.Eq, "=")
	f.forward(true, cst.Value.StartPos())
	f.encodeConstValue(cst.Value)
	f.encodeEndSeparator(cst.End)
}

func (f *Formatter) encodeConstValue(value ast.ConstValue) {
	switch x := value.(type) {
	case ast.Number:
		f.encodeNumber(x)
	case ast.Literal:
		f.encodeLiteral(x)
	case ast.Identifier:
		f.encodeIdentifier(&x)
	case ast.ConstList:
		f.encodeConstList(x)
	case ast.ConstMap:
		f.encodeConstMap(x)
	default:
		shouldNotReach()
	}
}

func (f *Formatter) encodeConstList(lst ast.ConstList) {
	end := lst.End
	end.Col--
	f.encodeBracket(lst.Start, end, ",", lst.Content, func(span ast.Span) {
		f.encodeConstListItem(span.(*ast.ConstListItem))
	})
}

func (f *Formatter) encodeConstListItem(item *ast.ConstListItem) {
	f.encodeConstValue(item.ConstValue)
	f.encodeEndSeparator(item.End)
}

func (f *Formatter) encodeConstMap(mp ast.ConstMap) {
	end := mp.EndPos()
	end.Col--
	f.encodeBrace(mp.StartPos(), end, ",", mp.Content, func(span ast.Span) {
		f.encodeConstMapItem(span.(*ast.ConstMapItem))
	})
}

func (f *Formatter) encodeConstMapItem(item *ast.ConstMapItem) {
	f.encodeConstValue(item.Key)
	f.forwardAndEmptySep(true, item.Colon)
	f.encodeKeyword(item.Colon, ":")
	f.forward(true, item.Value.StartPos())
	f.encodeConstValue(item.Value)
	f.encodeEndSeparator(item.End)
}
