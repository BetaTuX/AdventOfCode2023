package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	inputFilename = "input.txt"
)

var (
	colorIsLength bool
)

type DigInstruction struct {
	Direction image.Point
	Length    int
}

func init() {
	flag.BoolVar(&colorIsLength, "color-as-length", false, "Sets the color in the input as both length and direction code")
}

func ParseInputLine(line string) (DigInstruction, error) {
	re := regexp.MustCompile(`([URDL]) ([0-9]+) \(#([0-9a-z]{6})\)`)
	parsed := re.FindStringSubmatch(line)
	directionMap := map[string]image.Point{
		"R": {1, 0},
		"D": {0, 1},
		"L": {-1, 0},
		"U": {0, -1},
	}
	// Order here is important for part 2
	directionArr := [4]image.Point{
		{1, 0},
		{0, 1},
		{-1, 0},
		{0, -1},
	}

	var direction image.Point
	var length int
	var parsingErr error

	if colorIsLength {
		length64, parsing64Err := strconv.ParseInt(parsed[3][:5], 16, strconv.IntSize)
		length = int(length64)
		parsingErr = parsing64Err
		direction = directionArr[parsed[3][5:][0]-'0']
	} else {
		length, parsingErr = strconv.Atoi(parsed[2])
		direction = directionMap[parsed[1]]
	}

	if parsingErr != nil {
		return DigInstruction{}, fmt.Errorf("line parsing error: %v", parsingErr)
	}
	return DigInstruction{
		Direction: direction,
		Length:    length,
	}, nil
}

func EvaluateArea(lines []string) int {
	result := 0
	a := image.Point{0, 0}

	for index := 0; index < len(lines); index++ {
		line := lines[index]
		point, _ := ParseInputLine(line)
		length := point.Length
		b := a.Add(point.Direction.Mul(length))

		result += ((a.X*b.Y - a.Y*b.X) + length)
		a = b
	}
	return result/2 + 1
}

func main() {
	file, openError := os.ReadFile(inputFilename)
	if openError != nil {
		log.Panicf("couldn't open input file '%s'\n%v\n", inputFilename, openError)
	}
	lines := strings.Split(string(file), "\n")
	area := EvaluateArea(lines)

	fmt.Printf("result: %d\n", area)
}
