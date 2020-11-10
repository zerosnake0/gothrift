package parser

import (
	"io"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

const bufferSize = 512

type exprLexerImplEmbedded struct {
	capture bool
	offset  int

	// the current index position
	head int

	// the current line number (begin from zero)
	lineNo int
	// the offset of the line beginning
	lineBegOffset int

	err error

	comments []ast.Comment

	Document *ast.Document
}

type exprLexerImpl struct {
	reader io.Reader
	buffer []byte
	fixbuf []byte

	tmpBuffer []byte

	tail int

	exprLexerImplEmbedded
}
