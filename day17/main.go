package main

import (
	"bufio"
	"container/heap"
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

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	h := len(grid)
	w := len(grid[0])
	start := Pos{
		x: 0,
		y: 0,
	}
	end := Pos{
		x: w - 1,
		y: h - 1,
	}
	fmt.Println(search(start, end, grid, w, h))
	fmt.Println(search2(start, end, grid, w, h))
}

type (
	Dir int

	Pos struct {
		x, y int
	}

	State struct {
		pos     Pos
		dir     Dir
		sameDir int
	}
)

const (
	DirNorth = 0
	DirEast  = 1
	DirSouth = 2
	DirWest  = 3
)

func search2(start, end Pos, grid [][]byte, w, h int) int {
	closedSet := NewClosedSet()
	openSet := NewOpenSet()
	openSet.Push(State{
		pos:     start,
		dir:     DirEast,
		sameDir: -1, // start at -1 to counteract initial facing dir
	},
		0,
		manhattanDistance(start, end),
	)
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		if cur.pos == end {
			return curg
		}
		for _, o := range getNeighbors2(cur, w, h) {
			if closedSet.Has(o) {
				continue
			}
			g := curg + int(grid[o.pos.y][o.pos.x]-'0')
			f := g + manhattanDistance(o.pos, end)
			if v, ok := openSet.Get(o); ok {
				if g < v.g {
					openSet.Update(o, g, f)
				}
				continue
			}
			openSet.Push(o, g, f)
		}
	}
	return -1
}

func search(start, end Pos, grid [][]byte, w, h int) int {
	closedSet := NewClosedSet()
	openSet := NewOpenSet()
	openSet.Push(State{
		pos:     start,
		dir:     DirEast,
		sameDir: -1, // start at -1 to counteract initial facing dir
	},
		0,
		manhattanDistance(start, end),
	)
	for !openSet.Empty() {
		cur, curg, _ := openSet.Pop()
		closedSet.Push(cur)
		if cur.pos == end {
			return curg
		}
		for _, o := range getNeighbors(cur, w, h) {
			if closedSet.Has(o) {
				continue
			}
			g := curg + int(grid[o.pos.y][o.pos.x]-'0')
			f := g + manhattanDistance(o.pos, end)
			if v, ok := openSet.Get(o); ok {
				if g < v.g {
					openSet.Update(o, g, f)
				}
				continue
			}
			openSet.Push(o, g, f)
		}
	}
	return -1
}

func getNeighbors2(state State, w, h int) []State {
	opts := make([]State, 0, 3)
	if s, ok := state.goForward2(); ok && isInBounds(s.pos, w, h) {
		opts = append(opts, s)
	}
	if s, ok := state.goLeft2(); ok && isInBounds(s.pos, w, h) {
		opts = append(opts, s)
	}
	if s, ok := state.goRight2(); ok && isInBounds(s.pos, w, h) {
		opts = append(opts, s)
	}
	return opts
}

func getNeighbors(state State, w, h int) []State {
	opts := make([]State, 0, 3)
	if s, ok := state.goForward(); ok && isInBounds(s.pos, w, h) {
		opts = append(opts, s)
	}
	if s := state.goLeft(); isInBounds(s.pos, w, h) {
		opts = append(opts, s)
	}
	if s := state.goRight(); isInBounds(s.pos, w, h) {
		opts = append(opts, s)
	}
	return opts
}

func isInBounds(pos Pos, w, h int) bool {
	return pos.x >= 0 && pos.y >= 0 && pos.x < w && pos.y < h
}

func (s State) goForward2() (State, bool) {
	if s.sameDir >= 9 {
		return State{}, false
	}
	next := s.pos
	switch s.dir {
	case DirNorth:
		next.y--
	case DirEast:
		next.x++
	case DirSouth:
		next.y++
	case DirWest:
		next.x--
	default:
		log.Fatalln("Invalid dir")
	}
	return State{
		pos:     next,
		dir:     s.dir,
		sameDir: s.sameDir + 1,
	}, true
}

func (s State) goLeft2() (State, bool) {
	if s.sameDir >= 0 && s.sameDir < 3 {
		return State{}, false
	}
	next := s.pos
	nextDir := s.dir
	switch s.dir {
	case DirNorth:
		next.x--
		nextDir = DirWest
	case DirEast:
		next.y--
		nextDir = DirNorth
	case DirSouth:
		next.x++
		nextDir = DirEast
	case DirWest:
		next.y++
		nextDir = DirSouth
	default:
		log.Fatalln("Invalid dir")
	}
	return State{
		pos:     next,
		dir:     nextDir,
		sameDir: 0,
	}, true
}

