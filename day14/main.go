package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	inputFilename = "input.txt"
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

func main() {
	file, openError := os.ReadFile(inputFilename)

	if openError != nil {
		log.Fatalln("couldn't open file input.txt")
	}
	linesAsString := strings.Split(string(file[:]), "\n")
	fileLines := make([][]byte, 0, len(linesAsString))

	for _, line := range linesAsString {
		fileLines = append(fileLines, []byte(line))
	}
	// Rotate Right
	rotateMatrix(fileLines)

	sum := 0
	for _, line := range fileLines {
		// Roll balls to the end of line
		BubbleSort(line, RollBalls)
		for index, char := range line {
			if char == 'O' {
				// Use index in line (actually column) to determine ball score
				sum += index + 1
			}
		}
	}
	fmt.Printf("result: %d\n", sum)
}
