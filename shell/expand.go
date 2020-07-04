package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StripQuotes(arg string) string {
	return strings.NewReplacer("'", "", "\"", "").Replace(arg)
}

func Split(arg string) []string {
	args := make([]string, 0)
	current := &strings.Builder{}

	scan := ScannerTool{}
	scan.Init(bufio.NewReader(strings.NewReader(arg)))
	scan.Advance()
	for !scan.Done() {
		switch scan.Char {
		case '\'', '"':
			quote := scan.Char
			scan.Advance()
			part := scan.ReadUntil(quote)
			fmt.Fprintf(current, "%c%s%c", quote, part, quote)
			scan.Advance()
		case ' ', '\t', '\n':
			scan.Advance()
			if current.Len() != 0 {
				args = append(args, current.String())
			}
			current.Reset()
		default:
			current.WriteRune(scan.Char)
			scan.Advance()
		}
	}
	if current.Len() != 0 {
		args = append(args, current.String())
	}

	return args
}

func Expand(ctx ExecContext, arg string) string {
	arg = ExpandTilde(ctx, arg)
	arg = ExpandVars(ctx, arg)
	arg = ExpandCommandSub(ctx, arg)
	return arg
}

func ExpandTilde(ctx ExecContext, arg string) string {
	builder := &strings.Builder{}
	scan := ScannerTool{}
	scan.Init(bufio.NewReader(strings.NewReader(arg)))
	scan.Advance()
	for !scan.Done() {
		switch scan.Char {
		case '\'', '"':
			quote := scan.Char
			scan.Advance()
			part := scan.ReadUntil(quote)
			fmt.Fprintf(builder, "%c%s%c", quote, part, quote)
			scan.Advance()
		case '~':
			scan.Advance()
			home, err := os.UserHomeDir()
			if err != nil {
				ctx.Log.Print(err)
			}
			builder.WriteString(home)
		default:
			builder.WriteRune(scan.Char)
			scan.Advance()
		}
	}
	return builder.String()
}

func ExpandCommandSub(ctx ExecContext, arg string) string {
	builder := &strings.Builder{}
	scan := ScannerTool{}
	scan.Init(bufio.NewReader(strings.NewReader(arg)))
	scan.Advance()
	for !scan.Done() {
		switch scan.Char {
		case '\'':
			quote := scan.Char
			scan.Advance()
			part := scan.ReadUntil(quote)
			fmt.Fprintf(builder, "%c%s%c", quote, part, quote)
			scan.Advance()
		case '$':
			scan.Advance()

			if scan.Char != '(' {
				builder.WriteRune('$')
				break
			}
			scan.Advance()
			command := scan.ReadUntil(')')
			scan.Advance()

			execBuilder := &strings.Builder{}
			ctx.Stdout = execBuilder
			ExecString(ctx, command)
			result := execBuilder.String()
			result = strings.TrimSpace(result)
			builder.WriteString(result)
		default:
			builder.WriteRune(scan.Char)
			scan.Advance()
		}
	}
	return builder.String()
}

func ExpandVars(ctx ExecContext, arg string) string {
	lookup := func(key string) string {
		if value, ok := ctx.Variables[key]; ok {
			return value
		} else if value, ok := os.LookupEnv(key); ok {
			return value
		} else {
			return ""
		}
	}

	builder := &strings.Builder{}
	scan := ScannerTool{}
	scan.Init(bufio.NewReader(strings.NewReader(arg)))
	scan.Advance()
	for !scan.Done() {
		switch scan.Char {
		case '\'':
			quote := scan.Char
			scan.Advance()
			part := scan.ReadUntil(quote)
			fmt.Fprintf(builder, "%c%s%c", quote, part, quote)
			scan.Advance()
		case '$':
			scan.Advance()
			if scan.Char == '(' {
				builder.WriteRune('$')
				break
			} else if scan.Char == '{' {
				scan.Advance()
				name := scan.ReadUntil('}')
				scan.Advance()
				builder.WriteString(lookup(name))
			} else {
				name := strings.Builder{}
				for !scan.Done() && scan.Char != '\'' && scan.Char != '"' && scan.Char != '$' {
					name.WriteRune(scan.Char)
					scan.Advance()
				}
				builder.WriteString(lookup(name.String()))
			}
		default:
			builder.WriteRune(scan.Char)
			scan.Advance()
		}
	}
	return builder.String()
}
