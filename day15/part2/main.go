package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	inputFilename = "../input.txt"
)

type Code uint8

type Lens struct {
	Label string
	Power int
}

func mapStringToCode(s string) Code {
	initialValue := Code(0)

	for charIndex := range s {
		initialValue += Code(s[charIndex])
		initialValue *= 17
	}
	return initialValue
}

func main() {
	file, openError := os.ReadFile(inputFilename)

	if openError != nil {
		log.Fatalln("couldn't open file input.txt")
	}
	codes := strings.Split(string(file[:]), ",")
	boxes := make(map[Code][]Lens)

	for _, code := range codes {
		parts := strings.Split(code, "=")
		label := parts[0]

		if len(label) >= 2 && label[len(label)-1] == '-' {
			label := label[0 : len(label)-1]
			boxIndex := mapStringToCode(label)
			_, mapCreated := boxes[boxIndex]

			if !mapCreated {
				continue
			}
			boxes[boxIndex] = slices.DeleteFunc(boxes[boxIndex], func(l Lens) bool {
				return l.Label == label
			})
		} else if len(parts) == 2 {
			boxIndex := mapStringToCode(label)
			_, mapCreated := boxes[boxIndex]

			if !mapCreated {
				boxes[boxIndex] = make([]Lens, 0)
			}
			if lensId, parseErr := strconv.Atoi(parts[1]); parseErr == nil {
				lens := Lens{
					Label: label,
					Power: lensId,
				}
				if lensIndex := slices.IndexFunc(boxes[boxIndex], func(l Lens) bool {
					return l.Label == label
				}); lensIndex != -1 {
					boxes[boxIndex][lensIndex] = lens
				} else {
					boxes[boxIndex] = append(boxes[boxIndex], lens)
				}
			}
		} else {
			log.Panicf("shouldn't reach this code: label has no '-' nor '={[0-9]}'?? (%s)", code)
		}
	}

	sum := 0
	for boxIndex, box := range boxes {
		for lensIndex, lens := range box {
			sum += (int(boxIndex) + 1) * (lensIndex + 1) * lens.Power
		}
	}
	fmt.Printf("result: %d\n", sum)
}
