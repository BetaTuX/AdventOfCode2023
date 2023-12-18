package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	inputFilename = "input.txt"
)

var (
	shouldUnfoldInstructions bool
)

func init() {
	flag.BoolVar(&shouldUnfoldInstructions, "unfold", false, "unfold instructions in input.txt")
	flag.Parse()
}

type Instruction struct {
	inputString string
	objective   []int
}

func parseInstruction(line string) (Instruction, error) {
	splitted := strings.Split(line, " ")
	stringAmounts := strings.Split(splitted[1], ",")
	var amounts []int
	var inputString string

	if shouldUnfoldInstructions {
		inputString = fmt.Sprintf("%s?%s?%s?%s?%s", splitted[0], splitted[0], splitted[0], splitted[0], splitted[0])
		amounts = make([]int, len(stringAmounts)*5)
	} else {
		inputString = splitted[0]
		amounts = make([]int, len(stringAmounts))
	}
	if len(splitted) != 2 {
		return Instruction{}, fmt.Errorf("identified more than 2 parts in input, it should be exactly 2 parts separated by a space")
	}
	for index, amountString := range stringAmounts {
		if amount, parseError := strconv.Atoi(amountString); parseError == nil {
			amounts[index] = amount
			if shouldUnfoldInstructions {
				amounts[index+len(stringAmounts)] = amount
				amounts[index+len(stringAmounts)*2] = amount
				amounts[index+len(stringAmounts)*3] = amount
				amounts[index+len(stringAmounts)*4] = amount
			}
		} else {
			return Instruction{}, fmt.Errorf("error parsing amount: %v", parseError)
		}
	}

	return Instruction{
		inputString: inputString,
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

type State [3]int

func mapsClear[M ~map[K]V, K comparable, V any](m M) {
	for k := range m {
		delete(m, k)
	}
}

func countPossibilities(instruction Instruction) int {
	total := 0
	src := []byte(instruction.inputString)
	states := map[State]int{{0, 0, 0}: 1}
	nStates := map[State]int{}

	for _, srcChar := range src {
		for state, quantity := range states {
			challengeSucceeded, hits, requiredEmpty := state[0], state[1], state[2]
			switch {
			case (srcChar == '#' || srcChar == '?') && challengeSucceeded < len(instruction.objective) && requiredEmpty == 0:
				if srcChar == '?' && hits <= 0 {
					nStates[[3]int{challengeSucceeded, hits, requiredEmpty}] += quantity
				}
				hits++
				if hits == instruction.objective[challengeSucceeded] {
					challengeSucceeded, hits, requiredEmpty = challengeSucceeded+1, 0, 1
				}
				nStates[[3]int{challengeSucceeded, hits, requiredEmpty}] += quantity
			case (srcChar == '.' || srcChar == '?') && hits <= 0:
				requiredEmpty = 0
				nStates[[3]int{challengeSucceeded, hits, requiredEmpty}] += quantity
			}
		}
		states, nStates = nStates, states
		mapsClear(nStates)
	}

	for state, amount := range states {
		if state[0] == len(instruction.objective) {
			total += amount
		}
	}
	return total
}

func main() {
	instructions := parseInputFile(inputFilename)

	total := 0
	for _, instruction := range instructions {
		total += countPossibilities(instruction)
	}
	fmt.Printf("result: %d\n", total)
}
