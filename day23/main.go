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

	var grid [][]byte
	var start Pos
	var end Pos
	first := true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if first {
			first = false
			x := strings.Index(scanner.Text(), ".")
			if x < 0 {
				log.Fatalln("No start")
			}
			start = Pos{
				x: x,
				y: 0,
			}
		} else {
			x := strings.Index(scanner.Text(), ".")
			if x < 0 {
				log.Fatalln("No end")
			}
			end = Pos{
				x: x,
				y: len(grid),
			}
		}
		grid = append(grid, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	height := len(grid)
	width := len(grid[0])
	fmt.Println(searchLongest(start, end, grid, width, height))
	fmt.Println(searchLongest2(start, end, grid, width, height))
}

func searchLongest(start, end Pos, grid [][]byte, w, h int) int {
	maxG := -1
	closedSet := map[Pos]struct{}{}
	openSet := NewStack[State]()
	openSet.Write(State{
		pos: start,
		g:   0,
	})
	curPath := NewStack[Pos]()
	for {
		cur, ok := openSet.Read()
		if !ok {
			break
		}

		if cur.pos == end {
			maxG = max(maxG, cur.g)
			continue
		}

		for curPath.Len() > cur.g {
			k, ok := curPath.Read()
			if !ok {
				log.Fatalln("cur path in bad state")
			}
			delete(closedSet, k)
		}
		curPath.Write(cur.pos)
		closedSet[cur.pos] = struct{}{}

		for _, o := range getNeighbors(cur.pos, grid, w, h) {
			if _, ok := closedSet[o]; ok {
				continue
			}
			openSet.Write(State{
				pos: o,
				g:   cur.g + 1,
			})
		}
	}
	return maxG
}

func searchLongest2(start, end Pos, grid [][]byte, w, h int) int {
	maxG := -1
	closedSet := map[Pos]struct{}{}
	openSet := NewStack[State]()
	openSet.Write(State{
		pos: start,
		g:   0,
	})
	curPath := NewStack[Pos]()
	for {
		cur, ok := openSet.Read()
		if !ok {
			break
		}

		if cur.pos == end {
			maxG = max(maxG, cur.g)
			continue
		}

		for curPath.Len() > cur.g {
			k, ok := curPath.Read()
			if !ok {
				log.Fatalln("cur path in bad state")
			}
			delete(closedSet, k)
		}
		curPath.Write(cur.pos)
		closedSet[cur.pos] = struct{}{}

		neighbors := getNeighbors2(cur.pos, grid, w, h)
		if len(neighbors) == 1 {
			fmt.Println("is dead end")
		}
		for _, o := range neighbors {
			if _, ok := closedSet[o]; ok {
				continue
			}
			openSet.Write(State{
				pos: o,
				g:   cur.g + 1,
			})
		}
	}
	return maxG
}

type (
	State struct {
		pos Pos
		g   int
	}
)

func getNeighbors(pos Pos, grid [][]byte, w, h int) []Pos {
	opts := make([]Pos, 0, 4)
	cur := grid[pos.y][pos.x]
	{
		v := pos
		v.y--
		if isInBounds(v, w, h) {
			if grid[v.y][v.x] != '#' && (cur == '.' || cur == '^') {
				opts = append(opts, v)
			}
		}
	}
	{
		v := pos
		v.x--
		if isInBounds(v, w, h) {
			if grid[v.y][v.x] != '#' && (cur == '.' || cur == '<') {
				opts = append(opts, v)
			}
		}
	}
	{
		v := pos
		v.y++
		if isInBounds(v, w, h) {
			if grid[v.y][v.x] != '#' && (cur == '.' || cur == 'v') {
				opts = append(opts, v)
			}
		}
	}
	{
		v := pos
		v.x++
		if isInBounds(v, w, h) {
			if grid[v.y][v.x] != '#' && (cur == '.' || cur == '>') {
				opts = append(opts, v)
			}
		}
	}
	return opts
}

func getNeighbors2(pos Pos, grid [][]byte, w, h int) []Pos {
	opts := make([]Pos, 0, 4)
	{
		v := pos
		v.y--
		if isInBounds(v, w, h) {
			if grid[v.y][v.x] != '#' {
				opts = append(opts, v)
			}
		}
	}
	{
		v := pos
		v.x--
		if isInBounds(v, w, h) {
			if grid[v.y][v.x] != '#' {
				opts = append(opts, v)
			}
		}
	}
	{
		v := pos
		v.y++
		if isInBounds(v, w, h) {
			if grid[v.y][v.x] != '#' {
				opts = append(opts, v)
			}
		}
	}
	{
		v := pos
		v.x++
		if isInBounds(v, w, h) {
			if grid[v.y][v.x] != '#' {
				opts = append(opts, v)
			}
		}
	}
	return opts
}

func isInBounds(pos Pos, w, h int) bool {
	return pos.x >= 0 && pos.y >= 0 && pos.x < w && pos.y < h
}

type (
	Pos struct {
		x, y int
	}
)

type (
	Stack[T any] struct {
		buf []T
	}
)

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (b *Stack[T]) Write(m T) {
	b.buf = append(b.buf, m)
}

func (b *Stack[T]) Read() (T, bool) {
	if len(b.buf) == 0 {
		var v T
		return v, false
	}
	top := len(b.buf) - 1
	m := b.buf[top]
	b.buf = b.buf[:top]
	return m, true
}

func (b *Stack[T]) Len() int {
	return len(b.buf)
}