func (s State) goRight2() (State, bool) {
	if s.sameDir >= 0 && s.sameDir < 3 {
		return State{}, false
	}
	next := s.pos
	nextDir := s.dir
	switch s.dir {
	case DirNorth:
		next.x++
		nextDir = DirEast
	case DirEast:
		next.y++
		nextDir = DirSouth
	case DirSouth:
		next.x--
		nextDir = DirWest
	case DirWest:
		next.y--
		nextDir = DirNorth
	default:
		log.Fatalln("Invalid dir")
	}
	return State{
		pos:     next,
		dir:     nextDir,
		sameDir: 0,
	}, true
}

func (s State) goForward() (State, bool) {
	if s.sameDir >= 2 {
		return State{}, false
	}
	next := s.pos
	switch s.dir {
	case DirNorth:
		next.y--
	case DirEast:
		next.x++
	case DirSouth:
		next.y++
	case DirWest:
		next.x--
	default:
		log.Fatalln("Invalid dir")
	}
	return State{
		pos:     next,
		dir:     s.dir,
		sameDir: s.sameDir + 1,
	}, true
}

func (s State) goLeft() State {
	next := s.pos
	nextDir := s.dir
	switch s.dir {
	case DirNorth:
		next.x--
		nextDir = DirWest
	case DirEast:
		next.y--
		nextDir = DirNorth
	case DirSouth:
		next.x++
		nextDir = DirEast
	case DirWest:
		next.y++
		nextDir = DirSouth
	default:
		log.Fatalln("Invalid dir")
	}
	return State{
		pos:     next,
		dir:     nextDir,
		sameDir: 0,
	}
}

func (s State) goRight() State {
	next := s.pos
	nextDir := s.dir
	switch s.dir {
	case DirNorth:
		next.x++
		nextDir = DirEast
	case DirEast:
		next.y++
		nextDir = DirSouth
	case DirSouth:
		next.x--
		nextDir = DirWest
	case DirWest:
		next.y--
		nextDir = DirNorth
	default:
		log.Fatalln("Invalid dir")
	}
	return State{
		pos:     next,
		dir:     nextDir,
		sameDir: 0,
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func manhattanDistance(a, b Pos) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

type (
	Item struct {
		value State
		g, f  int
		index int
	}

	PriorityQueue struct {
		q []*Item
		s map[State]int
	}

	OpenSet struct {
		q *PriorityQueue
	}

	ClosedSet map[State]struct{}
)

func NewOpenSet() *OpenSet {
	return &OpenSet{
		q: NewPriorityQueue(),
	}
}

func (s *OpenSet) Empty() bool {
	return s.q.Len() == 0
}

func (s *OpenSet) Has(val State) bool {
	_, ok := s.q.s[val]
	return ok
}

func (s *OpenSet) Get(val State) (*Item, bool) {
	idx, ok := s.q.s[val]
	if !ok {
		return nil, false
	}
	return s.q.q[idx], true
}

func (s *OpenSet) Push(value State, g, f int) {
	heap.Push(s.q, &Item{
		value: value,
		g:     g,
		f:     f,
	})
}

func (s *OpenSet) Pop() (State, int, int) {
	item := heap.Pop(s.q).(*Item)
	return item.value, item.g, item.f
}

func (s *OpenSet) Update(value State, g, f int) bool {
	return s.q.Update(value, g, f)
}

func NewClosedSet() ClosedSet {
	return ClosedSet{}
}

func (cs ClosedSet) Has(val State) bool {
	_, ok := cs[val]
	return ok
}

func (cs ClosedSet) Push(val State) {
	cs[val] = struct{}{}
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		s: map[State]int{},
	}
}

func (q PriorityQueue) Len() int { return len(q.q) }
func (q PriorityQueue) Less(i, j int) bool {
	return q.q[i].f < q.q[j].f
}

func (q PriorityQueue) Swap(i, j int) {
	q.q[i], q.q[j] = q.q[j], q.q[i]
	q.q[i].index = i
	q.q[j].index = j
	q.s[q.q[i].value] = i
	q.s[q.q[j].value] = j
}

func (q *PriorityQueue) Push(x interface{}) {
	n := len(q.q)
	item := x.(*Item)
	item.index = n
	q.q = append(q.q, item)
	q.s[item.value] = n
}

func (q *PriorityQueue) Pop() interface{} {
	n := len(q.q)
	item := q.q[n-1]
	q.q[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	q.q = q.q[:n-1]
	delete(q.s, item.value)
	return item
}

func (q *PriorityQueue) Update(value State, g, f int) bool {
	idx, ok := q.s[value]
	if !ok {
		return false
	}
	item := q.q[idx]
	item.g = g
	item.f = f
	heap.Fix(q, item.index)
	return true
}
