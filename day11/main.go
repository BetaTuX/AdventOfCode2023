package main

import (
	"flag"
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

var (
	shouldUseOlderGalaxies bool
	galaxyOffset           int
)

func init() {
	flag.BoolVar(&shouldUseOlderGalaxies, "older", true, "Use older galaxies")
	flag.Parse()

	if shouldUseOlderGalaxies {
		galaxyOffset = 1000000
	} else {
		galaxyOffset = 2
	}
}

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

func pushStarsAfterColumn(stars []Star, x int) int {
	toPush := make([]*Star, 0)

	for indexOffset := 0; true; indexOffset++ {
		starIndex := slices.IndexFunc(stars[indexOffset:], func(s Star) bool { return s.X > x })
		if starIndex < 0 {
			break
		}
		toPush = append(toPush, &stars[starIndex+indexOffset])
		indexOffset += starIndex
	}
	for _, starPtr := range toPush {
		starPtr.X += galaxyOffset - 1
	}
	return galaxyOffset - 1
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
	yTotalOffset := 0

	for y, line := range fileLines {
		starOnLine := 0
		for x, b := range line {
			if b == '#' {
				starOnLine++
				stars = append(stars, Star{
					X: x,
					Y: y + yTotalOffset,
				})
			}
		}
		if starOnLine == 0 {
			yTotalOffset += galaxyOffset - 1
		}
	}
	xTotalOffset := 0
	for x := 0; x < params.Width; x++ {
		if slices.IndexFunc(stars, func(s Star) bool {
			return s.X-xTotalOffset == x
		}) == -1 {
			xTotalOffset += pushStarsAfterColumn(stars, x+xTotalOffset)
		}
	}
	return GameParams{
		Width:  params.Width + xTotalOffset,
		Height: params.Height + yTotalOffset,
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
