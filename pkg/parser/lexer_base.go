package parser

import (
	"bytes"
	"fmt"
	"io"
	"runtime"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

func (lx *exprLexerImpl) reset() {
	lx.reader = nil
	lx.buffer = nil
	lx.tail = 0

	lx.exprLexerImplEmbedded = exprLexerImplEmbedded{}
}

func (lx *exprLexerImpl) Reset(r io.Reader) {
	switch v := r.(type) {
	case nil:
		lx.reset()
		return
	case *bytes.Buffer:
		lx.ResetBytes(v.Bytes())
		return
	}
	lx.reader = r
	lx.buffer = lx.fixbuf[:cap(lx.fixbuf)]
	lx.tail = 0

	lx.exprLexerImplEmbedded = exprLexerImplEmbedded{}
}

func (lx *exprLexerImpl) ResetBytes(data []byte) {
	lx.reader = nil
	lx.buffer = data
	lx.tail = len(data)

	lx.exprLexerImplEmbedded = exprLexerImplEmbedded{}
}

func (lx *exprLexerImpl) readMore() {
	if lx.reader == nil {
		lx.err = io.EOF
		return
	}
	var (
		n   int
		err error
	)
	for {
		if lx.capture {
			var buf [bufferSize]byte
			n, err = lx.reader.Read(buf[:])
			lx.buffer = append(lx.buffer[:lx.tail], buf[:n]...)
			lx.tail += n
			// save internal buffer for reuse
			lx.fixbuf = lx.buffer
		} else {
			if debug {
				if lx.head != lx.tail {
					panic(fmt.Errorf("head %d, tail %d", lx.head, lx.tail))
				}
			}
			n, err = lx.reader.Read(lx.buffer)
			lx.offset += lx.tail
			lx.head = 0
			lx.tail = n
		}
		if err != nil {
			if err == io.EOF && n > 0 {
				return
			}
			lx.err = err
			return
		}
		if n > 0 {
			return
		}
		// n == 0 && err == nil
		// the implementation of the reader is wrong
		runtime.Gosched()
	}
}

func (lx *exprLexerImpl) Buffer() []byte {
	return lx.buffer[lx.head:lx.tail]
}

// charHead: the head location of the newline character
func (lx *exprLexerImpl) onNewLine(charHead int) {
	lx.lineNo++
	lx.lineBegOffset = lx.offset + charHead + 1
}

// the offset should be sure that it is on the same line
func (lx *exprLexerImpl) pos(offset int) ast.Pos {
	return ast.Pos{
		Line: lx.lineNo,
		Col:  lx.offset + lx.head + offset - lx.lineBegOffset,
	}
}

func (lx *exprLexerImpl) nextToken() (ret byte) {
	for {
		i := lx.head
		for ; i < lx.tail; i++ {
			c := lx.buffer[i]
			if c <= ' ' {
				if valueTypeMap[c] == WhiteSpaceValue {
					if c == '\n' {
						lx.onNewLine(i)
					}
					continue
				}
			}
			lx.head = i
			return c
		}
		lx.head = i
		if lx.readMore(); lx.err != nil {
			return
		}
	}
}

func (lx *exprLexerImpl) expectBytes(s string) {
	last := len(s) - 1
	j := 0
	for {
		i := lx.head
		for ; i < lx.tail; i++ {
			c := lx.buffer[i]
			if c != s[j] {
				lx.err = UnexpectedByteError{exp: s[j], got: c}
				return
			}
			if c == '\n' {
				lx.onNewLine(i)
			}
			if j == last {
				lx.head = i + 1
				return
			}
			j++
		}
		lx.head = i
		if lx.readMore(); lx.err != nil {
			return
		}
	}
}

func (lx *exprLexerImpl) nextByte() (ret byte) {
	if lx.head == lx.tail {
		if lx.readMore(); lx.err != nil {
			return
		}
	}
	return lx.buffer[lx.head]
}
