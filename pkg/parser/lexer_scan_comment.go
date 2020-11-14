package parser

import (
	"io"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

// the first '//' should already be consumed
func (lx *exprLexerImpl) scanLineComment(prefix string) {
	oldCapture := lx.capture
	lx.capture = true

	begin := lx.head

	cmt := ast.LineComment{
		Prefix: prefix,
	}
	cmt.Start = lx.pos(-2)
	for {
		i := lx.head
		for ; i < lx.tail; i++ {
			c := lx.buffer[i]
			if c != '\n' {
				continue
			}
			// c == '\n'
			cmt.Text = string(lx.buffer[begin:i])
			cmt.End = lx.pos(0)
			lx.onNewLine(i)
			lx.head = i + 1
			goto End
		}
		lx.head = i
		if lx.readMore(); lx.err != nil {
			if lx.err == io.EOF {
				cmt.Text = string(lx.buffer[begin:i])
				cmt.End = lx.pos(0)
				goto End
			}
			return
		}
	}
End:
	lx.comments = append(lx.comments, cmt)
	lx.capture = oldCapture
	return
}

// the first '/*' should already be consumed
func (lx *exprLexerImpl) scanBlockComment() {
	oldCapture := lx.capture
	lx.capture = true

	begin := lx.head
	var cmt ast.BlockComment
	cmt.Start = lx.pos(-2)
	for {
		i := lx.head
		for ; i < lx.tail; i++ {
			c := lx.buffer[i]
			switch c {
			case '\n':
				lx.onNewLine(i)
			case '/':
				if i != begin && lx.buffer[i-1] == '*' {
					// we are in capture mode so this is safe
					lx.head = i + 1
					cmt.Text = string(lx.buffer[begin : i-1])
					cmt.End = lx.pos(0)
					lx.comments = append(lx.comments, cmt)
					goto End
				}
			}
		}
		lx.head = i
		if lx.readMore(); lx.err != nil {
			goto End
		}
	}
End:
	lx.capture = oldCapture
	return
}
