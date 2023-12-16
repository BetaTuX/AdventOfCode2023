package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strings"
)

const (
	inputFilename = "input.txt"
)

type Star struct {
	X, Y int
}

type GameParams struct {
	Width, Height int
}

func evaluateDistanceBetweenStars(start, destination Star) int {
	xOffset := int(math.Abs(float64(destination.X - start.X)))
	yOffset := int(math.Abs(float64(destination.Y - start.Y)))

	return xOffset + yOffset
}

func pushStarsAfterColumn(stars []Star, x int) {
	toPush := make([]*Star, 0)

	for offset := 0; true; offset++ {
		starIndex := slices.IndexFunc(stars[offset:], func(s Star) bool { return s.X > x })
		if starIndex < 0 {
			break
		}
		toPush = append(toPush, &stars[starIndex+offset])
		offset += starIndex
	}
	for _, starPtr := range toPush {
		starPtr.X++
	}
}

func initStarIndex() (GameParams, []Star) {
	file, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines := strings.Split(string(file[:]), "\n")
	stars := make([]Star, 0)
	params := GameParams{
		Height: len(fileLines),
		Width:  len(fileLines[0]),
	}
	yOffset := 0

	for y, line := range fileLines {
		starOnLine := 0
		for x, b := range line {
			if b == '#' {
				starOnLine++
				stars = append(stars, Star{
					X: x,
					Y: y + yOffset,
				})
			}
		}
		if starOnLine == 0 {
			yOffset++
		}
	}
	xOffset := 0
	for x := 0; x < params.Width; x++ {
		if slices.IndexFunc(stars, func(s Star) bool {
			return s.X-xOffset == x
		}) == -1 {
			pushStarsAfterColumn(stars, x+xOffset)
			xOffset++
		}
	}
	return GameParams{
		Width:  params.Width + xOffset,
		Height: params.Height + yOffset,
	}, stars
}

func main() {
	_, stars := initStarIndex()
	totalDistance := 0

	for refIndex := range stars {
		for starIndex := range stars[refIndex+1:] {
			totalDistance += evaluateDistanceBetweenStars(stars[refIndex], stars[refIndex+starIndex+1])
		}
	}
	fmt.Printf("result: %d\n", totalDistance)
}
