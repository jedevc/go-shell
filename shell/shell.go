package shell

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

type ExecContext struct {
	Variables map[string]string
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
	Log       *log.Logger
}

func ExecString(ctx ExecContext, text string) int {
	return Exec(ctx, strings.NewReader(text), false)
}

func Exec(ctx ExecContext, source io.Reader, interactive bool) int {
	parser := &Parser{}
	parser.Init(bufio.NewReader(source))

	code := 0
	for !parser.Done() {
		if interactive {
			fmt.Print("> ")
		}

		node := parser.Parse()
		if err := parser.Error(); err != nil {
			ctx.Log.Printf("parse error: %s", err)
			continue
		}
		if node == nil {
			continue
		}

		code = node.Exec(ctx)
	}

	return code
}
