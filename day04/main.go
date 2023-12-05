package main

import (
	"flag"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Card struct {
	number, matchAmount, copies int
}

const (
	inputFilename = "input.txt"
)

var (
	useCardCopyRule bool
	fileLines       []string
	cardRegex       = regexp.MustCompile(`(?m)Card\s+([0-9]+): ([0-9\s]+) \| ([0-9\s]+)$`)
)

func init() {
	flag.BoolVar(&useCardCopyRule, "match-wins-copy", false, "Every match on your card gives you an extra copy of the n next cards (where n is the amount of matches for your card)")
	flag.Parse()
}

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

func parseCard(line string) (int, int) {
	var winningNumbers []string
	var playedNumbers []string
	match := 0

	results := cardRegex.FindAllStringSubmatch(line, -1)[0]
	cardNumber, conversionError := strconv.Atoi(results[1])
	winningNumbers = strings.Split(results[2], " ")
	playedNumbers = strings.Split(results[3], " ")

	if conversionError != nil {
		log.Panicf("error: card number (%s) couldn't be parsed:\n%v", results[1], conversionError)
	}
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
	return cardNumber, match
}

func main() {
	cards := make([]Card, 0, len(fileLines))

	for _, line := range fileLines {
		cardNumber, cardMatchAmount := parseCard(line)
		cards = append(cards, Card{
			number:      cardNumber,
			matchAmount: cardMatchAmount,
			copies:      1,
		})
	}

	pointTotal := 0
	cardTotal := 0
	for index, card := range cards {
		pointTotal += evaluateCardPoints(card.matchAmount) * card.copies
		if useCardCopyRule {
			startIndex := index + 1
			endIndex := int(math.Min(float64(startIndex+card.matchAmount), float64(len(cards))))

			for i := range cards[startIndex:endIndex] {
				cards[startIndex+i].copies += card.copies
			}
		}
		cardTotal += card.copies
	}

	if useCardCopyRule {
		println("result :", cardTotal)
	} else {
		println("result :", pointTotal)
	}
}
