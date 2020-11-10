package parser

type ValueType int

const (
	// WhiteSpaceValue the next token is whitespace
	WhiteSpaceValue ValueType = iota

	InvalidValue
)

type bits [charNum / 64]int64

func (bits *bits) set(c byte) {
	bits[c>>6] |= (1 << (c & 0x3F))
}

func (bits *bits) get(c byte) bool {
	return bits[c>>6]&(1<<(c&0x3F)) != 0
}

var (
	valueTypeMap     [charNum]ValueType
	identifierStart  bits
	identifierMiddle bits
)

func init() {
	for i := 0; i < charNum; i++ {
		valueTypeMap[i] = InvalidValue
	}
	for _, c := range " \n\t\r" {
		valueTypeMap[c] = WhiteSpaceValue
	}
	var c byte
	for c = 'A'; c <= 'Z'; c++ {
		identifierStart.set(c)
		identifierMiddle.set(c)
	}
	for c = 'a'; c <= 'z'; c++ {
		identifierStart.set(c)
		identifierMiddle.set(c)
	}
	for c = '0'; c <= '9'; c++ {
		identifierMiddle.set(c)
	}
	identifierStart.set('_')
	for _, c := range "._-" {
		identifierMiddle.set(byte(c))
	}
}
