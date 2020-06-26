package shell

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

func Exec(source io.Reader) int {
	reader := bufio.NewReader(source)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			fmt.Print(line)
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("read error: %s", err)
				return 1
			}
			break
		}
	}

	return 0
}
