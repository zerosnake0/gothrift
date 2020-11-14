package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

// will forward with true, no need to forward before calling this
func (f *Formatter) encodeAnnotations(annotations *ast.Annotations) {
	end := annotations.End
	if len(annotations.Annotations) == 0 {
		if f.cmtIdx >= len(f.Doc.Comments) {
			// no more comment
			goto End
		}
		cmt := f.Doc.Comments[f.cmtIdx]
		if !cmt.StartPos().Before(annotations.End) {
			goto End
		}
		// there is comment between two chevrons
	}
	f.forward(true, annotations.Start)
	end.Col--
	f.encodeChevron(annotations.Start, end, ",", annotations.Annotations,
		func(span ast.Span) {
			f.encodeAnnotation(span.(*ast.Annotation))
		})
End:
	f.lastEnd = annotations.End
}

func (f *Formatter) encodeAnnotation(anno *ast.Annotation) {
	f.encodeIdentifier(&anno.Key)
	f.forward(true, anno.Value.Start)
	f.encodeAnnotationValue(anno.Value)
	f.encodeEndSeparator(anno.End)
}

func (f *Formatter) encodeAnnotationValue(annoValue ast.AnnotationValue) {
	f.encodeKeyword(annoValue.Start, "=")
	f.forward(true, annoValue.Start)
	f.encodeLiteral(annoValue.Value)
}
