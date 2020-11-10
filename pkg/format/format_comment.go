package format

import (
	"strings"
	"unicode"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) outputComments(before ast.Pos) int {
	count := 0
	for ; f.cmtIdx < len(f.Doc.Comments); f.cmtIdx++ {
		cmt := f.Doc.Comments[f.cmtIdx]
		if !cmt.StartPos().Before(before) {
			break
		}
		f.outputComment(cmt)
		count++
	}
	return count
}

func (f *Formatter) outputComment(cmt ast.Comment) {
	start := cmt.StartPos()
	f.outputBlankLine(start)

	switch x := cmt.(type) {
	case ast.LineComment:
		f.outputLineComment(x)
	case ast.BlockComment:
		f.outputBlockComment(x)
	default:
		panic("should not reach")
	}
}

func (f *Formatter) outputLineComment(cmt ast.LineComment) {
	f.outputIndent()

	txt := strings.TrimRightFunc(cmt.Text, unicode.IsSpace)
	if len(txt) > 0 && txt[0] != ' ' {
		f.print("// %s\n", txt)
	} else {
		f.print("//%s\n", txt)
	}

	end := cmt.EndPos()
	f.lastEnd = ast.Pos{
		Line: end.Line + 1, // new line
	}
}

func (f *Formatter) outputBlockComment(cmt ast.BlockComment) {
	txt := cmt.Text
	for count := 0; ; count++ {
		// index of new line
		inl := strings.IndexByte(txt, '\n')
		var line string
		if inl < 0 { // no more newline
			line = txt
		} else {
			line = txt[:inl]
			txt = txt[inl+1:]
		}
		line = strings.TrimRightFunc(line, unicode.IsSpace)

		f.outputIndent()
		if count == 0 { // first line
			f.print("/*")
			if len(line) > 0 && line[0] != ' ' {
				f.print(" %s", line)
			} else { // empty line or start with space
				f.print("%s", line)
			}
			if inl < 0 { // single line block comment
				f.print(" */")
				break
			}
		} else {
			f.print(" *")
			if len(line) == 0 { // empty line
				if inl < 0 {
					f.print("/")
					break
				}
			} else { // non empty line
				i := strings.Index(line, " *")
				if i >= 0 {
					f.print(" %s", line[i+2:])
				} else {
					f.print(" %s", line)
				}
				if inl < 0 { // last line
					f.print(" */")
					break
				}
			}
		}
		f.print("\n")
		f.lastEnd = ast.Pos{
			Line: f.lastEnd.Line + 1,
		} // move to new line
	}
	f.lastEnd = cmt.EndPos()
}
