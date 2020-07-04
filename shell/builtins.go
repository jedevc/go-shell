package shell

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Command func(ctx ExecContext, name string, args ...string) int

var Builtins = map[string]Command{
	"exit":   BuiltinExit,
	"cd":     BuiltinChangeDirectory,
	"source": BuiltinSource,
	"export": BuiltinExport,
}

func BuiltinExit(ctx ExecContext, name string, args ...string) int {
	switch len(args) {
	case 0:
		os.Exit(0)
		return 0
	case 1:
		n, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(ctx.Stderr, "%s: numeric argument required\n", name)
			return 1
		}
		os.Exit(n)
		return 0
	default:
		fmt.Fprintf(ctx.Stderr, "%s: too many arguments\n", name)
		return 1
	}
}

func BuiltinChangeDirectory(ctx ExecContext, name string, args ...string) int {
	var target string
	switch len(args) {
	case 0:
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(ctx.Stderr, "%s: %s\n", name, err)
			return 1
		}
		target = home
	case 1:
		target = args[0]
	default:
		fmt.Fprintf(ctx.Stderr, "%s: too many arguments\n", name)
		return 1
	}

	err := os.Chdir(target)
	if err != nil {
		fmt.Fprintf(ctx.Stderr, "%s: %s\n", name, err)
		return 1
	}

	return 0
}

func BuiltinSource(ctx ExecContext, name string, args ...string) int {
	switch len(args) {
	case 0:
		fmt.Fprintf(ctx.Stderr, "%s: filename argument required\n", name)
		return 1
	case 1:
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(ctx.Stderr, "%s: %s\n", name, err)
		}

		return Exec(ctx, file, false)
	default:
		fmt.Fprintf(ctx.Stderr, "%s: too many arguments\n", name)
		return 1
	}
}

func BuiltinExport(ctx ExecContext, name string, args ...string) int {
	switch len(args) {
	case 0:
		for _, env := range os.Environ() {
			fmt.Fprintln(ctx.Stdout, env)
		}
	default:
		for _, env := range args {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			} else {
				if value, ok := ctx.Variables[parts[0]]; ok {
					os.Setenv(parts[0], value)
				}
			}
			delete(ctx.Variables, parts[0])
		}
	}

	return 0
}
