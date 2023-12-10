package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	inputFilename = "input.txt"
)

var (
	fileLines []string
)

func init() {
	file, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines = strings.Split(string(file[:]), "\n")
}

type IntSequence []int

func (sequence IntSequence) generateDiffArray() []int {
	diffArray := make([]int, 0, len(sequence)-1)

	for i := 0; i < len(sequence)-1; i++ {
		diffArray = append(diffArray, sequence[i+1]-sequence[i])
	}
	return diffArray
}

func (sequence IntSequence) allEqual(ref int) bool {
	for _, value := range sequence {
		if value != ref {
			return false
		}
	}
	return true
}

func (initialSequence IntSequence) extrapolate() int {
	subSequence := make([]IntSequence, 1)
	subSequence[0] = initialSequence

	for !subSequence[len(subSequence)-1].allEqual(0) {
		subSequence = append(subSequence, subSequence[len(subSequence)-1].generateDiffArray())
	}
	bufferNumber := 0
	for i := len(subSequence) - 2; i >= 0; i-- {
		bufferNumber = bufferNumber + subSequence[i][len(subSequence[i])-1]
	}
	return bufferNumber
}

func parseHistory(line string) (IntSequence, error) {
	numbersString := strings.Split(line, " ")
	sequence := make(IntSequence, 0, len(numbersString))

	for _, v := range numbersString {
		parsed, parsingError := strconv.Atoi(v)

		if parsingError == nil {
			sequence = append(sequence, parsed)
		} else {
			return IntSequence{}, fmt.Errorf("parsing error: %v", parsingError)
		}
	}
	return sequence, nil
}

func main() {
	// histories := make([]IntSequence, len(fileLines))
	sum := 0

	for _, line := range fileLines {
		sequence, err := parseHistory(line)

		if err == nil {
			// histories = append(histories, sequence)
			sum += sequence.extrapolate()
		} else {
			log.Panicf("error during number parsing: %v", err)
		}
	}

	fmt.Printf("result: %d\n", sum)
}
