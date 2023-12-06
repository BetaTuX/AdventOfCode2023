package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	inputFilename = "input.txt"
)

var (
	fileLines []string
)

type Range struct {
	// Destination range start
	destination int
	// Source range start
	source int
	// Ranges length
	length int
}

type PuzzleMap struct {
	sourceId      string
	destinationId string
	ranges        []Range
}

type MapChain struct {
	startId string
	maps    []PuzzleMap
}

func (chain *MapChain) setBlockPosition(blockId string, position int) error {
	var blockIndex int

	if blockIndex = slices.IndexFunc(chain.maps, func(m PuzzleMap) bool {
		return m.sourceId == chain.startId
	}); blockIndex < 0 {
		return fmt.Errorf("couldn't sort chain because it is missing its starting block (id: %v)", chain.startId)
	}
	chain.maps[0], chain.maps[blockIndex] = chain.maps[blockIndex], chain.maps[0]
	return nil
}

func (chain *MapChain) Sort() error {
	chain.setBlockPosition(chain.startId, 0)
	for index := range chain.maps {
		if err := chain.setBlockPosition(chain.maps[index].destinationId, index+1); err != nil {
			return err
		}
	}
	return nil
}

func (chain MapChain) Evaluate(id int) int {
	for _, block := range chain.maps {
		id = block.getNumber(id)
	}
	return id
}

func (r Range) isInRange(number int) bool {
	return r.source <= number && number <= r.source+r.length
}

func (r Range) getNumber(number int) int {
	if r.isInRange(number) {
		return number + r.destination - r.source
	}
	return number
}

func rangeFromString(input string) (Range, error) {
	inputs := strings.Split(input, " ")
	errors := make([]error, 3)
	var destinationRangeStart, sourceRangeStart, rangeLength int

	if len(inputs) != 3 {
		return Range{}, fmt.Errorf("range parsing error: found more than 3 numbers")
	}

	destinationRangeStart, errors[0] = strconv.Atoi(inputs[0])
	sourceRangeStart, errors[1] = strconv.Atoi(inputs[1])
	rangeLength, errors[2] = strconv.Atoi(inputs[2])

	for _, err := range errors {
		if err != nil {
			return Range{}, fmt.Errorf("range parsing error: %v", err)
		}
	}

	return Range{
		destination: destinationRangeStart,
		source:      sourceRangeStart,
		length:      rangeLength,
	}, nil
}

func puzzleMapFromLines(input *[]string) (PuzzleMap, error) {
	var reg = regexp.MustCompile(`(?m)([[:alpha:]]+)-to-([[:alpha:]]+).*$`)

	if len((*input)) <= 0 {
		return PuzzleMap{}, fmt.Errorf("cannot parse map: input is empty")
	}

	for !reg.MatchString((*input)[0]) {
		*input = (*input)[1:]
	}

	regResults := reg.FindStringSubmatch((*input)[0])
	*input = (*input)[1:]
	if len(regResults) != 3 {
		return PuzzleMap{}, fmt.Errorf("map parsing error: couldn't identify map resource ids")
	}

	ranges := make([]Range, 0)

	for ; len(*input) > 0 && (*input)[0] != ""; *input = (*input)[1:] {
		if r, err := rangeFromString((*input)[0]); err == nil {
			ranges = append(ranges, r)
		} else {
			return PuzzleMap{}, err
		}
	}

	return PuzzleMap{
		sourceId:      regResults[1],
		destinationId: regResults[2],
		ranges:        ranges,
	}, nil
}

func (m PuzzleMap) getNumber(number int) int {
	for _, r := range m.ranges {
		if r.isInRange(number) {
			return r.getNumber(number)
		}
	}
	return number
}

func init() {
	file, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines = strings.Split(string(file[:]), "\n")
}

func parseSeeds(line string) ([]int, error) {
	reg := regexp.MustCompile(`seeds: ([[:digit:] ]+)$`)
	regResults := reg.FindStringSubmatch(line)
	var seeds []int

	for _, parsedNumber := range strings.Split(regResults[1], " ") {
		if seedNumber, parsingError := strconv.Atoi(parsedNumber); parsingError != nil {
			return seeds, fmt.Errorf("couldn't parse seed number\n%v", parsingError)
		} else {
			seeds = append(seeds, seedNumber)
		}
	}

	return seeds, nil
}

func main() {
	seeds, seedParsingErr := parseSeeds(fileLines[0])
	mapPuzzles := make([]PuzzleMap, 0)

	if seedParsingErr != nil {
		log.Panic(seedParsingErr)
	}

	for {
		if p, err := puzzleMapFromLines(&fileLines); err == nil {
			mapPuzzles = append(mapPuzzles, p)
		} else {
			break
		}
	}

	chain := MapChain{
		startId: "seed",
		maps:    mapPuzzles,
	}

	if err := chain.Sort(); err != nil {
		log.Fatalln(err)
	}

	possibleLocations := make([]int, 0, len(mapPuzzles))
	for _, seed := range seeds {
		possibleLocations = append(possibleLocations, chain.Evaluate(seed))
	}
	closestLocation := possibleLocations[0]
	for _, opportunity := range possibleLocations[1:] {
		if opportunity < closestLocation {
			closestLocation = opportunity
		}
	}

	fmt.Printf("result: %d\n", closestLocation)
}
