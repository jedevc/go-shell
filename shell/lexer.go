package shell

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	TokenWord           = iota
	TokenEnd            = iota
	TokenEOF            = iota
	TokenRedirIn        = iota
	TokenRedirOut       = iota
	TokenRedirAppendOut = iota
	TokenPipe           = iota
)

type Token struct {
	Ttype  uint
	Lexeme string
}

type Lexer struct {
	ScannerTool

	advanceLater bool
}

func (lexer *Lexer) Init(reader *bufio.Reader) {
	lexer.ScannerTool.Init(reader)
	lexer.advanceLater = true
}

func (lexer *Lexer) Next() *Token {
	if lexer.advanceLater {
		lexer.Advance()
		lexer.advanceLater = false
	}
	if lexer.eof {
		return &Token{TokenEOF, ""}
	}

	switch lexer.Char {
	case ' ', '\t':
		lexer.Advance()
		return nil
	case ';', '\n':
		lexer.advanceLater = true
		return &Token{TokenEnd, string(lexer.Char)}
	case '|':
		lexer.Advance()
		return &Token{TokenPipe, ""}
	default:
		if redir := lexer.readRedir(); redir != nil {
			return redir
		} else {
			return lexer.readWord()
		}
	}
}

func (lexer *Lexer) readWord() *Token {
	builder := &strings.Builder{}

	for !lexer.eof && !isBreak(lexer.Char) {
		switch lexer.Char {
		case '\'', '"':
			quote := lexer.Char
			lexer.Advance()
			part := lexer.ReadUntil(quote)
			lexer.Advance()
			fmt.Fprintf(builder, "%c%s%c", quote, part, quote)
		case '$':
			builder.WriteRune(lexer.Char)
			switch lexer.Advance() {
			case '{':
				lexer.Advance()
				part := lexer.ReadUntil('}')
				fmt.Fprintf(builder, "{%s}", part)
				lexer.Advance()
			}
		default:
			builder.WriteRune(lexer.Char)
			lexer.Advance()
		}
	}

	return &Token{TokenWord, builder.String()}
}

func (lexer *Lexer) readRedir() *Token {
	switch {
	case lexer.Char == '>':
		lexer.Advance()
		if lexer.Char == '>' {
			lexer.Advance()
			return &Token{TokenRedirAppendOut, ""}
		} else {
			return &Token{TokenRedirOut, ""}
		}
	case lexer.Char == '<':
		lexer.Advance()
		return &Token{TokenRedirIn, ""}
	case '0' <= lexer.Char && lexer.Char <= '9':
		n := lexer.Char

		peeked := lexer.Peek()
		if peeked == '>' || peeked == '<' {
			lexer.Advance()
			token := lexer.readRedir()
			token.Lexeme = string(n)
			return token
		}
		return nil
	default:
		return nil
	}
}

func isBreak(ch rune) bool {
	return strings.ContainsRune(" \n\t;<>|", ch)
}
