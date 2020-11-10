package parser

import (
	"io"
)

// must in capture mode
// returns
// - the current byte at head
// - has digit or not
// - is eof or not
func (lx *exprLexerImpl) trySkipInteger(c byte) (byte, bool, bool) {
	hasDigit := false
	// skip sign
	if c == '+' || c == '-' {
		lx.head++ // capture +-
		c = lx.nextByte()
		if lx.err != nil {
			goto AllowEOF // allow single +-
		}
	}
	if c < '0' || c > '9' {
		return c, false, false
	}
	hasDigit = true
	// try skip integer part
	for i := lx.head + 1; ; i++ {
		if i == lx.tail {
			if lx.readMore(); lx.err != nil {
				lx.head = i   // capture the current integer part
				goto AllowEOF // allow EOF
			}
		}
		c = lx.buffer[i]
		if c < '0' || c > '9' {
			lx.head = i
			return c, hasDigit, false
		}
	}
AllowEOF:
	if lx.err == io.EOF {
		lx.err = nil
		return 0, hasDigit, true
	}
	// lx.err != nil
	return 0, hasDigit, false
}

// do not forward head before calling this
func (lx *exprLexerImpl) readNumber(c byte) (num string) {
	oldCapture := lx.capture
	lx.capture = true

	begin := lx.head
	eof := false
	// skip pure integer
	c, _, eof = lx.trySkipInteger(c)
	if lx.err != nil || eof {
		goto End
	}
	// try skip fractional part
	if c == '.' {
		lx.head++ // we must consume the dot
		for i := lx.head; ; i++ {
			if i == lx.tail {
				if lx.readMore(); lx.err != nil {
					if i == lx.head {
						goto End // we have only one dot, error
					}
					lx.head = i
					goto AllowEOF // allow EOF
				}
			}
			c = lx.buffer[i]
			if c < '0' || c > '9' {
				if i == lx.head {
					lx.err = UnexpectedByteError{
						got: c,
					}
					goto End // we have only one dot, error
				}
				lx.head = i
				break
			}
		}
	}
	// try skip exponential part
	if c == 'e' || c == 'E' {
		oldHead := lx.head // save head location
		lx.head++
		c = lx.nextByte()
		if lx.err != nil {
			lx.head = oldHead
			goto AllowEOF
		}
		_, hasDigit, _ := lx.trySkipInteger(c)
		if lx.err != nil {
			goto End
		}
		// any way we should end wheat eof or not
		if !hasDigit {
			lx.head = oldHead // reset head location
		}
	}
	goto End
AllowEOF:
	if lx.err == io.EOF {
		lx.err = nil
	}
End:
	num = string(lx.buffer[begin:lx.head])
	lx.capture = oldCapture
	return
}
