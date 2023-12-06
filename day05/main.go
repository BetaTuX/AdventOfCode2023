package main

import (
	"flag"
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
	fileLines     []string
	useSeedRanges bool
)

func init() {
	flag.BoolVar(&useSeedRanges, "seed-ranges", false, "seeds are ranges instead of simple seeds")
	flag.Parse()
}

type Range struct {
	// Destination range start
	start int
	// Ranges length
	length int
}

type ByStartIndex []Range

func (ranges ByStartIndex) Len() int           { return len(ranges) }
func (ranges ByStartIndex) Swap(i, j int)      { ranges[i], ranges[j] = ranges[j], ranges[i] }
func (ranges ByStartIndex) Less(i, j int) bool { return ranges[i].start < ranges[j].start }

type Mapper struct {
	source, destination Range
}

type PuzzleMap struct {
	sourceId      string
	destinationId string
	mappers       []Mapper
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
		newId := block.GetNumber(id)
		id = newId
	}
	return id
}

func (chain MapChain) ReverseEvaluate(id int) int {
	for i := len(chain.maps) - 1; i >= 0; i-- {
		block := chain.maps[i]
		newId := block.GetRoot(id)
		id = newId
	}
	return id
}

func (r Range) IsInRange(number int) bool {
	return r.start <= number && number < r.start+r.length
}

func (r Mapper) GetNumber(number int) int {
	if r.source.IsInRange(number) {
		return number + r.destination.start - r.source.start
	}
	return number
}

func (r Mapper) GetRoot(number int) int {
	if r.destination.IsInRange(number) {
		return number - (r.destination.start - r.source.start)
	}
	return number
}

func mapperFromString(input string) (Mapper, error) {
	inputs := strings.Split(input, " ")
	errors := make([]error, 3)
	var destinationRangeStart, sourceRangeStart, rangeLength int

	if len(inputs) != 3 {
		return Mapper{}, fmt.Errorf("range parsing error: found more than 3 numbers")
	}

	destinationRangeStart, errors[0] = strconv.Atoi(inputs[0])
	sourceRangeStart, errors[1] = strconv.Atoi(inputs[1])
	rangeLength, errors[2] = strconv.Atoi(inputs[2])

	for _, err := range errors {
		if err != nil {
			return Mapper{}, fmt.Errorf("range parsing error: %v", err)
		}
	}

	return Mapper{
		destination: Range{start: destinationRangeStart, length: rangeLength},
		source:      Range{start: sourceRangeStart, length: rangeLength},
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

	mappers := make([]Mapper, 0)

	for ; len(*input) > 0 && (*input)[0] != ""; *input = (*input)[1:] {
		if m, err := mapperFromString((*input)[0]); err == nil {
			mappers = append(mappers, m)
		} else {
			return PuzzleMap{}, err
		}
	}

	return PuzzleMap{
		sourceId:      regResults[1],
		destinationId: regResults[2],
		mappers:       mappers,
	}, nil
}

func (m PuzzleMap) GetNumber(number int) int {
	result := -1

	for _, r := range m.mappers {
		if mapped := r.GetNumber(number); (r.source.IsInRange(number) && result < 0) || mapped < result {
			result = mapped
		}
	}
	if result >= 0 {
		return result
	} else {
		return number
	}
}

func (m PuzzleMap) GetRoot(number int) int {
	for _, r := range m.mappers {
		if mapped := r.GetRoot(number); r.destination.IsInRange(number) {
			return mapped
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

func mapSeedsToRangeList(seeds []int) []Range {
	ranges := make([]Range, len(seeds)/2)

	for i := 0; i < len(seeds); i += 2 {
		ranges[i/2] = Range{
			start:  seeds[i],
			length: seeds[i+1],
		}
	}
	return ranges
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

	var closestLocation = -1
	var rangeList []Range
	if useSeedRanges {
		rangeList = mapSeedsToRangeList(seeds)
	} else {
		rangeList = make([]Range, 0, len(seeds))
		for _, seed := range seeds {
			rangeList = append(rangeList, Range{
				start:  seed,
				length: 1,
			})
		}
	}

	for location := 1; closestLocation < 0; location++ {
		root := chain.ReverseEvaluate(location)

		for _, r := range rangeList {
			if r.IsInRange(root) && (location < closestLocation || closestLocation < 0) {
				closestLocation = location
				break
			}
		}
	}

	fmt.Printf("result: %d\n", closestLocation)
}
