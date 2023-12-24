package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	inputFilename   = "input.txt"
	DEFAULT_MAXLOOP = 1000000000
)

func rotateMatrix[T any](matrix [][]T) [][]T {
	// reverse the matrix
	for i, j := 0, len(matrix)-1; i < j; i, j = i+1, j-1 {
		matrix[i], matrix[j] = matrix[j], matrix[i]
	}

	// transpose it
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < i; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}
	return matrix
}

// Reimplemented sorting to make sure of the sorting algorithm
func BubbleSort[T any](array []T, sortFx func(a, b T) bool) []T {
	for i := 0; i < len(array)-1; i++ {
		for j := 0; j < len(array)-i-1; j++ {
			if sortFx(array[j], array[j+1]) {
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
	}
	return array
}
func RollBalls(a, b byte) bool {
	return b == '.' && a == 'O'
}

func evaluateBallWeight(m [][]byte) int {
	sum := 0
	for _, line := range m {
		for index, char := range line {
			if char == 'O' {
				// Use index in line (actually column) to determine ball score
				sum += index + 1
			}
		}
	}
	return sum
}

func cleanRecurrenceMap(r map[[4]int]int) {
	for tuple, occurence := range r {
		if occurence <= 1 {
			delete(r, tuple)
		}
	}
}

func findEndLoopValue(start, loopLength, end int) int {
	return loopLength - (((end - start) / loopLength) % loopLength)
}

func main() {
	file, openError := os.ReadFile(inputFilename)

	if openError != nil {
		log.Fatalln("couldn't open file input.txt")
	}
	linesAsString := strings.Split(string(file[:]), "\n")
	fileLines := make([][]byte, 0, len(linesAsString))
	recurrenceMap := make(map[[4]int]int)

	for _, line := range linesAsString {
		fileLines = append(fileLines, []byte(line))
	}
	// Rotate Right (North on right)
	rotateMatrix(fileLines)

	loopLimit := DEFAULT_MAXLOOP

	for i := 0; i < loopLimit; i++ {
		sequenceTuple := [4]int{}

		for direction := 0; direction < 4; direction++ {
			for _, line := range fileLines {
				// Roll balls to the end of line
				BubbleSort(line, RollBalls)
			}
			rotateMatrix(fileLines)
			evaluated := evaluateBallWeight(fileLines)
			sequenceTuple[direction] = evaluated
		}
		recurrenceMap[sequenceTuple]++
		if recurrenceMap[sequenceTuple] > 2 && loopLimit == DEFAULT_MAXLOOP {
			// Remove every occurence that happens only once (not in the loop)
			cleanRecurrenceMap(recurrenceMap)
			// We are already in the loop so remove 1 to i to place cursor before the beginning of the loop
			evaluatedSolutionIndex := findEndLoopValue(i-1, len(recurrenceMap), DEFAULT_MAXLOOP)
			loopLimit = i + evaluatedSolutionIndex
		}
	}
	sum := evaluateBallWeight(fileLines)
	fmt.Printf("result: %d\n", sum)
}
