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

func Expand(ctx ExecContext, arg string) string {
	arg = ExpandVars(ctx, arg)
	return arg
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
