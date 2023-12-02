package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	inputFile = "./calibration_input.txt"
)

var (
	redBalls, blueBalls, greenBalls int
	ballBinding                     = map[string]*int{
		"green": &greenBalls,
		"red":   &redBalls,
		"blue":  &blueBalls,
	}
	gameRegex = regexp.MustCompile(`(?m)Game (?P<game>[0-9]+): (?P<line>.*)$`)
)

func init() {
	flag.IntVar(&redBalls, "reds-limit", 12, "Sets reds ball amount limit (default: 12)")
	flag.IntVar(&greenBalls, "green-limit", 13, "Sets green ball amount limit (default: 13)")
	flag.IntVar(&blueBalls, "blue-limit", 14, "Sets blue ball amount limit (default: 14)")
	flag.Parse()
}

func analyseLine(turns []string) bool {
	for _, turn := range turns {
		details := strings.Split(turn, ",")

		for _, ballDetails := range details {
			ballDetails = strings.Trim(ballDetails, " ")
			// parts[0] should be the ball amount, [1] should be the color
			parts := strings.Split(ballDetails, " ")
			ballAmount, _ := strconv.Atoi(parts[0])
			ballType := ballBinding[parts[1]]

			if *ballType < ballAmount {
				return false
			}
		}
	}
	return true
}

func main() {
	file, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines := strings.Split(string(file[:]), "\n")
	validLines := make([]int, 0, len(fileLines))
	regexRes := gameRegex.FindAllStringSubmatch(string(file[:]), -1)

	for _, line := range regexRes {
		gameNb, _ := strconv.Atoi(line[1])
		gameTurns := strings.Split(line[2], ";")
		isGameValid := analyseLine(gameTurns)

		if isGameValid {
			validLines = append(validLines, gameNb)
		}
	}

	sum := 0
	for _, gameNb := range validLines {
		sum += gameNb
	}
	fmt.Printf("Game result: %d\n", sum)
}
