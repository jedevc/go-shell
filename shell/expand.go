package shell

import "strings"

func StripQuotes(arg string) string {
	return strings.NewReplacer("'", "", "\"", "").Replace(arg)
}
