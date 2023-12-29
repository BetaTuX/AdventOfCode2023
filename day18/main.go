package main

import (
	"fmt"
	"image"
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

type DigInstruction struct {
	Direction image.Point
	Length    int
	Color     string
}

func ParseInputLine(line string) (DigInstruction, error) {
	re := regexp.MustCompile(`([URDL]) ([0-9]+) \(#([0-9a-z]{6})\)`)
	parsed := re.FindStringSubmatch(line)
	directionMap := map[string]image.Point{
		"U": {0, -1},
		"R": {1, 0},
		"D": {0, 1},
		"L": {-1, 0},
	}
	length, parsingErr := strconv.Atoi(parsed[2])

	if parsingErr != nil {
		return DigInstruction{}, fmt.Errorf("line parsing error: %v", parsingErr)
	}
	return DigInstruction{
		Direction: directionMap[parsed[1]],
		Length:    length,
		Color:     parsed[3],
	}, nil
}

func PartOne(lines []string) int {
	result := 0
	a := image.Point{0, 0}
	posSum := 0
	negSum := 0

	for index := 0; index < len(lines); index++ {
		line := lines[index]
		point, _ := ParseInputLine(line)
		b := a.Add(point.Direction.Mul(point.Length))

		posSum += a.X*b.Y + point.Length
		negSum += a.Y * b.X
		a = b
	}
	result = int(math.Abs(float64(posSum - negSum)))
	return result/2 + 1
}

func main() {
	file, openError := os.ReadFile(inputFilename)
	if openError != nil {
		log.Panicf("couldn't open input file '%s'\n%v\n", inputFilename, openError)
	}
	lines := strings.Split(string(file), "\n")
	partOneResult := PartOne(lines)

	fmt.Printf("result: %d\n", partOneResult)
}
