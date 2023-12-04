package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	x, y int
}

type EngineNumber struct {
	value    int
	position Position
	length   int
	marked   bool
}

const (
	inputFile = "./calibration_input.txt"
)

var (
	shouldFindGearRatios bool

	numbers []EngineNumber
	symbols []Position
)

func init() {
	flag.BoolVar(&shouldFindGearRatios, "find-gear-ratios", false, "Finds gear ratio instead of every engine numbers")
	flag.Parse()
}

func (symbol Position) findAdjacent(ratios []EngineNumber) []*EngineNumber {
	results := make([]*EngineNumber, 0)

	for index := range ratios {
		if ratios[index].isAdjacent(symbol) {
			results = append(results, &ratios[index])
		}
	}
	return results
}

func (number EngineNumber) isAdjacent(pos Position) bool {
	numberLeftBound := number.position.x
	numberRightBound := number.position.x + number.length - 1
	isHorizontalAdjacent := pos.y-1 <= number.position.y && number.position.y <= pos.y+1
	isVerticalAdjacent := pos.x-1 <= numberRightBound && numberLeftBound <= pos.x+1

	return isHorizontalAdjacent && isVerticalAdjacent
}

func markAdjacents(engine []EngineNumber, pos Position) {
	for index := range engine {
		if engine[index].isAdjacent(pos) {
			engine[index].marked = true
		}
	}
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isEngineSymbol(r rune) bool {
	return r != '.' && !isDigit(r)
}

func cutString(base string, isValid func(rune) bool) string {
	for index, char := range base {
		if !isValid(char) {
			return base[0:index]
		}
	}
	return base
}

func parseDigit(slice string, x, y int) int {
	numberString := cutString(slice[x:], isDigit)
	numberLength := len(numberString)

	if number, err := strconv.Atoi(numberString); err == nil {
		numbers = append(numbers, EngineNumber{
			value: number,
			position: Position{
				x: x,
				y: y,
			},
			length: numberLength,
		})
		return numberLength - 1
	} else {
		fmt.Printf("error during number parsing: %s\n", err)
	}
	return 0
}

func parseLine(line string, lineIndex int) {
	slice := line

	for i := 0; i < len(slice); i++ {
		char := rune(slice[i])

		if char == '.' {
			continue
		}
		if isDigit(char) {
			i += parseDigit(slice, i, lineIndex)
		} else if isEngineSymbol(char) {
			symbols = append(symbols, Position{
				x: i,
				y: lineIndex,
			})
		}
	}
}

func main() {
	file, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines := strings.Split(string(file[:]), "\n")

	for rowIndex, row := range fileLines {
		parseLine(row, rowIndex)
	}
	sum := 0
	for _, pos := range symbols {
		if !shouldFindGearRatios {
			markAdjacents(numbers, pos)
		} else if fileLines[pos.y][pos.x] == '*' {
			if adjacents := pos.findAdjacent(numbers); len(adjacents) == 2 {
				sum += adjacents[0].value * adjacents[1].value
			}
		}
	}

	if !shouldFindGearRatios {
		for _, v := range numbers {
			if v.marked {
				sum += v.value
			}
		}
	}

	println("Result :", sum)
}
