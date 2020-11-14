package parser

import (
	"io"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

var _ exprLexer = &exprLexerImpl{}

func setExprValDoc(exprlex exprLexer, document ast.Document) {
	exprlex.(*exprLexerImpl).Document = &document
}

const errWidth = 20

func (lx *exprLexerImpl) Error(ex string) {
	buf := lx.Buffer()
	var s string
	if len(buf) < errWidth {
		s = string(buf)
	} else {
		s = string(buf[:errWidth])
	}
	lx.err = LexerError{
		pos: lx.pos(0),
		buf: s,
		ex:  ex,
		err: lx.err,
	}
}

var keywords = map[string]int{
	"include":       INCLUDE,
	"cpp_include":   INCLUDE,
	"namespace":     NAMESPACE,
	"const":         CONST,
	"map":           MAP,
	"set":           SET,
	"list":          LIST,
	"cpp_type":      CPPTYPE,
	"typedef":       TYPEDEF,
	"enum":          ENUM,
	"senum":         SENUM,
	"struct":        STRUCTHEAD,
	"union":         STRUCTHEAD,
	"required":      REQUIRED,
	"optional":      OPTIONAL,
	"xsd_all":       XSDALL,
	"xsd_optional":  XSDOPTIONAL,
	"xsd_nillable":  XSDNILLABLE,
	"xsd_namespace": XSDNAMESPACE,
	"xsd_attrs":     XSDATTRS,
	"exception":     EXCEPTION,
	"service":       SERVICE,
	"extends":       EXTENDS,
	"oneway":        ONEWAY,
	"void":          VOID,
	"throws":        THROWS,
}

func init() {
	for _, s := range []string{
		"bool", "byte", "i8", "i16", "i32", "i64",
		"double", "string", "binary", "slist",
	} {
		keywords[s] = BASETYPE
	}
}

func (lx *exprLexerImpl) Lex(lval *exprSymType) int {
	eofCode := exprEofCode // error code when EOF encountered
	if lx.err != nil {
		goto AllowEOF
	}
	for {
		c := lx.nextToken()
		if lx.err != nil {
			goto AllowEOF
		}
		switch c {
		case '#':
			lx.head++
			lx.scanLineComment("#")
			if lx.err != nil {
				goto AllowEOF
			}
		case '/':
			lx.head++
			if c = lx.nextByte(); lx.err != nil {
				goto Err
			}
			switch c {
			case '/':
				lx.head++
				lx.scanLineComment("//")
				if lx.err != nil {
					goto AllowEOF
				}
			case '*':
				lx.head++
				lx.scanBlockComment()
				if lx.err != nil {
					goto Err
				}
			default:
				lx.err = UnexpectedByteError{got: c, exp: '/'}
				return exprErrCode
			}
			// continue outer for loop
		case '"', '\'':
			lval.literal.Start = lx.pos(0)
			lx.head++
			quoted := lx.readQuoted(c)
			if lx.err != nil {
				goto Err
			}
			lval.literal.End = lx.pos(0)
			lval.literal.Text = string(quoted)
			return Literal
		case '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			lval.number.Start = lx.pos(0)
			// do not forward head
			number := lx.readNumber(c)
			if lx.err != nil {
				goto Err
			}
			lval.number.End = lx.pos(0)
			lval.number.Text = number
			return Number
		case '(', ')', '=', ',', ';', '*', '<', '>', '[', ']', '{', '}', ':', '&':
			lval.identifier = &ast.Identifier{}
			lval.identifier.Node = ast.Node{
				Start: lx.pos(0),
				End:   lx.pos(1),
			}
			lval.identifier.Text = string(c)
			lx.head++
			return int(c)
		default:
			// identifier
			lval.identifier = &ast.Identifier{}
			lval.identifier.Start = lx.pos(0)
			lx.head++
			identifier := lx.readIdentifier(c)
			if lx.err != nil {
				if lx.err != io.EOF {
					return exprErrCode
				}
				// allow eof but some process still remains
			}
			s := string(identifier)

			lval.identifier.End = lx.pos(0)
			lval.identifier.Text = s

			v, ok := keywords[s]
			if ok {
				return v
			}
			return Identifier
		}
	}
AllowEOF:
	eofCode = 0
Err:
	if lx.err == io.EOF {
		return eofCode
	}
	return exprErrCode
}

func (lx *exprLexerImpl) readQuoted(terminator byte) []byte {
	buf := lx.tmpBuffer[:0]
	buf = append(buf, terminator)
	lastC := byte(0)
	for {
		for i := lx.head; i < lx.tail; i++ {
			c := lx.buffer[i]
			if c == terminator {
				if lastC != '\\' {
					buf = append(buf, lx.buffer[lx.head:i]...)
					lx.head = i + 1
					goto End
				}
			} else if c == '\n' { // newline
				lx.err = UnexpectedByteError{got: c}
				goto End
			}
			lastC = c
		}
		buf = append(buf, lx.buffer[lx.head:lx.tail]...)
		lx.head = lx.tail
		if lx.readMore(); lx.err != nil {
			goto End
		}
	}
End:
	buf = append(buf, terminator)
	lx.tmpBuffer = buf
	return buf
}

func (lx *exprLexerImpl) readIdentifier(c byte) []byte {
	buf := append(lx.tmpBuffer[:0], c)
	if !identifierStart.get(c) {
		lx.err = UnexpectedByteError{got: c}
		goto End
	}
	for {
		for i := lx.head; i < lx.tail; i++ {
			c := lx.buffer[i]
			if !identifierMiddle.get(c) {
				buf = append(buf, lx.buffer[lx.head:i]...)
				lx.head = i
				goto End
			}
		}
		buf = append(buf, lx.buffer[lx.head:lx.tail]...)
		lx.head = lx.tail
		if lx.readMore(); lx.err != nil {
			goto End
		}
	}
End:
	lx.tmpBuffer = buf
	return buf
}
