package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
)

const (
	inputFilename = "input.txt"
)

var (
	shouldSearchForMax bool
)

func init() {
	flag.BoolVar(&shouldSearchForMax, "maximize", false, "Search for the maximum energized tiles")
	flag.Parse()
}

type Direction uint8

const (
	DIR_RIGHT Direction = 0
	DIR_DOWN  Direction = 1
	DIR_LEFT  Direction = 2
	DIR_UP    Direction = 3
)

func (d Direction) RotateClockwise() Direction {
	return (d + 1) % 4
}

func (d Direction) RotateCounterClockwise() Direction {
	returnValue := d - 1

	if returnValue > 3 {
		return DIR_UP
	}
	return returnValue
}

func (d Direction) IsVertical() bool {
	return d == DIR_UP || d == DIR_DOWN
}

func (d Direction) IsHorizontal() bool {
	return d == DIR_LEFT || d == DIR_RIGHT
}

type Cursor struct {
	x, y      int
	direction Direction
}

func (c *Cursor) Move() {
	switch c.direction {
	case DIR_DOWN:
		c.y++
	case DIR_UP:
		c.y--
	case DIR_RIGHT:
		c.x++
	case DIR_LEFT:
		c.x--
	}
}

func (c Cursor) IsInBoundary(m MirrorMap) bool {
	return c.x >= 0 && c.y >= 0 && c.x < m.GetWidth() && c.y < m.GetHeight()
}

func (c Cursor) IsLooping(m MirrorMap) bool {
	activeTile := &m[c.y][c.x]

	return activeTile.visitedFrom[c.direction]
}

type MirrorTile struct {
	tile        ETile
	energized   bool
	visitedFrom [4]bool
}

type ETile byte

const (
	TIL_EMPTY           ETile = '.'
	TIL_MIRROR_FORWARD  ETile = '/'
	TIL_MIRROR_BACKWARD ETile = '\\'
	TIL_SPLITTER_VER    ETile = '|'
	TIL_SPLITTER_HOR    ETile = '-'
)

func (t ETile) MapDirection(input Direction) []Direction {
	if slices.Contains([]ETile{TIL_EMPTY}, t) {
		return []Direction{input}
	} else if (t == TIL_SPLITTER_HOR && input.IsHorizontal()) || (t == TIL_SPLITTER_VER && input.IsVertical()) {
		return []Direction{input}
	} else if t == TIL_SPLITTER_HOR {
		return []Direction{DIR_LEFT, DIR_RIGHT}
	} else if t == TIL_SPLITTER_VER {
		return []Direction{DIR_UP, DIR_DOWN}
	} else if t == TIL_MIRROR_FORWARD {
		if input.IsVertical() {
			return []Direction{input.RotateClockwise()}
		} else {
			return []Direction{input.RotateCounterClockwise()}
		}
	} else if t == TIL_MIRROR_BACKWARD {
		if input.IsVertical() {
			return []Direction{input.RotateCounterClockwise()}
		} else {
			return []Direction{input.RotateClockwise()}
		}
	}
	return []Direction{}
}

type MirrorMap [][]MirrorTile

func newMirrorMapFromByteArray(b []byte) MirrorMap {
	m := make(MirrorMap, 0)

	for len(b) > 0 {
		cutIndex := slices.Index(b, '\n')

		if cutIndex != -1 {
			tileArr := make([]MirrorTile, cutIndex)
			for index := range tileArr {
				tileArr[index] = MirrorTile{
					tile:        ETile(b[index]),
					energized:   false,
					visitedFrom: [4]bool{false, false, false, false},
				}
			}
			b = b[cutIndex+1:]
			m = append(m, tileArr)
		} else {
			tileArr := make([]MirrorTile, len(b))
			for index := range tileArr {
				tileArr[index] = MirrorTile{
					tile:        ETile(b[index]),
					energized:   false,
					visitedFrom: [4]bool{false, false, false, false},
				}
			}
			m = append(m, tileArr)
			break
		}
	}
	return m
}

func (m MirrorMap) Display() {
	for _, line := range m {
		for _, tile := range line {
			if tile.energized {
				fmt.Printf("#")
			} else {
				fmt.Printf("%c", tile.tile)
			}
		}
		fmt.Printf("\n")
	}
}

func (m MirrorMap) CountEnergized() int {
	sum := 0

	for _, line := range m {
		for _, tile := range line {
			if tile.energized {
				sum++
			}
		}
	}
	return sum
}

func (m MirrorMap) Reset() int {
	sum := 0

	for lineIndex := range m {
		for tileIndex := range m[lineIndex] {
			m[lineIndex][tileIndex].energized = false
			m[lineIndex][tileIndex].visitedFrom = [4]bool{false, false, false, false}
		}
	}
	return sum
}

func (m MirrorMap) GetWidth() int {
	if len(m) <= 0 {
		return 0
	}
	return len(m[0])
}

func (m MirrorMap) GetHeight() int {
	return len(m)
}

func (m MirrorMap) RunSimulation(startCursor Cursor) {
	cursorArray := make([]Cursor, 1)
	cursorArray[0] = startCursor

	for len(cursorArray) > 0 {
		nCursorArray := make([]Cursor, 0, len(cursorArray))

		for _, cursor := range cursorArray {
			cursorActiveTile := &(m[cursor.y][cursor.x])
			newDirections := cursorActiveTile.tile.MapDirection(cursor.direction)

			for _, newDirection := range newDirections {
				newCursor := Cursor{
					x:         cursor.x,
					y:         cursor.y,
					direction: newDirection,
				}
				newCursor.Move()

				if newCursor.IsInBoundary(m) && !newCursor.IsLooping(m) {
					nCursorArray = append(nCursorArray, newCursor)
				}
			}
			cursorActiveTile.energized = true
			cursorActiveTile.visitedFrom[cursor.direction] = true
		}
		cursorArray = nCursorArray
	}
}

func (m MirrorMap) SearchMax() int {

	mapWidth := m.GetWidth()
	mapHeight := m.GetHeight()
	perimeter := mapWidth*2 + mapHeight*2

	max := 0
	for i := 0; i < perimeter; i++ {
		var x, y int
		var dir Direction

		if i < mapWidth*2 {
			x = i % mapWidth
			if i < mapWidth {
				y = 0
				dir = DIR_DOWN
			} else {
				y = mapHeight - 1
				dir = DIR_UP
			}
		} else {
			y = (i - (mapWidth * 2)) % mapHeight
			if i-(mapWidth*2) < mapHeight {
				dir = DIR_RIGHT
				x = 0
			} else {
				dir = DIR_LEFT
				x = mapWidth - 1
			}
		}

		m.RunSimulation(Cursor{
			x:         x,
			y:         y,
			direction: dir,
		})
		if energized := m.CountEnergized(); energized > max {
			max = energized
		}
		m.Reset()
	}
	return max
}

func main() {
	file, openError := os.ReadFile(inputFilename)

	if openError != nil {
		log.Fatalln("couldn't open file input.txt")
	}
	mirrorMap := newMirrorMapFromByteArray(file[:])

	if shouldSearchForMax {
		fmt.Printf("max result: %d\n", mirrorMap.SearchMax())
	} else {
		mirrorMap.RunSimulation(Cursor{
			x:         0,
			y:         0,
			direction: DIR_RIGHT,
		})
		fmt.Printf("result: %d\n", mirrorMap.CountEnergized())

	}
}
