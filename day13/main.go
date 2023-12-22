package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
)

const (
	inputFilename = "input.txt"
)

type GroundMap []string

func (m GroundMap) checkHorizontalMirrorAt(index int) bool {
	for i := 0; (index+i < len(m)) && (index-i > 0); i++ {
		if m[index+i] != m[index-(i+1)] {
			return false
		}
	}
	return index > 0
}

func (m GroundMap) makeColumnBuffer() []byte {
	return make([]byte, len(m))
}

func (m GroundMap) fillColumnBuffer(index int, buff []byte) {
	for i := 0; i < len(m); i++ {
		buff[i] = m[i][index]
	}
}

func (m GroundMap) checkVerticalMirrorAt(index int) bool {
	verticalBuffer := [][]byte{
		m.makeColumnBuffer(),
		m.makeColumnBuffer(),
	}
	for i := 0; (index+i < len(m[0])) && (index-i > 0); i++ {
		m.fillColumnBuffer(index-(i+1), verticalBuffer[0])
		m.fillColumnBuffer(index+i, verticalBuffer[1])
		if !slices.Equal(verticalBuffer[0], verticalBuffer[1]) {
			return false
		}
	}
	return index > 0
}

func (m GroundMap) Solve() (bool, int) {
	verticalLimit := len(m)
	horizontalLimit := len(m[0])
	limit := int(math.Max(float64(verticalLimit), float64(horizontalLimit)))
	hLineBuffer := [1]string{}
	vLineBuffer := [2][]byte{
		m.makeColumnBuffer(),
		m.makeColumnBuffer(),
	}

	for i := 0; i < limit; i++ {
		if i < verticalLimit {
			if hLineBuffer[0] == m[i] && m.checkHorizontalMirrorAt(i) {
				return true, i * 100
			}
			hLineBuffer[0] = m[i]
		}
		if i < horizontalLimit {
			m.fillColumnBuffer(i, vLineBuffer[1])
			if slices.Equal(vLineBuffer[0], vLineBuffer[1]) && m.checkVerticalMirrorAt(i) {
				return true, i
			}
			vLineBuffer[0] = slices.Insert(vLineBuffer[0][:0], 0, vLineBuffer[1]...)
		}
	}
	return false, -1
}

func main() {
	file, openError := os.Open(inputFilename)

	if openError != nil {
		log.Fatalln("couldn't open file input.txt")
	}
	reader := bufio.NewScanner(file)
	lineBuffer := make([]string, 0)
	results := make([]int, 0)

	for i := 0; true; i++ {
		shouldExit := !reader.Scan()
		line := reader.Text()

		if line == "" {
			groundMap := GroundMap(lineBuffer)

			if hasMirror, index := groundMap.Solve(); hasMirror {
				results = append(results, index)
			} else {
				log.Printf("no mirror? %d\n", i)
			}
			lineBuffer = lineBuffer[:0]
			i = -1
		} else {
			lineBuffer = append(lineBuffer, line)
		}
		if shouldExit {
			break
		}
	}

	sum := 0
	for _, value := range results {
		sum += value
	}
	fmt.Printf("result: %d\n", sum)
}
