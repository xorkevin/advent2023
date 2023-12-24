package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}

	height := len(grid)
	width := len(grid[0])
	start := Pos{
		x: 0,
		y: 0,
	}
	start.x = bytes.IndexByte(grid[0], '.')
	if start.x < 0 {
		log.Fatalln("No start")
	}
	end := Pos{
		x: 0,
		y: height - 1,
	}
	end.x = bytes.IndexByte(grid[height-1], '.')
	if end.x < 0 {
		log.Fatalln("No end")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	undirectedGraph, directedGraph := contractPaths(start, end, grid, width, height)
	startID := start.y*width + start.x
	endID := end.y*width + end.x
	startID, startCost := findBranch(startID, undirectedGraph)
	endID, endCost := findBranch(endID, undirectedGraph)
	prunedCost := startCost + endCost
	fmt.Println("Part 1:", prunedCost+searchLongestGraph(startID, endID, directedGraph, width, height))
	fmt.Println("Part 2:", prunedCost+searchLongestGraph(startID, endID, undirectedGraph, width, height))
}

func searchLongestGraph(start, end int, graph map[int]map[int]int, w, h int) int {
	maxG := -1
	closedSet := make([]bool, w*h)
	openSet := NewStack[State]()
	openSet.Write(State{
		id:    start,
		g:     0,
		depth: 0,
	})
	curPath := NewStack[State]()
	for {
		cur, ok := openSet.Read()
		if !ok {
			break
		}

		if cur.id == end {
			maxG = max(maxG, cur.g)
			continue
		}

		for curPath.Len() > cur.depth {
			k, ok := curPath.Read()
			if !ok {
				log.Fatalln("cur path in bad state")
			}
			closedSet[k.id] = false
		}
		curPath.Write(cur)
		closedSet[cur.id] = true

		for k, v := range graph[cur.id] {
			if closedSet[k] {
				continue
			}
			openSet.Write(State{
				id:    k,
				g:     cur.g + v,
				depth: cur.depth + 1,
			})
		}
	}
	return maxG
}

func findBranch(pos int, graph map[int]map[int]int) (int, int) {
	edges := graph[pos]
	if len(edges) != 1 {
		return pos, 0
	}
	for k, v := range edges {
		return k, v
	}
	log.Fatalln("Unreachable")
	return -1, -1
}

func contractPaths(start, end Pos, grid [][]byte, w, h int) (undirected, directed map[int]map[int]int) {
	contractedGraph := map[int]map[int]int{}
	contractedDirectedGraph := map[int]map[int]int{}
	closedSet := make([]bool, w*h)
	openSet := NewRing[Pos](w * h)
	openSet.Write(start)
	exploredSet := make([]bool, w*h)
	exploredSet[start.y*w+start.x] = true
	var edges [4]Edge
	for {
		cur, ok := openSet.Read()
		if !ok {
			break
		}

		curKey := cur.y*w + cur.x
		closedSet[curKey] = true
		n := getPathEdges(cur, end, grid, w, h, closedSet, edges[:])
		closedSet[curKey] = false

		if cur != start && n < 2 && edges[0].pos != end {
			// eliminate dead ends
			continue
		}

		{
			var graphEdges map[int]int
			if m, ok := contractedGraph[curKey]; ok {
				graphEdges = m
			} else {
				graphEdges = map[int]int{}
				contractedGraph[curKey] = graphEdges
			}
			for _, o := range edges[:n] {
				key := o.pos.y*w + o.pos.x
				if o.pos != end && !exploredSet[key] {
					openSet.Write(o.pos)
					exploredSet[key] = true
				}
				if v, ok := graphEdges[key]; ok {
					if o.cost > v {
						graphEdges[key] = o.cost
						contractedGraph[key][curKey] = o.cost
					}
				} else {
					graphEdges[key] = o.cost
					var revGraphEdges map[int]int
					if m, ok := contractedGraph[key]; ok {
						revGraphEdges = m
					} else {
						revGraphEdges = map[int]int{}
						contractedGraph[key] = revGraphEdges
					}
					revGraphEdges[curKey] = o.cost
				}
			}
		}
		{
			var graphEdges map[int]int
			if m, ok := contractedDirectedGraph[curKey]; ok {
				graphEdges = m
			} else {
				graphEdges = map[int]int{}
				contractedDirectedGraph[curKey] = graphEdges
			}
			for _, o := range edges[:n] {
				key := o.pos.y*w + o.pos.x
				if o.forward {
					if v, ok := graphEdges[key]; ok {
						if o.cost > v {
							graphEdges[key] = o.cost
						}
					} else {
						graphEdges[key] = o.cost
					}
				}
				if o.rev {
					var revGraphEdges map[int]int
					if m, ok := contractedDirectedGraph[key]; ok {
						revGraphEdges = m
					} else {
						revGraphEdges = map[int]int{}
						contractedDirectedGraph[key] = revGraphEdges
					}
					if v, ok := revGraphEdges[curKey]; ok {
						if o.cost > v {
							revGraphEdges[curKey] = o.cost
						}
					} else {
						revGraphEdges[curKey] = o.cost
					}
				}
			}
		}
	}
	return contractedGraph, contractedDirectedGraph
}

