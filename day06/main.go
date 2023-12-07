package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	inputFilename = "input.txt"
)

var (
	fileLines       []string
	mergeRacesInput bool
)

func init() {
	flag.BoolVar(&mergeRacesInput, "one-race", true, "there is simply one race with all the concatenated numbers")
	flag.Parse()
}

func init() {
	file, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines = strings.Split(string(file[:]), "\n")
}

type RaceRecord struct {
	time, distance int
}

// Equation is as follows f(x)=-x²*tx-d (a: -1, b: t, c: -d)
// t is race.time
// d is race.distance
func solveRace(race RaceRecord) (float64, float64, error) {
	time := float64(race.time)
	distance := float64(race.distance)
	delta := time*time - 4*distance // ∆ = b^2-4ac

	if delta >= 0 {
		x0, x1 := (-time-math.Sqrt(delta))/-2, (-time+math.Sqrt(delta))/-2

		if x0 < x1 {
			if x0 == math.Ceil(x0) {
				x0++
			}
			if x1 == math.Floor(x1) {
				x1--
			}
			return math.Ceil(x0), math.Floor(x1), nil
		} else {
			if x0 == math.Floor(x0) {
				x0--
			}
			if x1 == math.Ceil(x1) {
				x1++
			}
			return math.Ceil(x1), math.Floor(x0), nil
		}
	}
	return 0, 0, fmt.Errorf("race has no real solutions")
}

func parseFile() ([]RaceRecord, error) {
	reg := regexp.MustCompile(`[[:digit:]]+`)
	timeResults := reg.FindAllString(fileLines[0], -1)
	distanceResults := reg.FindAllString(fileLines[1], -1)

	if mergeRacesInput {
		timeResults = []string{
			strings.Join(timeResults, ""),
		}
		distanceResults = []string{
			strings.Join(distanceResults, ""),
		}
	}

	races := make([]RaceRecord, len(timeResults))

	for index := range timeResults {
		time, timeError := strconv.Atoi(timeResults[index])
		distance, distanceError := strconv.Atoi(distanceResults[index])

		if timeError != nil || distanceError != nil {
			return nil, fmt.Errorf("parsing error of the input file\ntime error: %v\ndistance error: %v", timeError, distanceError)
		}

		races[index] = RaceRecord{
			time:     time,
			distance: distance,
		}
	}
	return races, nil
}

func main() {
	races, parsingError := parseFile()

	if parsingError != nil {
		log.Panicf("%v\n", parsingError)
	}

	sum := 1
	for _, race := range races {
		raceRootLeft, raceRootRight, err := solveRace(race)

		if err == nil {
			ye := math.Abs(raceRootLeft - raceRootRight)
			wat := int(ye) + 1
			sum *= wat
		}
	}
	fmt.Printf("result: %d\n", sum)
}
