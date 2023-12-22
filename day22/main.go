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
			name: len(lines),
			a:    lhsPos,
			b:    rhsPos,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	slices.SortFunc(lines, func(a, b Line) int {
		return posLess(a.a, b.a)
	})

	count := 0
	sum := 0
	fullTower := getTower(lines, -1)
	for n := range lines {
		t := getTower(lines, n)
		delta := towerDelta(fullTower, t, n)
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

func getTower(lines []Line, except int) []Line {
	next := make([]Line, 0, len(lines))
	heightMap := map[Vec2]int{}
	for n, i := range lines {
		if n == except {
			next = append(next, Line{})
			continue
		}
		highest := 0
		for y := i.a.y; y <= i.b.y; y++ {
			for x := i.a.x; x <= i.b.x; x++ {
				h, ok := heightMap[Vec2{x: x, y: y}]
				if !ok {
					continue
				}
				if h > highest {
					highest = h
				}
			}
		}
		height := i.b.z - i.a.z
		i.a.z = highest + 1
		i.b.z = i.a.z + height
		next = append(next, Line{
			a: i.a,
			b: i.b,
		})
		for y := i.a.y; y <= i.b.y; y++ {
			for x := i.a.x; x <= i.b.x; x++ {
				heightMap[Vec2{x: x, y: y}] = i.b.z
			}
		}
	}
	return next
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
		name int
		a, b Pos
	}

	Pos struct {
		x, y, z int
	}

	Vec2 struct {
		x, y int
	}

	Terrain struct {
		name int
		h    int
	}
)
