package main

import (
	"os"

	"github.com/jedevc/go-shell/shell"
)

func main() {
	os.Exit(sh(os.Args))
}

func sh(argv []string) int {
	if len(argv) == 0 {
		return 0
	}

	return shell.Exec(os.Stdin)
}
