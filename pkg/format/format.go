package format

import (
	"bytes"
	"fmt"
	"io"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

type stack struct {
	first  bool
	indent int
}

type formatterEmbedded struct {
	cmtIdx  int     // current comment index
	lastEnd ast.Pos // end position of last element
}

type Formatter struct {
	Doc    *ast.Document
	Writer io.Writer

	formatterEmbedded

	stack []stack
	buf   bytes.Buffer
	sep   string
}

func (f *Formatter) print(format string, args ...interface{}) {
	f.buf.Reset()
	fmt.Fprintf(&f.buf, format, args...)
	if bytes.IndexByte(f.buf.Bytes(), '\n') >= 0 {
		f.pushCurrentIndent()
	}
	io.Copy(f.Writer, &f.buf)
}

func (f *Formatter) printWithIndent(format string, args ...interface{}) {
	f.printIndent()
	f.print(format, args...)
}

func (f *Formatter) printIndent() {
	for _, elem := range f.stack {
		if elem.indent != 0 {
			f.Writer.Write([]byte{'\t'})
		}
	}
}

const defaultSep = " "

func (f *Formatter) Encode() {
	f.formatterEmbedded = formatterEmbedded{}
	f.stack = f.stack[:0]
	f.sep = defaultSep
	f.encodeHeaders()
	f.encodeDefinitions()
}

func (f *Formatter) pushStackStep() {
	f.stack[len(f.stack)-1].first = false
}

func (f *Formatter) pushCurrentIndent() {
	l := len(f.stack)
	if l == 0 {
		return
	}
	l -= 1
	if f.stack[l].indent == 0 {
		f.stack[l].indent = 1
	}
}

func (f *Formatter) newScope() {
	if len(f.stack) == 0 {
		f.print("\n")
		f.lastEnd.Col = 0 // TODO: is this safe?
	}
	f.stack = append(f.stack, stack{
		first: true,
	})
}

func (f *Formatter) endScope() {
	f.stack = f.stack[:len(f.stack)-1]
}

func (f *Formatter) outputBlankLine(next ast.Pos) {
	last := f.lastEnd
	if last.Line == 0 && last.Col == 0 { // beginning of the doc, nothing to do
		return
	}
	d := next.Line - last.Line
	if d <= 0 {
		return // same line, nothing to do
	}
	f.print("\n")
	f.lastEnd = ast.Pos{
		Line: next.Line,
	}
}

func (f *Formatter) outputIndent() {
	if f.lastEnd.Col == 0 {
		f.printIndent()
	} else {
		f.print(f.sep)
	}
	f.sep = defaultSep
}

func (f *Formatter) forward(follow bool, next ast.Pos) (count int) {
	last := f.lastEnd // save the previous token's end position
	count = f.outputComments(next)
	// if there is no comment, all is good, otherwise
	// 1. line comment, the indent will be pushed, no problem
	// 2. block comment
	//   - <prev> /**/ <next>
	//   - <prev> /*
	//     */ <next>
	//   - <prev>
	//     /* */
	//     <next>
	if follow {
		if f.lastEnd.Line != last.Line { // comment changed line
			if f.lastEnd.Line == next.Line { // next same line with comment end
				return
			}
		}
	}
	f.outputBlankLine(next)
	return
}

func (f *Formatter) forwardAndEmptySep(follow bool, next ast.Pos) {
	count := f.forward(follow, next)
	if count == 0 {
		f.sep = ""
	}
}

func (f *Formatter) encodeKeyword(start ast.Pos, kw string) {
	f.outputIndent()
	f.print("%s", kw)
	start.Col += len(kw)
	f.lastEnd = start
}

func (f *Formatter) encodeKeywordEnd(end ast.Pos, kw string) {
	f.outputIndent()
	f.print("%s", kw)
	f.lastEnd = end
}

func (f *Formatter) encodeIdentifier(ident *ast.Identifier) {
	f.outputIndent()
	f.print("%s", ident.Text)
	f.lastEnd = ident.End
}

func (f *Formatter) encodeLiteral(literal ast.Literal) {
	f.outputIndent()
	f.print("%s", literal.Text)
	f.lastEnd = literal.End
}

func (f *Formatter) encodeNumber(num ast.Number) {
	f.outputIndent()
	f.print("%s", num.Text)
	f.lastEnd = num.End
}
