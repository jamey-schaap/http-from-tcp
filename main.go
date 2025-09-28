package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)

		var currentLine string
		for {
			bytes := make([]byte, 8)
			_, err := f.Read(bytes)

			currentLine += string(bytes)
			parts := strings.Split(currentLine, "\n")
			if len(parts) > 1 {
				for i := 0; i < len(parts)-1; i++ {
					lines <- parts[i]
				}
				currentLine = parts[len(parts)-1]
			}

			if err == io.EOF {
				lines <- currentLine
				return
			}
		}
	}()

	return lines
}
