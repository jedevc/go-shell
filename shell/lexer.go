package shell

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	TokenWord = iota
	TokenEnd  = iota
	TokenEOF  = iota
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
	default:
		return lexer.readWord()
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
		default:
			builder.WriteRune(lexer.Char)
			lexer.Advance()
		}
	}

	return &Token{TokenWord, builder.String()}
}

func isBreak(ch rune) bool {
	return strings.ContainsRune(" \n\t;", ch)
}
