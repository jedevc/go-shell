package shell

import (
	"bufio"
	"fmt"
)

type Parser struct {
	lexer   *Lexer
	last    *Token
	current *Token

	eof bool
	err error
}

func (parser *Parser) Init(reader *bufio.Reader) {
	parser.lexer = &Lexer{}
	parser.lexer.Init(reader)
}

func (parser *Parser) Done() bool {
	return parser.lexer.Done()
}

func (parser *Parser) Error() error {
	// Find error (if there are any)
	err := parser.lexer.Error()
	if err == nil {
		err = parser.err
	}

	// Handle error
	parser.err = nil
	return err
}

func (parser *Parser) Parse() Node {
	if parser.eof {
		return nil
	}
	return parser.line()
}

func (parser *Parser) line() Node {
	nodes := make([]Node, 0)

	for {
		node := parser.simple()
		if node != nil {
			nodes = append(nodes, node)
		}

		token := parser.expect(TokenEnd, TokenEOF)
		if token == nil {
			parser.discardTo(TokenEnd, TokenEOF)
			break
		} else if token.Ttype == TokenEOF || token.Lexeme == "\n" {
			break
		}
	}

	if len(nodes) == 0 {
		return nil
	} else {
		return &GroupNode{Children: nodes}
	}
}

func (parser *Parser) simple() Node {
	words := make([]string, 0)
	for token := parser.accept(TokenWord); token != nil; token = parser.accept(TokenWord) {
		words = append(words, token.Lexeme)
	}

	if len(words) == 0 {
		return nil
	} else {
		return &SimpleNode{Words: words}
	}
}

func (parser *Parser) consume() {
	parser.last = parser.current
	parser.current = nil
}

func (parser *Parser) fill() {
	for parser.current == nil {
		parser.current = parser.lexer.Next()
	}
}

func (parser *Parser) accept(ttypes ...uint) *Token {
	if parser.match(ttypes...) {
		parser.consume()
		return parser.last
	}
	return nil
}

func (parser *Parser) expect(ttypes ...uint) *Token {
	if token := parser.accept(ttypes...); token != nil {
		return token
	}
	if parser.err == nil {
		parser.err = fmt.Errorf("expected token, but didn't get it")
	}
	return nil
}

func (parser *Parser) discardTo(ttypes ...uint) {
	for parser.accept(ttypes...) == nil {
		parser.consume()
	}
}

func (parser *Parser) match(ttypes ...uint) bool {
	parser.fill()
	if parser.current == nil {
		return false
	}

	for _, ttype := range ttypes {
		if parser.current.Ttype == ttype {
			return true
		}
	}

	return false
}
