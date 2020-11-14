.PHONY: parser
parser:
	go generate ./pkg/parser

.PHONY: install
install:
	go install ./cmd/thriftfmt
