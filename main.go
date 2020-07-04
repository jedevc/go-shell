package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jedevc/go-shell/shell"
)

const NAME = "go-shell"

func main() {
	commandString := flag.String("c", "", "command string")
	flag.Parse()
	rest := flag.Args()

	logger := log.New(os.Stderr, fmt.Sprintf("%s: ", NAME), 0)

	ctx := shell.ExecContext{
		Variables: map[string]string{},
		Stdin:     os.Stdin,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		Log:       logger,
	}

	var code int
	switch {
	case len(*commandString) > 0:
		code = shell.ExecString(ctx, *commandString)
	case len(rest) == 0:
		code = shell.Exec(ctx, os.Stdin, true)
	default:
		file, err := os.Open(rest[0])
		if err != nil {
			logger.Fatal(err)
		}
		defer file.Close()

		code = shell.Exec(ctx, file, false)
	}

	os.Exit(code)
}
