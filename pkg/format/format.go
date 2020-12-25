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
	col     int     // current output cursor
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
	data := f.buf.Bytes()
	for {
		i := bytes.IndexByte(data, '\n')
		data = data[i+1:]
		if i < 0 {
			f.col += len(data)
			break
		}
		f.col = 0
		f.pushCurrentIndent()
	}
	io.Copy(f.Writer, &f.buf)
}

func (f *Formatter) printWithIndent(format string, args ...interface{}) {
	f.printIndent()
	f.print(format, args...)
}

func (f *Formatter) printIndent() {
	if f.col != 0 {
		shouldNotReach()
	}
	for _, elem := range f.stack {
		if elem.indent != 0 {
			f.Writer.Write([]byte{'\t'})
			f.col++
		}
	}
}

// returns: if a new line is printed
func (f *Formatter) newLineIfNot() bool {
	if f.col == 0 {
		return false
	}
	f.print("\n")
	return true
}

const defaultSep = " "

func (f *Formatter) Encode() {
	f.formatterEmbedded = formatterEmbedded{}
	f.stack = f.stack[:0]
	f.sep = defaultSep
	f.encodeHeaders()
	f.encodeDefinitions()
	f.outputComments(nil, nil)
	f.newLineIfNot()
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

func (f *Formatter) pushStack() {
	f.stack = append(f.stack, stack{
		first: true,
	})
}

func (f *Formatter) popStack() {
	f.stack = f.stack[:len(f.stack)-1]
}

func (f *Formatter) endScope() {
	f.popStack()
}

func (f *Formatter) outputRemainingComments(next ast.Pos) {
	if f.lastEnd.Col > 0 {
		f.pushStack() // treat the following comment as a single scope
		f.outputComments(&f.lastEnd, &next)
		f.popStack()
	}
}

func (f *Formatter) encodeEndSeparator(sep *ast.Identifier) {
	if sep != nil {
		f.outputRemainingComments(sep.Start)
		f.lastEnd = sep.End
	}
}

func (f *Formatter) newScope(next ast.Pos, attach bool) {
	// consume all the remaining comments
	if attach {
		f.sep = ""
	}
	last := f.lastEnd
	f.outputRemainingComments(next)
	f.sep = defaultSep

	atRoot := len(f.stack) == 0
	if atRoot { // if we are at the root level
		f.outputComments(nil, &next)
		f.newLineIfNot()
		f.outputBlankLine(next)
	}

	// real new scope here
	f.pushStack()
	if !atRoot {
		// we were already at a non root level
		// we need to indent if line change happened
		if f.lastEnd.Line != last.Line {
			f.pushCurrentIndent()
		}
	}
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
	// case1: f.col == 0, f.lastEnd.col > 0, d == 1 no need to output
	if f.col == 0 {
		// we printed a newline already
		if last.Col == 0 {
			// there is really a gap to be printed
			f.print("\n")
		} else {
			if d > 1 {
				f.print("\n")
			}
		}
	} else {
		// case2: f.col > 0, f.lastEnd.col > 0, d > 1
		f.print("\n")
		if d > 1 {
			// the gap is more than one line
			f.print("\n")
		}
	}
	f.lastEnd = ast.Pos{
		Line: next.Line,
	}
}

func (f *Formatter) outputIndent() {
	if f.col == 0 {
		f.printIndent()
	} else {
		f.print(f.sep)
	}
	f.sep = defaultSep
}

func (f *Formatter) forward(follow bool, next ast.Pos) (count int) {
	last := f.lastEnd // save the previous token's end position
	count = f.outputComments(nil, &next)
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
		if f.lastEnd.Line == last.Line || // comment didn't change line
			f.lastEnd.Line == next.Line { // next same line with comment end
			return
		}
		// comment changed line and the next element
		// is not on the same line as the comment
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
