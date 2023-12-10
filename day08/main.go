package main

import (
	"flag"
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

	instructions       string
	nodes              map[string]BTNode
	useGhostNavigation bool
)

func init() {
	flag.BoolVar(&useGhostNavigation, "ghost", true, "Uses ghosts navigation rules")
	flag.Parse()
}

type BTNode struct {
	id          string
	left, right string
}

type NodeGroup []*BTNode

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

func (group NodeGroup) walkLeft() {
	for nodeIndex := range group {
		group[nodeIndex] = group[nodeIndex].walkLeft()
	}
}

func (group NodeGroup) walkRight() {
	for nodeIndex := range group {
		group[nodeIndex] = group[nodeIndex].walkRight()
	}
}

func (group NodeGroup) hasReachedDestination() bool {
	if !useGhostNavigation {
		return group[0].id == "ZZZ"
	} else {
		for _, node := range group {
			if node.id[len(node.id)-1] != 'Z' {
				return false
			}
		}
		return true
	}
}

func parseNode(input string) (BTNode, error) {
	var re = regexp.MustCompile(`(?m)([0-9A-Z]{3}) = \(([0-9A-Z]{3}), ([0-9A-Z]{3})\)$`)
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

func buildNodeGroup() NodeGroup {
	group := make(NodeGroup, 0)

	if !useGhostNavigation {
		node, exists := nodes["AAA"]
		if exists {
			group = append(group, &node)
		}
	} else {
		for nodeIndex := range nodes {
			node := nodes[nodeIndex]

			if node.id[len(node.id)-1] == 'A' {
				group = append(group, &node)
			}
		}
	}
	return group
}

func main() {
	group := buildNodeGroup()
	loop := 0

	for instructionIndex := 0; instructionIndex < len(instructions); {
		instruction := instructions[instructionIndex]

		switch instruction {
		case 'R':
			group.walkRight()
		case 'L':
			group.walkLeft()
		default:
			log.Panicf("instruction unrecognized: %c", instruction)
		}
		loop++
		if group.hasReachedDestination() {
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
