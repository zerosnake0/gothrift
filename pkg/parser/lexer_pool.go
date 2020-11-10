package parser

import "sync"

var (
	defaultLexerPool = newLexerPool()
)

type lexerPool struct {
	pool sync.Pool
}

func newLexerPool() *lexerPool {
	return &lexerPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &exprLexerImpl{
					fixbuf:    make([]byte, 64),
					tmpBuffer: make([]byte, 64),
				}
			},
		},
	}
}

func (p *lexerPool) borrowLexer() *exprLexerImpl {
	return p.pool.Get().(*exprLexerImpl)
}

func (p *lexerPool) returnLexer(lx *exprLexerImpl) {
	lx.reset()
	p.pool.Put(lx)
}
