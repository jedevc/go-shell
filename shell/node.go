package shell

import (
	"os/exec"
)

type Node interface {
	Exec(ctx ExecContext) int
}

type SimpleNode struct {
	Words []string
}

func (node *SimpleNode) Exec(ctx ExecContext) int {
	// Strip quotes from words
	args := make([]string, 0)
	for _, word := range node.Words {
		arg := StripQuotes(word)
		args = append(args, arg)
	}

	// Execute a builtin
	if builtin, ok := Builtins[args[0]]; ok {
		return builtin(ctx, args[0], args[1:]...)
	}

	// Execute an external command
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = ctx.Stdin
	cmd.Stdout = ctx.Stdout
	cmd.Stderr = ctx.Stderr
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

type GroupNode struct {
	Children []Node
}

func (node *GroupNode) Exec(ctx ExecContext) int {
	// Execute all children, returning the last's exit code
	code := 0
	for _, child := range node.Children {
		code = child.Exec(ctx)
	}
	return code
}
