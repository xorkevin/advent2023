package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
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

	connGraph := map[string]map[string]int{}
	scanner := bufio.NewScanner(file)
	var edges [][2]string
	for scanner.Scan() {
		lhs, rhs, ok := strings.Cut(scanner.Text(), ": ")
		if !ok {
			log.Fatalln("Invalid line")
		}
		var conns map[string]int
		if m, ok := connGraph[lhs]; ok {
			conns = m
		} else {
			conns = map[string]int{}
			connGraph[lhs] = conns
		}
		for _, i := range strings.Split(rhs, " ") {
			conns[i] = 1
			if m, ok := connGraph[i]; ok {
				m[lhs] = 1
			} else {
				m := map[string]int{}
				m[lhs] = 1
				connGraph[i] = m
			}
			edges = append(edges, [2]string{lhs, i})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	connNames := make([]string, 0, len(connGraph))
	for k := range connGraph {
		connNames = append(connNames, k)
	}
	slices.Sort(connNames)

	slices.SortFunc(edges, func(a, b [2]string) int {
		sa := len(connGraph[a[0]]) + len(connGraph[a[1]])
		sb := len(connGraph[b[0]]) + len(connGraph[b[1]])
		return sb - sa
	})

	openSet := NewRing[string](len(connGraph))
	closedSet := map[string]struct{}{}
	product, count := rmCandidatesAndCount([3][2]string{{"zlx", "chr"}, {"cpq", "hlx"}, {"hqp", "spk"}}, connNames, connGraph, openSet, closedSet)
	if count != 2 {
		log.Fatalln("Invalid cut")
	}
	fmt.Println("Part 1:", product)
}

func rmCandidatesAndCount(candidates [3][2]string, connNames []string, connGraph map[string]map[string]int, openSet *Ring[string], closedSet map[string]struct{}) (int, int) {
	for _, i := range candidates {
		if _, ok := connGraph[i[0]][i[1]]; !ok {
			return -1, -1
		}
		if _, ok := connGraph[i[1]][i[0]]; !ok {
			return -1, -1
		}
	}
	for _, i := range candidates {
		delete(connGraph[i[0]], i[1])
		delete(connGraph[i[1]], i[0])
	}
	defer func() {
		for _, i := range candidates {
			connGraph[i[0]][i[1]] = 1
			connGraph[i[1]][i[0]] = 1
		}
	}()
	prevSize := 0
	count := 0
	product := 1
	for {
		start := slices.IndexFunc(connNames, func(e string) bool {
			_, ok := closedSet[e]
			return !ok
		})
		if start < 0 {
			return product, count
		}
		addReachable(connNames[start], openSet, closedSet, connGraph)
		count++
		product *= len(closedSet) - prevSize
		prevSize = len(closedSet)
		if count >= 2 {
			if len(closedSet) < len(connNames) {
				return -1, -1
			}
			return product, count
		}
	}
}

func getNumSets(connNames []string, connGraph map[string]map[string]int) int {
	openSet := NewRing[string](len(connGraph))
	closedSet := map[string]struct{}{}
	count := 0
	for {
		start := slices.IndexFunc(connNames, func(e string) bool {
			_, ok := closedSet[e]
			return !ok
		})
		if start < 0 {
			return count
		}
		addReachable(connNames[start], openSet, closedSet, connGraph)
		count++
	}
}

func addReachable(start string, openSet *Ring[string], closedSet map[string]struct{}, connGraph map[string]map[string]int) {
	openSet.Write(start)
	closedSet[start] = struct{}{}
	for {
		cur, ok := openSet.Read()
		if !ok {
			return
		}
		for k := range connGraph[cur] {
			if _, ok := closedSet[k]; !ok {
				closedSet[k] = struct{}{}
				openSet.Write(k)
			}
		}
	}
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
