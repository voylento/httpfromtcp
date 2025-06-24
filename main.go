package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	currentLineContents := ""
	go func() {
		defer close(ch)
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					ch <- currentLineContents
					currentLineContents = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("Error: %s\n", err.Error())
				break
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				ch <- fmt.Sprintf("%s%s\n", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents = parts[len(parts)-1]
		}
	}()

	return ch
}


func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("failed to open file %s: %s\n", inputFilePath, err)
	}
	defer file.Close()

	fmt.Printf("Reading contents from %s\n", inputFilePath)
	fmt.Println("========================================")

	ch := getLinesChannel(file)
	for {
		str, ok := <- ch
		if !ok {
			break
		}
		fmt.Printf("%s\n", str)
	}
}
