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
		panic("should not reach")
	}
}

func (f *Formatter) encodeConstList(lst ast.ConstList) {
	f.encodeKeyword(lst.Start, "[")
	f.newScope()
	for i := range lst.Content {
		item := &lst.Content[i]
		if i == 0 {
			f.forwardAndEmptySep(false, item.StartPos())
		} else {
			f.print(",")
			f.forward(false, item.StartPos())
		}
		f.encodeConstListItem(item)
	}
	f.forwardAndEmptySep(false, lst.End)
	f.endScope()
	f.encodeKeywordEnd(lst.End, "]")
}

func (f *Formatter) encodeConstListItem(item *ast.ConstListItem) {
	f.encodeConstValue(item.ConstValue)
}

func (f *Formatter) encodeConstMap(mp ast.ConstMap) {
	f.startBrace(mp.Start)
	for i := range mp.Content {
		item := &mp.Content[i]
		if i == 0 {
			f.forwardAndEmptySep(false, item.StartPos())
		} else {
			f.print(",")
			f.forward(false, item.StartPos())
		}
		f.encodeConstMapItem(item)
	}
	f.endBrace(mp.End)
}

func (f *Formatter) encodeConstMapItem(item *ast.ConstMapItem) {
	f.encodeConstValue(item.Key)
	f.forwardAndEmptySep(true, item.Colon)
	f.encodeKeyword(item.Colon, ":")
	f.forward(true, item.Value.StartPos())
	f.encodeConstValue(item.Value)
}
