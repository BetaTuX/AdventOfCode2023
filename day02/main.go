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
	isSecondPart                                   bool
	redBallsLimit, blueBallsLimit, greenBallsLimit int
	gameRegex                                      = regexp.MustCompile(`(?m)Game (?P<game>[0-9]+): (?P<line>.*)$`)
)

func init() {
	flag.IntVar(&redBallsLimit, "reds-limit", 12, "Sets reds ball amount limit (default: 12)")
	flag.IntVar(&greenBallsLimit, "green-limit", 13, "Sets green ball amount limit (default: 13)")
	flag.IntVar(&blueBallsLimit, "blue-limit", 14, "Sets blue ball amount limit (default: 14)")
	flag.BoolVar(&isSecondPart, "alternative-calculations", false, "Uses ball amount power instead of game number")
	flag.Parse()
}

func checkBallAmountIsValid(red, green, blue int) bool {
	return red <= redBallsLimit && green <= greenBallsLimit && blue <= blueBallsLimit
}

func countBalls(turns []string) (int, int, int) {
	redBalls, greenBalls, blueBalls := 0, 0, 0
	ballBinding := map[string]*int{
		"green": &greenBalls,
		"red":   &redBalls,
		"blue":  &blueBalls,
	}

	for _, turn := range turns {
		details := strings.Split(turn, ",")

		for _, ballDetails := range details {
			ballDetails = strings.Trim(ballDetails, " ")
			// parts[0] should be the ball amount, [1] should be the color
			parts := strings.Split(ballDetails, " ")
			parsedBallAmount, _ := strconv.Atoi(parts[0])
			savedBallAmount := ballBinding[parts[1]]

			if *savedBallAmount < parsedBallAmount {
				*savedBallAmount = parsedBallAmount
			}
		}
	}
	return redBalls, greenBalls, blueBalls
}

func main() {
	file, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines := strings.Split(string(file[:]), "\n")
	games := make([]int, 0, len(fileLines))
	regexRes := gameRegex.FindAllStringSubmatch(string(file[:]), -1)

	for _, line := range regexRes {
		gameNb, _ := strconv.Atoi(line[1])
		gameTurns := strings.Split(line[2], ";")
		redBallAmount, greenBallAmount, blueBallAmount := countBalls(gameTurns)

		if !isSecondPart && checkBallAmountIsValid(redBallAmount, greenBallAmount, blueBallAmount) {
			games = append(games, gameNb)
		} else if isSecondPart {
			games = append(games, redBallAmount*greenBallAmount*blueBallAmount)
		}
	}

	sum := 0
	for _, gamePower := range games {
		sum += gamePower
	}
	fmt.Printf("Game result: %d\n", sum)
}
