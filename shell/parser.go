package shell

import (
	"bufio"
	"io"
	"strings"
)

type Parser struct {
	reader *bufio.Reader
	eof    bool
	err    error
}

func (parser *Parser) Init(reader *bufio.Reader) {
	parser.reader = reader
}

func (parser *Parser) Done() bool {
	return parser.eof
}

func (parser *Parser) Error() error {
	return parser.err
}

func (parser *Parser) Parse() Node {
	if parser.eof {
		return nil
	}
	return parser.line()
}

func (parser *Parser) line() Node {
	line, err := parser.reader.ReadString('\n')
	if err == io.EOF {
		parser.eof = true
	} else if err != nil {
		parser.err = err
		return nil
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil
	} else {
		words := strings.Split(line, " ")
		return &SimpleNode{Words: words}
	}
}
