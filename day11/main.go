package main

import (
	"bufio"
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

	var emptyRows []int
	var coords []Coord
	var cols [][]byte

	scanner := bufio.NewScanner(file)
	for y := 0; scanner.Scan(); y++ {
		isEmpty := true
		line := scanner.Bytes()
		if len(cols) == 0 {
			cols = make([][]byte, len(line))
		}
		for x, i := range line {
			if i == '#' {
				isEmpty = false
				coords = append(coords, Coord{
					x: x,
					y: y,
				})
				cols[x] = append(cols[x], i)
			}
		}
		if isEmpty {
			emptyRows = append(emptyRows, y)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	var emptyColumns []int
	for x, i := range cols {
		isEmpty := true
		for _, c := range i {
			if c == '#' {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			emptyColumns = append(emptyColumns, x)
		}
	}

	sd := 0
	se := 0
	for n, i := range coords {
		for _, j := range coords[n+1:] {
			sd += manhattanDistance(i, j)
			se += calcExpansion(emptyRows, i.y, j.y) + calcExpansion(emptyColumns, i.x, j.x)
		}
	}

	fmt.Println("Part 1:", sd+se)
	fmt.Println("Part 2:", sd+se*999999)
}

type (
	Coord struct {
		x, y int
	}
)

func calcExpansion(emptyRows []int, a, b int) int {
	if a > b {
		a, b = b, a
	}
	left := 0
	for left < len(emptyRows) && emptyRows[left] < a {
		left++
	}
	right := len(emptyRows) - 1
	for right >= 0 && emptyRows[right] > b {
		right--
	}
	return right - left + 1
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func manhattanDistance(a, b Coord) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}
