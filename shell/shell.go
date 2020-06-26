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
		if err != nil {
			if err != io.EOF {
				log.Printf("read error: %s", err)
				return 1
			}
			break
		}

		fmt.Print(line)
	}

	return 0
}
