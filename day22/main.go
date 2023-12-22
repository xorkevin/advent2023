package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
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

	var lines []Line

	first := true
	minX := 0
	maxX := 0
	minY := 0
	maxY := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lhs, rhs, ok := strings.Cut(line, "~")
		if !ok {
			log.Fatalln("Invalid line")
		}
		var lhsPos Pos
		var rhsPos Pos
		lhsNumStrs := strings.Split(lhs, ",")
		if len(lhsNumStrs) != 3 {
			log.Fatalln("Invalid line")
		}
		var err error
		lhsPos.x, err = strconv.Atoi(lhsNumStrs[0])
		if err != nil {
			log.Fatalln(err)
		}
		lhsPos.y, err = strconv.Atoi(lhsNumStrs[1])
		if err != nil {
			log.Fatalln(err)
		}
		lhsPos.z, err = strconv.Atoi(lhsNumStrs[2])
		if err != nil {
			log.Fatalln(err)
		}
		rhsNumStrs := strings.Split(rhs, ",")
		if len(rhsNumStrs) != 3 {
			log.Fatalln("Invalid line")
		}
		rhsPos.x, err = strconv.Atoi(rhsNumStrs[0])
		if err != nil {
			log.Fatalln(err)
		}
		rhsPos.y, err = strconv.Atoi(rhsNumStrs[1])
		if err != nil {
			log.Fatalln(err)
		}
		rhsPos.z, err = strconv.Atoi(rhsNumStrs[2])
		if err != nil {
			log.Fatalln(err)
		}

		if posLess(rhsPos, lhsPos) < 0 {
			lhsPos, rhsPos = rhsPos, lhsPos
		}
		lines = append(lines, Line{
			a:      lhsPos,
			b:      rhsPos,
			height: rhsPos.z - lhsPos.z,
		})

		if first {
			first = false
			minX = lhsPos.x
			minY = lhsPos.y
			maxX = rhsPos.x
			maxY = rhsPos.y
		} else {
			minX = min(minX, lhsPos.x)
			minY = min(minY, lhsPos.y)
			maxX = max(maxX, rhsPos.x)
			maxY = max(maxY, rhsPos.y)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	slices.SortFunc(lines, func(a, b Line) int {
		return posLess(a.a, b.a)
	})

	xwidth := maxX - minX + 1
	ywidth := maxY - minY + 1

	count := 0
	sum := 0
	heightMap := make([]int, xwidth*ywidth)
	fullTower := make([]Line, len(lines))
	getTower(minX, minY, xwidth, heightMap, lines, -1, fullTower)
	candidate := make([]Line, len(lines))
	for n := range lines {
		getTower(minX, minY, xwidth, heightMap, lines, n, candidate)
		delta := towerDelta(fullTower, candidate, n)
		if delta == 0 {
			count++
		}
		sum += delta
	}
	fmt.Println("Part 1:", count)
	fmt.Println("Part 2:", sum)
}

func towerDelta(a, b []Line, except int) int {
	count := 0
	for n, i := range a {
		if n == except {
			continue
		}
		if i != b[n] {
			count++
		}
	}
	return count
}

func getTower(minx, miny, xwidth int, heightMap []int, lines []Line, except int, next []Line) {
	for n, i := range lines {
		if n == except {
			continue
		}
		highest := 0
		for y := i.a.y; y <= i.b.y; y++ {
			for x := i.a.x; x <= i.b.x; x++ {
				key := posKey(x, y, minx, miny, xwidth)
				h := heightMap[key]
				if h > highest {
					highest = h
				}
			}
		}
		i.a.z = highest + 1
		i.b.z = i.a.z + i.height
		next[n] = Line{
			a: i.a,
			b: i.b,
		}
		for y := i.a.y; y <= i.b.y; y++ {
			for x := i.a.x; x <= i.b.x; x++ {
				key := posKey(x, y, minx, miny, xwidth)
				heightMap[key] = i.b.z
			}
		}
	}
	clearHeightMap(heightMap)
}

func posKey(x, y, minx, miny, xwidth int) int {
	return (y-miny)*xwidth + x - minx
}

func clearHeightMap(a []int) {
	for i := range a {
		a[i] = 0
	}
}

func posLess(a, b Pos) int {
	if a.z == b.z {
		if a.y == b.y {
			return a.x - b.x
		}
		return a.y - b.y
	}
	return a.z - b.z
}

type (
	Line struct {
		a, b   Pos
		height int
	}

	Pos struct {
		x, y, z int
	}
)
