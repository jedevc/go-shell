package shell

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type ScannerTool struct {
	reader *bufio.Reader

	Char rune
	Last rune

	eof bool
	err error
}

func (scanner *ScannerTool) Init(reader *bufio.Reader) {
	scanner.reader = reader
}

func (scanner *ScannerTool) Done() bool {
	return scanner.eof
}

func (scanner *ScannerTool) Error() error {
	err := scanner.err
	scanner.err = nil
	return err
}

func (scanner *ScannerTool) Advance() rune {
	ch, _, err := scanner.reader.ReadRune()
	if err == io.EOF {
		scanner.eof = true
	} else if err != nil {
		scanner.eof = true
		if scanner.err == nil {
			scanner.err = err
		}
	}

	scanner.Last = scanner.Char
	scanner.Char = ch
	return ch
}

func (scanner *ScannerTool) ReadUntil(terminator rune) string {
	builder := strings.Builder{}
	for !scanner.eof && scanner.Char != terminator {
		builder.WriteRune(scanner.Char)
		scanner.Advance()
	}
	if scanner.eof && scanner.err == nil {
		scanner.err = fmt.Errorf("unexpected EOF while looking for %c", terminator)
	}
	return builder.String()
}

func (scanner *ScannerTool) Peek() rune {
	ch, _, err := scanner.reader.ReadRune()
	if err != nil {
		scanner.eof = true
		if scanner.err == nil {
			scanner.err = err
		}
		return ch
	}
	_ = scanner.reader.UnreadRune()
	return ch
}
