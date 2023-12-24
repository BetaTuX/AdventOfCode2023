package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	inputFilename = "input.txt"
)

type Code uint8

func mapStringToCode(s string, initialValue Code) Code {
	for charIndex := range s {
		initialValue += Code(s[charIndex])
		initialValue *= 17
	}
	return initialValue
}

func main() {
	file, openError := os.ReadFile(inputFilename)

	if openError != nil {
		log.Fatalln("couldn't open file input.txt")
	}
	codes := strings.Split(string(file[:]), ",")
	result := 0

	for _, code := range codes {
		result += int(mapStringToCode(code, 0))
	}
	fmt.Printf("result hash: %d\n", result)
}
