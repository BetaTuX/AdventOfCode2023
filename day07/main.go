package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	inputFilename = "input.txt"
	LEGAL_CARDS   = "AKQJT98765432"
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

type Hand struct {
	Cards string
	Type  int
	Bid   int
}

func identifyHandType(hand string) (int, error) {
	cards := make(map[rune]int)

	for _, b := range hand {
		if !strings.ContainsRune(LEGAL_CARDS, b) {
			return -1, fmt.Errorf("hand contains illegal cards")
		}
		prevValue := cards[b]
		cards[b] = prevValue + 1
	}

	maxSameCard := 0
	differentCardKind := 0
	pairAmount := 0
	for _, amount := range cards {
		differentCardKind++
		if amount == 2 {
			pairAmount++
		}
		if amount > maxSameCard {
			maxSameCard = amount
		}
	}

	switch {
	case maxSameCard == 5:
		return 0, nil
	case maxSameCard == 4:
		return 1, nil
	case differentCardKind == 2 && maxSameCard == 3 && pairAmount == 1:
		return 2, nil
	case maxSameCard == 3:
		return 3, nil
	case pairAmount == 2:
		return 4, nil
	case pairAmount == 1:
		return 5, nil
	default:
		return 6, nil
	}
}

func parseHand(source string) (Hand, error) {
	parts := strings.Split(source, " ")

	if len(parts) != 2 || len(parts[0]) != 5 {
		return Hand{}, fmt.Errorf("parsing error, source must have 2 parts separated by a ' ' and left part should be a string of 5 characters representing the hand of the player")
	}

	bid, convError := strconv.Atoi(parts[1])
	cardType, identifyError := identifyHandType(parts[0])

	if convError != nil {
		return Hand{}, fmt.Errorf("parsing error on bid part, %v", convError)
	}
	if identifyError != nil {
		return Hand{}, fmt.Errorf("parsing error on hand part, %v", identifyError)
	}

	return Hand{
		Cards: parts[0],
		Type:  cardType,
		Bid:   bid,
	}, nil
}

type ByHandPower []Hand

func (ranges ByHandPower) Len() int      { return len(ranges) }
func (ranges ByHandPower) Swap(i, j int) { ranges[i], ranges[j] = ranges[j], ranges[i] }
func (ranges ByHandPower) Less(i, j int) bool {
	if ranges[i].Type > ranges[j].Type {
		return true
	} else if ranges[i].Type < ranges[j].Type {
		return false
	}
	for index := range ranges[i].Cards {
		left := strings.IndexByte(LEGAL_CARDS, ranges[i].Cards[index])
		right := strings.IndexByte(LEGAL_CARDS, ranges[j].Cards[index])

		if left == right {
			continue
		}
		return left > right
	}
	return true
}

func main() {
	hands := make([]Hand, 0, len(fileLines))

	for _, line := range fileLines {
		if hand, err := parseHand(line); err == nil {
			hands = append(hands, hand)
		} else {
			log.Panicf("Parsing error:\n%v", err)
		}
	}

	sort.Sort(ByHandPower(hands))

	sum := 0
	for index, hand := range hands {
		sum += hand.Bid * (index + 1)
	}
	fmt.Printf("result: %d\n", sum)
}
