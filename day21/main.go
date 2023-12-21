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

	var start Pos
	var grid [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if pos := strings.Index(line, "S"); pos >= 0 {
			start = Pos{
				x: pos,
				y: len(grid),
			}
		}
		grid = append(grid, []byte(line))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", getNum(grid, start, 64))

	height := len(grid)
	width := len(grid[0])
	if height != width {
		log.Fatalln("Grid is not square")
	}
	target := 26501365
	multiple := target / height
	rem := target % height
	a0 := getNum(grid, start, rem)
	a1 := getNum(grid, start, height+rem)
	a2 := getNum(grid, start, height*2+rem)

	delta1 := a1 - a0
	delta2 := a2 - a1
	delta3 := delta2 - delta1
	a := delta3 / 2
	b := delta1 - 3*a
	c := a0 - a - b

	seqNum := multiple + 1

	fmt.Println("Part 2:", seqNum*seqNum*a+seqNum*b+c)
}

func getNum(grid [][]byte, start Pos, target int) int {
	closedSet := map[Pos]struct{}{}
	openSet := NewStateVec()
	openSet.Push(State{pos: start, g: 0})
	for !openSet.IsEmpty() {
		s, _ := openSet.Pop()
		if s.g >= target {
			closedSet[s.pos] = struct{}{}
			continue
		}
		for _, i := range getNeighbors(grid, s.pos) {
			openSet.Push(State{
				pos: i,
				g:   s.g + 1,
			})
		}
	}
	return len(closedSet)
}

type (
	State struct {
		pos Pos
		g   int
	}

	Pos struct {
		x, y int
	}
)

func getNeighbors(grid [][]byte, pos Pos) []Pos {
	res := make([]Pos, 0, 4)
	{
		v := Pos{
			x: pos.x,
			y: pos.y - 1,
		}
		if getOnTorus(grid, v) != '#' {
			res = append(res, v)
		}
	}
	{
		v := Pos{
			x: pos.x - 1,
			y: pos.y,
		}
		if getOnTorus(grid, v) != '#' {
			res = append(res, v)
		}
	}
	{
		v := Pos{
			x: pos.x,
			y: pos.y + 1,
		}
		if getOnTorus(grid, v) != '#' {
			res = append(res, v)
		}
	}
	{
		v := Pos{
			x: pos.x + 1,
			y: pos.y,
		}
		if getOnTorus(grid, v) != '#' {
			res = append(res, v)
		}
	}
	return res
}

func getOnTorus(grid [][]byte, pos Pos) byte {
	height := len(grid)
	width := len(grid[0])
	y := pos.y % height
	if y < 0 {
		y += height
	}
	x := pos.x % width
	if x < 0 {
		x += width
	}
	return grid[y][x]
}

type (
	StateVec struct {
		states  []State
		holding map[State]struct{}
	}
)

func NewStateVec() *StateVec {
	return &StateVec{
		holding: map[State]struct{}{},
	}
}

func (s *StateVec) Push(v State) {
	if _, ok := s.holding[v]; ok {
		return
	}
	s.states = append(s.states, v)
	s.holding[v] = struct{}{}
}

func (s *StateVec) Pop() (State, bool) {
	if len(s.states) == 0 {
		return State{}, false
	}
	v := s.states[len(s.states)-1]
	s.states = s.states[:len(s.states)-1]
	return v, true
}

func (s *StateVec) IsEmpty() bool {
	return len(s.states) == 0
}