type (
	State struct {
		id    int
		g     int
		depth int
	}

	Edge struct {
		pos     Pos
		cost    int
		forward bool
		rev     bool
	}

	Pos struct {
		x, y int
	}
)

func getPathEdges(pos, end Pos, grid [][]byte, w, h int, closedSet []bool, res []Edge) int {
	var unwindQueue []int
	count := getEdges(pos, grid, w, h, closedSet, res)
	var edges [4]Edge
	for idx, o := range res[:count] {
		cur := o
		for {
			if cur.pos == end {
				for _, i := range unwindQueue {
					closedSet[i] = false
				}
				unwindQueue = unwindQueue[:0]
				res[idx] = cur
				break
			}
			key := cur.pos.y*w + cur.pos.x
			closedSet[key] = true
			unwindQueue = append(unwindQueue, key)
			n := getEdges(cur.pos, grid, w, h, closedSet, edges[:])
			if n != 1 {
				for _, i := range unwindQueue {
					closedSet[i] = false
				}
				unwindQueue = unwindQueue[:0]
				res[idx] = cur
				break
			}
			first := edges[0]
			cur = Edge{
				pos:     first.pos,
				cost:    cur.cost + first.cost,
				forward: cur.forward && first.forward,
				rev:     cur.rev && first.rev,
			}
		}
	}
	return count
}

func getEdges(pos Pos, grid [][]byte, w, h int, closedSet []bool, res []Edge) int {
	count := 0
	cur := grid[pos.y][pos.x]
	{
		v := pos
		v.y--
		if isInBounds(v, w, h) && !closedSet[v.y*w+v.x] {
			if b := grid[v.y][v.x]; b != '#' {
				res[count] = Edge{
					pos:     v,
					cost:    1,
					forward: cur == '.' || cur == '^',
					rev:     b == '.' || b == 'v',
				}
				count++
			}
		}
	}
	{
		v := pos
		v.x--
		if isInBounds(v, w, h) && !closedSet[v.y*w+v.x] {
			if b := grid[v.y][v.x]; b != '#' {
				res[count] = Edge{
					pos:     v,
					cost:    1,
					forward: cur == '.' || cur == '<',
					rev:     b == '.' || b == '>',
				}
				count++
			}
		}
	}
	{
		v := pos
		v.y++
		if isInBounds(v, w, h) && !closedSet[v.y*w+v.x] {
			if b := grid[v.y][v.x]; b != '#' {
				res[count] = Edge{
					pos:     v,
					cost:    1,
					forward: cur == '.' || cur == 'v',
					rev:     b == '.' || b == '^',
				}
				count++
			}
		}
	}
	{
		v := pos
		v.x++
		if isInBounds(v, w, h) && !closedSet[v.y*w+v.x] {
			if b := grid[v.y][v.x]; b != '#' {
				res[count] = Edge{
					pos:     v,
					cost:    1,
					forward: cur == '.' || cur == '>',
					rev:     b == '.' || b == '<',
				}
				count++
			}
		}
	}
	return count
}

func isInBounds(pos Pos, w, h int) bool {
	return pos.x >= 0 && pos.y >= 0 && pos.x < w && pos.y < h
}

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

func (b *Stack[T]) Peek() (T, bool) {
	if len(b.buf) == 0 {
		var v T
		return v, false
	}
	top := len(b.buf) - 1
	m := b.buf[top]
	return m, true
}

func (b *Stack[T]) Len() int {
	return len(b.buf)
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
