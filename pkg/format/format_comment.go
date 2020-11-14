package format

import (
	"strings"
	"unicode"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (f *Formatter) outputComments(prev, before *ast.Pos) int {
	if prev != nil {
		prev = prev.Ptr()
	}
	count := 0
	for ; f.cmtIdx < len(f.Doc.Comments); f.cmtIdx++ {
		cmt := f.Doc.Comments[f.cmtIdx]
		start := cmt.StartPos()
		if before != nil {
			if !start.Before(*before) {
				break
			}
		}
		if prev != nil {
			if start.Line > prev.Line {
				break
			}
		}
		f.outputComment(cmt)
		if prev != nil {
			*prev = cmt.EndPos()
		}
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
		shouldNotReach()
	}
}

func (f *Formatter) outputLineComment(cmt ast.LineComment) {
	f.outputIndent()

	txt := strings.TrimRightFunc(cmt.Text, unicode.IsSpace)
	if len(txt) > 0 && txt[0] != ' ' {
		f.print("%s %s\n", cmt.Prefix, txt)
	} else {
		f.print("%s%s\n", cmt.Prefix, txt)
	}

	f.lastEnd = ast.Pos{
		Line: cmt.End.Line + 1,
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
