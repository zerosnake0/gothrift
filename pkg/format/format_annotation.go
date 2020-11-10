package format

import (
	"github.com/zerosnake0/gothrift/pkg/ast"
)

// will forward with true, no need to forward before calling this
func (f *Formatter) encodeAnnotations(annotations *ast.Annotations) {
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
	f.startChevron(annotations.Start)
	for i := range annotations.Annotations {
		anno := &annotations.Annotations[i]
		if i == 0 {
			f.forwardAndEmptySep(false, anno.StartPos())
		} else {
			f.print(",")
			f.forward(false, anno.StartPos())
		}
		f.encodeAnnotation(anno)
	}
	f.endChevron(annotations.End)
End:
	f.lastEnd = annotations.End
}

func (f *Formatter) encodeAnnotation(anno *ast.Annotation) {
	f.encodeIdentifier(&anno.Key)
	f.forward(true, anno.Value.Start)
	f.encodeAnnotationValue(anno.Value)
}

func (f *Formatter) encodeAnnotationValue(annoValue ast.AnnotationValue) {
	f.encodeKeyword(annoValue.Start, "=")
	f.forward(true, annoValue.Start)
	f.encodeLiteral(annoValue.Value)
}
