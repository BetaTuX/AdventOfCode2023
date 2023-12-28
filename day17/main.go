package main

import (
	"container/heap"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"strings"
)

const (
	inputFilename = "input.txt"
)

type QueueItem[T any] struct {
	Value    T
	Priority int
}

type PriorityQueue[T any] []QueueItem[T]

func (q PriorityQueue[_]) Len() int {
	return len(q)
}

func (q PriorityQueue[_]) Less(i, j int) bool {
	return q[i].Priority < q[j].Priority
}

func (q PriorityQueue[_]) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *PriorityQueue[T]) Push(item any) {
	(*q) = append(*q, item.(QueueItem[T]))
}

func (q *PriorityQueue[_]) Pop() any {
	old := *q
	n := len(old)
	item := (*q)[n-1]
	*q = (*q)[0 : n-1]
	return item
}

func (q *PriorityQueue[T]) BetterPush(item T, priority int) {
	heap.Push(q, QueueItem[T]{Value: item, Priority: priority})
}

func (q *PriorityQueue[T]) BetterPop() (T, int) {
	val := heap.Pop(q).(QueueItem[T])

	return val.Value, val.Priority
}

type Cursor struct {
	Coords image.Point
	Dir    image.Point
}

func findPath(grid map[image.Point]int, end image.Point, minMove, maxMove int) int {
	queue, visitedRecord := PriorityQueue[Cursor]{}, map[Cursor]bool{}

	queue.BetterPush(Cursor{Coords: image.Point{0, 0}, Dir: image.Point{0, 1}}, 0)
	queue.BetterPush(Cursor{Coords: image.Point{0, 0}, Dir: image.Point{1, 0}}, 0)

	for len(queue) > 0 {
		cursor, heatloss := queue.BetterPop()

		if cursor.Coords == end {
			return heatloss
		}
		if _, alreadyVisited := visitedRecord[cursor]; alreadyVisited {
			continue
		}
		visitedRecord[cursor] = true
		for i := -maxMove; i <= maxMove; i++ {
			n := cursor.Coords.Add(cursor.Dir.Mul(i))
			if _, ok := grid[n]; !ok || i > -minMove && i < minMove {
				continue
			}
			heatlossStreak, sign := 0, int(math.Copysign(1, float64(i)))
			for j := sign; j != i+sign; j += sign {
				heatlossStreak += grid[cursor.Coords.Add(cursor.Dir.Mul(j))]
			}
			queue.BetterPush(Cursor{n, image.Point{cursor.Dir.Y, cursor.Dir.X}}, heatloss+heatlossStreak)
		}
	}
	return -1
}

func main() {
	file, openError := os.ReadFile(inputFilename)
	if openError != nil {
		log.Panicf("couldn't open input file '%s'\n%v\n", inputFilename, openError)
	}
	lines := strings.Split(string(file), "\n")

	grid, end := map[image.Point]int{}, image.Point{}

	for y, line := range lines {
		for x, chr := range line {
			actualCoord := image.Point{x, y}

			grid[actualCoord] = int(chr - '0')
			end = actualCoord
		}
	}

	fmt.Printf("result (1, 3): %d\n", findPath(grid, end, 1, 3))
	fmt.Printf("result (4, 10): %d\n", findPath(grid, end, 4, 10))
}
