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
	PIP_EMPTY PipeType = '.'
	PIP_VER   PipeType = '|'
	PIP_HOR   PipeType = '-'
	PIP_NTE   PipeType = 'L'
	PIP_NTW   PipeType = 'J'
	PIP_STE   PipeType = 'F'
	PIP_STW   PipeType = '7'
)

type Color int8

const (
	COLOR_UNMARKED Color = 0
	COLOR_RED      Color = 1
	COLOR_BLUE     Color = -1
)

type Tile struct {
	Type           PipeType
	X, Y           int
	TunnelProgress int
	mark           Color
}

type TunnelMap struct {
	Map         []Tile
	MapHeight   int
	MapWidth    int
	StartingPos *Tile
}

func (t Tile) To(direction Direction) (x, y int) {
	switch direction {
	case DIR_NORTH:
		return t.X, t.Y - 1
	case DIR_EAST:
		return t.X + 1, t.Y
	case DIR_SOUTH:
		return t.X, t.Y + 1
	case DIR_WEST:
		return t.X - 1, t.Y
	default:
		return -1, -1
	}
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

	x, y = t.To(direction)
	return x, y, direction
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

func stepNavigationForward(m *TunnelMap, t **Tile, headingDirection *Direction, execOnTile func(*Tile) (bool, int)) (bool, int) {
	x, y, dir := (*t).Go(*headingDirection)

	*t = m.tileAt(x, y)
	*headingDirection = dir

	return execOnTile(*t)
}

func (m *TunnelMap) identifyForward() (*Tile, Direction, []*Tile) {
	neighboors := m.findNeighboors(m.StartingPos)
	forwardDirection := Direction(slices.IndexFunc(neighboors, func(elem *Tile) bool { return elem != nil }))

	return neighboors[forwardDirection], forwardDirection, neighboors
}

func (m *TunnelMap) navigate() (int, error) {
	m.StartingPos.TunnelProgress = 0
	forwardTile, forwardDirection, neighboors := m.identifyForward()

	tilePathProgress := func(actualProgress int) func(*Tile) (bool, int) {
		return func(t *Tile) (bool, int) {
			if (*t).TunnelProgress < 0 {
				(*t).TunnelProgress = actualProgress + 1
				return false, (*t).TunnelProgress
			}
			return true, (*t).TunnelProgress
		}
	}

	if forwardDirection < 0 {
		return -1, fmt.Errorf("starting tile as no connection")
	}
	backwardDirection := Direction(int(forwardDirection) + 1 + slices.IndexFunc(neighboors[int(forwardDirection)+1:], func(elem *Tile) bool { return elem != nil }))

	if backwardDirection <= forwardDirection {
		return -1, fmt.Errorf("starting tile as only one connection")
	}

	backwardTile := neighboors[backwardDirection]
	pathProgress := 1

	forwardTile.TunnelProgress = pathProgress
	backwardTile.TunnelProgress = pathProgress
	for ; forwardTile != backwardTile; pathProgress++ {
		if looped, progress := stepNavigationForward(m, &forwardTile, &forwardDirection, tilePathProgress(pathProgress)); looped {
			return progress, nil
		}
		if looped, progress := stepNavigationForward(m, &backwardTile, &backwardDirection, tilePathProgress(pathProgress)); looped {
			return progress, nil
		}
	}
	return pathProgress, nil
}

func (m *TunnelMap) identifyStartTileType() PipeType {
	neighboors := m.findNeighboors(m.StartingPos)

	switch {
	case neighboors[0] == neighboors[2] && neighboors[2] == nil:
		return PIP_HOR
	case neighboors[1] == neighboors[3] && neighboors[3] == nil:
		return PIP_VER
	case neighboors[0] == neighboors[1] && neighboors[1] == nil:
		return PIP_STW
	case neighboors[0] == neighboors[3] && neighboors[3] == nil:
		return PIP_STE
	case neighboors[2] == neighboors[3] && neighboors[3] == nil:
		return PIP_NTE
	case neighboors[2] == neighboors[1] && neighboors[1] == nil:
		return PIP_NTW
	default:
		return PIP_EMPTY
	}
}

func (m *TunnelMap) markZones() Color {
	lastFromTop := false
	isInside := false

	for tileIndex := range m.Map {
		currentTile := &m.Map[tileIndex]

		if tileIndex%m.MapWidth == 0 {
			isInside = false
		}
		if currentTile.TunnelProgress < 0 {
			if isInside {
				currentTile.mark = COLOR_RED
			} else {
				currentTile.mark = COLOR_BLUE
			}
		} else {
			currentType := currentTile.Type

			if currentType == PIP_START {
				currentType = m.identifyStartTileType()
			}
			if currentType == PIP_HOR ||
				(lastFromTop && (currentType == PIP_STE || currentType == PIP_STW)) ||
				(!lastFromTop && (currentType == PIP_NTE || currentType == PIP_NTW)) {
				continue
			}
			isInside = !isInside
		}
	}
	return COLOR_RED
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
				mark:           COLOR_UNMARKED,
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
	identifiedColor := tunnelMap.markZones()

	enclosedTiles := 0

	for index, tile := range tunnelMap.Map {
		switch tile.mark {
		case COLOR_UNMARKED:
			print("X")
		case COLOR_RED:
			print("r")
			if identifiedColor == COLOR_RED {
				enclosedTiles++
			}
		case COLOR_BLUE:
			print("b")
			if identifiedColor == COLOR_BLUE {
				enclosedTiles++
			}
		}
		if (index+1)%tunnelMap.MapWidth == 0 {
			println("")
		}
	}

	if error == nil {
		fmt.Printf("result: %d\nenclosed tiles: %d\n", furthestTileDistance, enclosedTiles)
	} else {
		fmt.Printf("An error occured: %v\n", error)
	}
}
