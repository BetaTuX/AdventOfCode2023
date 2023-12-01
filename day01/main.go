package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	inputFile = "./calibration_input.txt"
)

var (
	numbersAsWord = []string{
		"zero",
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
	}
	enableNumbersAsLetters bool
)

func init() {
	flag.BoolVar(&enableNumbersAsLetters, "numbers-can-be-words", false, "Enable this flag if you want to parse coordinates using also numbers written in letters (one, two, three...)")
	flag.Parse()
}

func identifyLinePrefix(line string) (int, error) {
	firstChar := line[0]

	if '0' <= firstChar && firstChar <= '9' {
		return int(firstChar - '0'), nil
	}
	if enableNumbersAsLetters {
		for index, numberWord := range numbersAsWord {
			if strings.HasPrefix(line, numberWord) {
				return index, nil
			}
		}
	}
	return -1, fmt.Errorf("no number could be identified in the following string: %s", line)
}

func main() {
	file, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines := strings.Split(string(file[:]), "\n")
	coordinatesArray := make([]int, 0, len(fileLines))
	for _, line := range fileLines {
		firstDigit, lastDigit := -1, -1

		for index := range line {
			number, _ := identifyLinePrefix(line[index:])

			if number < 0 {
				continue
			}

			if firstDigit < 0 {
				firstDigit = number
			}
			lastDigit = number
		}
		coordinatesArray = append(coordinatesArray, firstDigit*10+lastDigit)
	}
	total := 0
	for _, v := range coordinatesArray {
		total += v
	}
	log.Println("processed coordinates:", total)
}
