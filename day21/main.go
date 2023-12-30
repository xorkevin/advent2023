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

	height := len(grid)
	width := len(grid[0])
	const target = 26501365
	multiple := target / height
	rem := target % height

	innerEven := 0
	innerOdd := 0
	cornerEven := 0
	cornerOdd := 0

	closedSet := make([]bool, width*height)
	openSet := NewRing[State](width * height)
	closedSet[start.y*width+start.x] = true
	openSet.Write(State{
		pos: start,
		g:   0,
	})
	const p1Target = 64
	const p1TargetIsEven = p1Target%2 == 0
	sum := 0
	for {
		s, ok := openSet.Read()
		if !ok {
			break
		}
		closedSet[s.pos.y*width+s.pos.x] = true
		curIsEven := s.g%2 == 0
		if manhattanDistance(s.pos, start) > rem {
			if curIsEven {
				cornerEven++
			} else {
				cornerOdd++
			}
		} else {
			if curIsEven {
				innerEven++
			} else {
				innerOdd++
			}
		}
		if s.g <= p1Target && curIsEven == p1TargetIsEven {
			sum++
		}
		var neighbors [4]Pos
		n := getNeighbors(grid, width, height, s.pos, neighbors[:])
		for _, i := range neighbors[:n] {
			key := i.y*width + i.x
			if closedSet[key] {
				continue
			}
			closedSet[key] = true
			openSet.Write(State{
				pos: i,
				g:   s.g + 1,
			})
		}
	}
	fmt.Println("Part 1:", sum)

	if height != width {
		log.Fatalln("Grid is not square")
	}
	if height%2 != 1 || start.y != (height-1)/2 || start.x != start.y {
		log.Fatalln("Start is not centered")
	}

	const targetIsEven = target%2 == 0
	multipleIsEven := multiple%2 == 0
	multiple1 := multiple + 1
	outerMultiple := multiple1 * multiple1
	innerMultiple := multiple * multiple

	outerDiamond := innerOdd
	innerDiamond := innerEven
	outerCorner := cornerOdd
	innerCorner := cornerEven
	if targetIsEven == multipleIsEven {
		outerDiamond, innerDiamond = innerDiamond, outerDiamond
		outerCorner, innerCorner = innerCorner, outerCorner
	}
	fmt.Println("Part 2:", outerMultiple*outerDiamond+innerMultiple*innerDiamond+(outerMultiple-multiple1)*outerCorner+(innerMultiple+multiple)*innerCorner)
}

func manhattanDistance(a, b Pos) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
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

func getNeighbors(grid [][]byte, w, h int, pos Pos, res []Pos) int {
	count := 0
	{
		v := pos
		v.y--
		if inBounds(v, w, h) && grid[v.y][v.x] != '#' {
			res[count] = v
			count++
		}
	}
	{
		v := pos
		v.x--
		if inBounds(v, w, h) && grid[v.y][v.x] != '#' {
			res[count] = v
			count++
		}
	}
	{
		v := pos
		v.y++
		if inBounds(v, w, h) && grid[v.y][v.x] != '#' {
			res[count] = v
			count++
		}
	}
	{
		v := pos
		v.x++
		if inBounds(v, w, h) && grid[v.y][v.x] != '#' {
			res[count] = v
			count++
		}
	}
	return count
}

func inBounds(pos Pos, w, h int) bool {
	return pos.x >= 0 && pos.y >= 0 && pos.x < w && pos.y < h
}

type (
	Ring[T any] struct {
		buf []T
		r   int
		w   int
	}
)

func NewRing[T any](size int) *Ring[T] {
	if size < 2 {
		size = 2
	}
	return &Ring[T]{
		buf: make([]T, size),
		r:   0,
		w:   0,
	}
}

func (b *Ring[T]) resize() {
	next := make([]T, len(b.buf)*2)
	if b.r == b.w {
		b.w = 0
	} else if b.r < b.w {
		b.w = copy(next, b.buf[b.r:b.w])
	} else {
		p := copy(next, b.buf[b.r:])
		q := 0
		if b.w > 0 {
			q = copy(next[p:], b.buf[:b.w])
		}
		b.w = p + q
	}
	b.buf = next
	b.r = 0
}

func (b *Ring[T]) Write(m T) {
	next := (b.w + 1) % len(b.buf)
	if next == b.r {
		b.resize()
		b.Write(m)
		return
	}
	b.buf[b.w] = m
	b.w = next
}

func (b *Ring[T]) Read() (T, bool) {
	if b.r == b.w {
		var v T
		return v, false
	}
	next := (b.r + 1) % len(b.buf)
	m := b.buf[b.r]
	b.r = next
	return m, true
}

func (b *Ring[T]) Peek() (T, bool) {
	if b.r == b.w {
		var v T
		return v, false
	}
	m := b.buf[b.r]
	return m, true
}
