package main

import (
	"log"
	"math"
	"os"
	"regexp"
	"strings"
)

const (
	inputFilename = "input.txt"
)

var (
	fileLines []string
	cardRegex = regexp.MustCompile(`(?m)Card\s+[0-9]+: ([0-9\s]+) \| ([0-9\s]+)$`)
)

func evaluateCardPoints(nbOfMatches int) int {
	if nbOfMatches <= 0 {
		return 0
	}
	return int(math.Pow(float64(2), float64(nbOfMatches-1)))
}

func init() {
	file, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines = strings.Split(string(file[:]), "\n")
}

func parseCard(line string) int {
	var winningNumbers []string
	var playedNumbers []string
	match := 0

	results := cardRegex.FindAllStringSubmatch(line, -1)[0]
	winningNumbers = strings.Split(results[1], " ")
	playedNumbers = strings.Split(results[2], " ")

	for _, winningRef := range winningNumbers {
		if winningRef == "" {
			continue
		}
		for _, number := range playedNumbers {
			if strings.Trim(number, " ") == strings.Trim(winningRef, " ") {
				match++
			}
		}
	}
	return evaluateCardPoints(match)
}

func main() {
	cardResults := make([]int, 0, len(fileLines))

	for _, line := range fileLines {
		cardResults = append(cardResults, parseCard(line))
	}

	sum := 0
	for _, result := range cardResults {
		sum += result
	}

	println("result :", sum)
}
