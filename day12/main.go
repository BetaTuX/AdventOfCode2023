package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	inputFilename = "input.txt"
)

type Instruction struct {
	inputString string
	objective   []int
}

func parseInstruction(line string) (Instruction, error) {
	splitted := strings.Split(line, " ")
	stringAmounts := strings.Split(splitted[1], ",")
	amounts := make([]int, 0, len(stringAmounts))

	if len(splitted) != 2 {
		return Instruction{}, fmt.Errorf("identified more than 2 parts in input, it should be exactly 2 parts separated by a space")
	}
	for _, amountString := range stringAmounts {
		if amount, parseError := strconv.Atoi(amountString); parseError == nil {
			amounts = append(amounts, amount)
		} else {
			return Instruction{}, fmt.Errorf("error parsing amount: %v", parseError)
		}
	}

	return Instruction{
		inputString: splitted[0],
		objective:   amounts,
	}, nil
}

func parseInputFile(filename string) []Instruction {
	file, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines := strings.Split(string(file[:]), "\n")
	instructions := make([]Instruction, 0, len(fileLines))

	for _, line := range fileLines {
		if instruction, error := parseInstruction(line); error == nil {
			instructions = append(instructions, instruction)
		} else {
			log.Panicf("error during parsing: %v\n", error)
		}
	}
	return instructions
}

func identifyArrangements(str string) []int {
	arrangements := make([]int, 0)
	inArrangements := false

	for i := 0; i < len(str); {
		inArrangements = str[i] == '#'

		if inArrangements {
			arrangementLength := strings.IndexRune(str[i:], '.')

			if arrangementLength < 0 {
				arrangementLength = len(str[i:])
			}
			arrangements = append(arrangements, arrangementLength)
			i += arrangementLength
		}
		if nextArrangementIndex := strings.IndexRune(str[i:], '#'); nextArrangementIndex < 0 {
			break
		} else {
			i += nextArrangementIndex
		}
	}
	return arrangements
}

func generateStringFromSeed(input string, seed int) string {
	test := []byte(input)

	for i := range test {
		if test[i] == '?' {
			if seed&0b1 == 1 {
				test[i] = '#'
			} else {
				test[i] = '.'
			}
			seed = seed >> 1
		}
	}
	return string(test)
}

func generatePossibleStrings(input string) []string {
	possibilitiesAmount := int(math.Pow(2.0, float64(strings.Count(input, "?"))))
	possibilities := make([]string, 0, possibilitiesAmount)

	for i := 0; i < possibilitiesAmount; i++ {
		possibilities = append(possibilities, generateStringFromSeed(input, i))
	}
	return possibilities
}

func main() {
	instructions := parseInputFile(inputFilename)

	total := 0
	for _, instruction := range instructions {
		possibilities := generatePossibleStrings(instruction.inputString)

		for _, possibility := range possibilities {
			identifiedArrangement := identifyArrangements(possibility)

			if slices.Equal(instruction.objective, identifiedArrangement) {
				total++
			}
		}
	}
	fmt.Printf("result: %d\n", total)
}
