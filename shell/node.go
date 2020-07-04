package shell

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
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

type RedirectInNode struct {
	Node
	Filename string
	Fd       int
}

type RedirectOutNode struct {
	Node
	Filename string
	Append   bool
	Fd       int
}

func (node *RedirectInNode) Exec(ctx ExecContext) int {
	file, err := os.Open(node.Filename)
	if err != nil {
		ctx.Log.Print(err)
		return 1
	}
	defer file.Close()

	switch node.Fd {
	case 0:
		ctx.Stdin = file
	default:
		ctx.Log.Printf("cannot redirect from %d", node.Fd)
		return 1
	}

	return node.Node.Exec(ctx)
}

func (node *RedirectOutNode) Exec(ctx ExecContext) int {
	flags := os.O_CREATE | os.O_WRONLY
	if node.Append {
		flags |= os.O_APPEND
	}

	file, err := os.OpenFile(node.Filename, flags, 0644)
	if err != nil {
		ctx.Log.Print(err)
		return 1
	}
	defer file.Close()

	switch node.Fd {
	case 1:
		ctx.Stdout = file
	case 2:
		ctx.Stderr = file
	default:
		ctx.Log.Printf("cannot redirect to %d", node.Fd)
		return 1
	}

	return node.Node.Exec(ctx)
}

type PipeNode struct {
	First  Node
	Second Node
}

func (node *PipeNode) Exec(ctx ExecContext) int {
	reader, writer := io.Pipe()

	code := 0

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		newCtx := ctx
		newCtx.Stdout = writer
		node.First.Exec(newCtx)
		writer.Close()

		wg.Done()
	}()
	go func() {
		newCtx := ctx
		newCtx.Stdin = reader
		code = node.Second.Exec(newCtx)

		_, err := io.Copy(ioutil.Discard, reader)
		if err != nil {
			ctx.Log.Print(err)
		}

		wg.Done()
	}()
	wg.Wait()

	return code
}
