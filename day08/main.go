package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	inputFilename = "input.txt"
)

var (
	fileLines []string

	instructions string
	nodes        map[string]BTNode
)

type BTNode struct {
	id          string
	left, right string
}

func (node *BTNode) walkLeft() *BTNode {
	leftNode, exists := nodes[node.left]

	if exists {
		return &leftNode
	} else {
		return nil
	}
}

func (node *BTNode) walkRight() *BTNode {
	rightNode, exists := nodes[node.right]

	if exists {
		return &rightNode
	} else {
		return nil
	}
}

func parseNode(input string) (BTNode, error) {
	var re = regexp.MustCompile(`(?m)([A-Z]{3}) = \(([A-Z]{3}), ([A-Z]{3})\)$`)
	regexResult := re.FindStringSubmatch(input)

	if len(regexResult) != 4 {
		return BTNode{}, fmt.Errorf("parsing of line resulted in missing parts")
	}

	return BTNode{
		id:    regexResult[1],
		left:  regexResult[2],
		right: regexResult[3],
	}, nil
}

func init() {
	file, err := os.ReadFile(inputFilename)
	if err != nil {
		log.Fatalln("couldn't open input file\n", err)
	}
	fileLines = strings.Split(string(file[:]), "\n")

	instructions = fileLines[0]

	nodes = make(map[string]BTNode)
	for _, line := range fileLines[2:] {
		node, err := parseNode(line)

		if err == nil {
			nodes[node.id] = node
		} else {
			log.Panicf("parsing error: err: %s", err)
		}
	}
}

func main() {
	startingNode := nodes["AAA"]
	currentNode := &startingNode
	loop := 0

	for instructionIndex := 0; instructionIndex < len(instructions); {
		instruction := instructions[instructionIndex]

		switch instruction {
		case 'R':
			currentNode = currentNode.walkRight()
		case 'L':
			currentNode = currentNode.walkLeft()
		default:
			log.Panicf("instruction unrecognized: %c", instruction)
		}
		loop++
		if currentNode.id == "ZZZ" {
			break
		}
		if instructionIndex == len(instructions)-1 {
			instructionIndex = 0
		} else {
			instructionIndex++
		}
	}

	fmt.Printf("result: %d\n", loop)
}
