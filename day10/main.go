package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

const (
	inputFilename = "input.txt"
)

type Direction int

const (
	DIR_ERROR Direction = -1
	DIR_NORTH Direction = 0
	DIR_EAST  Direction = 1
	DIR_SOUTH Direction = 2
	DIR_WEST  Direction = 3
)

type PipeType rune

const (
	PIP_START PipeType = 'S'
	PIP_VER   PipeType = '|'
	PIP_HOR   PipeType = '-'
	PIP_NTE   PipeType = 'L'
	PIP_NTW   PipeType = 'J'
	PIP_STE   PipeType = 'F'
	PIP_STW   PipeType = '7'
)

type Tile struct {
	Type           PipeType
	X, Y           int
	TunnelProgress int
}

type TunnelMap struct {
	Map         []Tile
	MapHeight   int
	MapWidth    int
	StartingPos *Tile
}

// Returns next tile's coordinate according to direction (eg: tile(X:0, Y:0) -> go(DIR_NORTH) would return (0, -1))
func (t Tile) Go(currentDirection Direction) (x, y int, newDirection Direction) {
	direction := currentDirection

	if direction == DIR_NORTH {
		if t.Type == PIP_STE {
			direction = DIR_EAST
		} else if t.Type == PIP_STW {
			direction = DIR_WEST
		} else if t.Type != PIP_VER {
			direction = DIR_ERROR
		}
	} else if direction == DIR_EAST {
		if t.Type == PIP_NTW {
			direction = DIR_NORTH
		} else if t.Type == PIP_STW {
			direction = DIR_SOUTH
		} else if t.Type != PIP_HOR {
			direction = DIR_ERROR
		}
	} else if direction == DIR_WEST {
		if t.Type == PIP_NTE {
			direction = DIR_NORTH
		} else if t.Type == PIP_STE {
			direction = DIR_SOUTH
		} else if t.Type != PIP_HOR {
			direction = DIR_ERROR
		}
	} else if direction == DIR_SOUTH {
		if t.Type == PIP_NTE {
			direction = DIR_EAST
		} else if t.Type == PIP_NTW {
			direction = DIR_WEST
		} else if t.Type != PIP_VER {
			direction = DIR_ERROR
		}
	}

	switch direction {
	case DIR_NORTH:
		return t.X, t.Y - 1, direction
	case DIR_EAST:
		return t.X + 1, t.Y, direction
	case DIR_SOUTH:
		return t.X, t.Y + 1, direction
	case DIR_WEST:
		return t.X - 1, t.Y, direction
	default:
		return -1, -1, direction
	}
}

// Returns tile pointer, nil if doesn't exist (out of bound)
func (m TunnelMap) tileAt(x, y int) *Tile {
	if 0 <= x && x < m.MapWidth && 0 <= y && y < m.MapHeight {

		return &m.Map[x+m.MapWidth*y]
	}
	return nil
}

// It goes like [N, E, S, W]
// This order is important as it uses Direction type as index value
func (m TunnelMap) findNeighboors(tile *Tile) []*Tile {
	tileUp := m.tileAt(tile.X, tile.Y-1)
	tileDown := m.tileAt(tile.X, tile.Y+1)
	tileLeft := m.tileAt(tile.X-1, tile.Y)
	tileRight := m.tileAt(tile.X+1, tile.Y)

	if tileUp != nil {
		if x, y, _ := tileUp.Go(DIR_NORTH); x < 0 && y < 0 {
			tileUp = nil
		}
	}

	if tileDown != nil {
		if x, y, _ := tileDown.Go(DIR_SOUTH); x < 0 && y < 0 {
			tileDown = nil
		}
	}

	if tileRight != nil {
		if x, y, _ := tileRight.Go(DIR_EAST); x < 0 && y < 0 {
			tileRight = nil
		}
	}

	if tileLeft != nil {
		if x, y, _ := tileLeft.Go(DIR_WEST); x < 0 && y < 0 {
			tileLeft = nil
		}
	}

	return []*Tile{
		tileUp,
		tileRight,
		tileDown,
		tileLeft,
	}
}

func stepNavigationForward(m *TunnelMap, t **Tile, headingDirection *Direction) (bool, int) {
	initialProgress := (*t).TunnelProgress
	x, y, dir := (*t).Go(*headingDirection)

	*t = m.tileAt(x, y)
	*headingDirection = dir

	if (*t).TunnelProgress < 0 {
		(*t).TunnelProgress = initialProgress + 1
		return false, (*t).TunnelProgress
	}
	return true, initialProgress + 1
}

func (m TunnelMap) navigate() (int, error) {
	m.StartingPos.TunnelProgress = 0
	neighboors := m.findNeighboors(m.StartingPos)
	forwardDirection := Direction(slices.IndexFunc(neighboors, func(elem *Tile) bool { return elem != nil }))

	if forwardDirection < 0 {
		return -1, fmt.Errorf("starting tile as no connection")
	}
	backwardDirection := Direction(int(forwardDirection) + 1 + slices.IndexFunc(neighboors[int(forwardDirection)+1:], func(elem *Tile) bool { return elem != nil }))

	if backwardDirection <= forwardDirection {
		return -1, fmt.Errorf("starting tile as only one connection")
	}

	forwardTile := neighboors[forwardDirection]
	backwardTile := neighboors[backwardDirection]
	pathProgress := 1

	forwardTile.TunnelProgress = pathProgress
	backwardTile.TunnelProgress = pathProgress
	for ; forwardTile != backwardTile; pathProgress++ {
		if looped, progress := stepNavigationForward(&m, &forwardTile, &forwardDirection); looped {
			return progress, nil
		}
		if looped, progress := stepNavigationForward(&m, &backwardTile, &backwardDirection); looped {
			return progress, nil
		}
	}
	return pathProgress, nil
}

func initTunnelMap() (TunnelMap, error) {
	file, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines := strings.Split(string(file[:]), "\n")

	if len(fileLines) <= 0 {
		return TunnelMap{}, fmt.Errorf("input file is empty")
	}

	mapWidth := len(fileLines[0])
	mapHeight := len(fileLines)

	tunnelMap := TunnelMap{
		Map:         make([]Tile, mapWidth*mapHeight),
		MapHeight:   mapHeight,
		MapWidth:    mapWidth,
		StartingPos: nil,
	}

	for y, row := range fileLines {
		for x, tile := range row {
			tunnelMap.Map[y*mapWidth+x] = Tile{
				X:              x,
				Y:              y,
				Type:           PipeType(tile),
				TunnelProgress: -1,
			}
			if tile == rune(PIP_START) {
				tunnelMap.StartingPos = &tunnelMap.Map[y*mapWidth+x]
			}
		}
	}
	if tunnelMap.StartingPos == nil {
		return TunnelMap{}, fmt.Errorf("couldn't find tunnel starting position")
	}
	return tunnelMap, nil
}

func main() {
	tunnelMap, error := initTunnelMap()

	if error != nil {
		log.Panicf("Error on file parsing: %v\n", error)
	}

	furthestTileDistance, error := tunnelMap.navigate()

	if error == nil {
		fmt.Printf("result: %d\n", furthestTileDistance)
	} else {
		fmt.Printf("An error occured: %v\n", error)
	}
}
