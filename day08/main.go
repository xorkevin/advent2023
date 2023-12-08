package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	Node struct {
		ID    string
		Left  string
		Right string
		IsEnd bool
	}

	NodeVisit struct {
		ID        string
		Count     int
		Cycle     int
		Candidate bool
	}
)

func main() {
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	var steps []byte
	nodes := map[string]Node{}
	var starts []string

	first := true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if first {
			first = false
			steps = []byte(scanner.Text())
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		lhs, rhs, ok := strings.Cut(line, " = ")
		if !ok {
			log.Fatalln("Invalid line")
		}
		a, b, ok := strings.Cut(strings.Trim(rhs, "()"), ", ")
		if !ok {
			log.Fatalln("Invalid line")
		}
		nodes[lhs] = Node{
			ID:    lhs,
			Left:  a,
			Right: b,
			IsEnd: strings.HasSuffix(lhs, "Z"),
		}
		if strings.HasSuffix(lhs, "A") {
			starts = append(starts, lhs)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	start, ok := nodes["AAA"]
	if !ok {
		log.Fatalln("No start node", "AAA")
	}
	end := "ZZZ"
	count := 0
	for start.ID != end {
		instr := steps[count%len(steps)]
		next := start.Left
		if instr == 'R' {
			next = start.Right
		}
		var ok bool
		start, ok = nodes[next]
		if !ok {
			log.Fatalln("Next node not found")
		}
		count++
	}
	fmt.Println("Part 1:", count)

	startNodes := make([]Node, 0, len(starts))
	for _, i := range starts {
		node, ok := nodes[i]
		if !ok {
			log.Fatalln("No start node", i)
		}
		startNodes = append(startNodes, node)
	}
	count = 0
	revisitNodes := make([]NodeVisit, len(startNodes))
	totalRevisits := 0
	totalStarts := len(startNodes)
	for totalRevisits < totalStarts {
		instr := steps[count%len(steps)]
		for i := range startNodes {
			next := startNodes[i].Left
			if instr == 'R' {
				next = startNodes[i].Right
			}
			node, ok := nodes[next]
			if !ok {
				log.Fatalln("Next node not found")
			}
			startNodes[i] = node
			if node.IsEnd {
				if revisitNodes[i].ID == "" {
					revisitNodes[i] = NodeVisit{
						ID:    node.ID,
						Count: count + 1,
					}
				} else if !revisitNodes[i].Candidate {
					if node.ID != revisitNodes[i].ID {
						log.Fatalln("Multiple terminals for cycle")
					}
					cycle := (count + 1) - revisitNodes[i].Count
					revisitNodes[i] = NodeVisit{
						ID:        revisitNodes[i].ID,
						Count:     count + 1,
						Cycle:     cycle,
						Candidate: true,
					}
					totalRevisits++
				} else {
					if node.ID != revisitNodes[i].ID {
						log.Fatalln("Multiple terminals for cycle")
					}
					if cycle := (count + 1) - revisitNodes[i].Count; cycle != revisitNodes[i].Cycle {
						log.Fatalln("Multiple cycle lengths")
					}
					revisitNodes[i] = NodeVisit{
						ID:        revisitNodes[i].ID,
						Count:     count + 1,
						Cycle:     revisitNodes[i].Cycle,
						Candidate: true,
					}
				}
			}
		}
		count++
	}
	l := 1
	for _, i := range revisitNodes {
		l = lcm(l, i.Cycle)
	}
	fmt.Println("Part 2:", l)
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	// b should be larger than a
	if a > b {
		a, b = b, a
	}
	for a != 0 {
		a, b = b%a, a
	}
	return b
}
