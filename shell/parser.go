package shell

import (
	"bufio"
	"fmt"
	"strconv"
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
			for {
				joiner := parser.accept(TokenPipe)
				if joiner == nil {
					break
				}

				node2 := parser.simple()
				if node2 == nil {
					if parser.err == nil {
						parser.err = fmt.Errorf("required command after joiner")
					}
					parser.discardTo(TokenEnd, TokenEOF)
					return nil
				}

				switch joiner.Ttype {
				case TokenPipe:
					node = &PipeNode{First: node, Second: node2}
				}
			}

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

	var node Node
	if len(words) != 0 {
		node = &SimpleNode{Words: words}
	}

	for redir := parser.redirect(node); redir != nil; redir = parser.redirect(node) {
		node = redir
	}

	return node
}

func (parser *Parser) redirect(node Node) Node {
	token := parser.accept(TokenRedirIn, TokenRedirOut, TokenRedirAppendOut)
	if token == nil {
		return nil
	}

	filename := parser.expect(TokenWord)
	if filename == nil {
		return nil
	}

	switch token.Ttype {
	case TokenRedirIn:
		var fd int
		if token.Lexeme == "" {
			fd = 0
		} else {
			fd, _ = strconv.Atoi(token.Lexeme)
		}

		return &RedirectOutNode{
			Filename: filename.Lexeme,
			Fd:       fd,
			Node:     node,
		}
	case TokenRedirAppendOut:
		fallthrough
	case TokenRedirOut:
		var fd int
		if token.Lexeme == "" {
			fd = 1
		} else {
			fd, _ = strconv.Atoi(token.Lexeme)
		}

		return &RedirectOutNode{
			Filename: filename.Lexeme,
			Fd:       fd,
			Append:   parser.last.Ttype == TokenRedirAppendOut,
			Node:     node,
		}
	default:
		panic("unreachable status")
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
