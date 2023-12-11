package main

import (
	"bufio"
	"fmt"
	"log"
	"math/bits"
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
		Rem       int
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
		count++
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
						Count: count,
					}
				} else if !revisitNodes[i].Candidate {
					if node.ID != revisitNodes[i].ID {
						log.Fatalln("Multiple terminals for cycle")
					}
					cycle := count - revisitNodes[i].Count
					revisitNodes[i] = NodeVisit{
						ID:        revisitNodes[i].ID,
						Count:     count,
						Cycle:     cycle,
						Rem:       count % cycle,
						Candidate: true,
					}
					totalRevisits++
				} else {
					if node.ID != revisitNodes[i].ID {
						log.Fatalln("Multiple terminals for cycle")
					}
					if cycle := count - revisitNodes[i].Count; cycle != revisitNodes[i].Cycle {
						log.Fatalln("Multiple cycle lengths")
					}
					revisitNodes[i] = NodeVisit{
						ID:        revisitNodes[i].ID,
						Count:     count,
						Cycle:     revisitNodes[i].Cycle,
						Rem:       revisitNodes[i].Rem,
						Candidate: true,
					}
				}
			}
		}
	}
	a := revisitNodes[0].Rem
	m := revisitNodes[0].Cycle
	for _, i := range revisitNodes[1:] {
		var ok bool
		a, m, ok = crt(a, m, i.Rem, i.Cycle)
		if !ok {
			log.Fatalln("Unsolvable constraints")
		}
	}
	if a <= 0 {
		a += m
	}
	fmt.Println("Part 2:", a)
}

func crt(a1, m1, a2, m2 int) (int, int, bool) {
	g, p, q := extGCD(m1, m2)
	if a1%g != a2%g {
		return 0, 0, false
	}
	m1g := m1 / g
	m2g := m2 / g
	lcm := m1g * m2
	// a1 * m2/g * q + a2 * m1/g * p (mod lcm)
	x := (mulmod(mulmod(a1, m2g, lcm), q, lcm) + mulmod(mulmod(a2, m1g, lcm), p, lcm)) % lcm
	if x < 0 {
		x += lcm
	}
	return x, lcm, true
}

func extGCD(a, b int) (int, int, int) {
	x2 := 1
	x1 := 0
	y2 := 0
	y1 := 1
	// a should be larger than b
	flip := false
	if a < b {
		a, b = b, a
		flip = true
	}
	for b > 0 {
		q := a / b
		a, b = b, a%b
		x2, x1 = x1, x2-q*x1
		y2, y1 = y1, y2-q*y1
	}
	if flip {
		x2, y2 = y2, x2
	}
	return a, x2, y2
}

func mulmod(a, b, m int) int {
	sign := 1
	if a < 0 {
		a = -a
		sign *= -1
	}
	if b < 0 {
		b = -b
		sign *= -1
	}
	a = a % m
	b = b % m
	hi, lo := bits.Mul(uint(a), uint(b))
	return sign * int(bits.Rem(hi, lo, uint(m)))
}
