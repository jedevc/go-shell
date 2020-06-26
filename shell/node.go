package shell

import (
	"os"
	"os/exec"
)

type Node interface {
	Exec(ctx ExecContext) int
}

type SimpleNode struct {
	Words []string
}

func (node *SimpleNode) Exec(ctx ExecContext) int {
	// Execute a builtin
	if builtin, ok := Builtins[node.Words[0]]; ok {
		return builtin(ctx, node.Words[0], node.Words[1:]...)
	}

	// Execute an external command
	cmd := exec.Command(node.Words[0], node.Words[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			// Some internal error occurred
			ctx.Log.Print(err)
			return 1
		}
	}

	return cmd.ProcessState.ExitCode()
}
