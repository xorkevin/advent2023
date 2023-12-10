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
	start := Coord{x: -1, y: -1}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		k := []byte(scanner.Text())
		x := bytes.IndexByte(k, 'S')
		if x >= 0 {
			start = Coord{
				x: x,
				y: len(grid),
			}
		}
		grid = append(grid, k)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	curPos, curDir, ok := getStartNeighbor(grid, start)
	if !ok {
		log.Fatalln("Missing start neighbor")
	}
	steps := 1
	area := 0
	switch curDir {
	case DirNorth:
	case DirEast:
		area += curPos.y
	case DirSouth:
	case DirWest:
		area -= curPos.y
	}
	for curPos != start {
		transform, ok := tileDirMap[grid[curPos.y][curPos.x]]
		if !ok {
			log.Fatalln("Invalid pipe path")
		}
		curDir = transform[curDir]
		switch curDir {
		case DirNorth:
			curPos.y--
		case DirEast:
			curPos.x++
		case DirSouth:
			curPos.y++
		case DirWest:
			curPos.x--
		default:
			log.Fatalln("Invalid pipe connection")
		}
		steps++
		switch curDir {
		case DirNorth:
		case DirEast:
			area += curPos.y
		case DirSouth:
		case DirWest:
			area -= curPos.y
		}
	}
	if steps%2 != 0 {
		log.Fatalln("Pipe path not aligned to grid")
	}
	halfSteps := steps / 2
	fmt.Println("Part 1:", halfSteps)
	area = abs(area)
	fmt.Println("Part 2:", area-halfSteps+1)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type (
	Coord struct {
		x, y int
	}

	Dir int
)

const (
	DirNorth Dir = 0
	DirEast      = 1
	DirSouth     = 2
	DirWest      = 3
)

var tileDirMap = map[byte][4]Dir{
	'|': {DirNorth, -1, DirSouth, -1},
	'-': {-1, DirEast, -1, DirWest},
	'L': {-1, -1, DirEast, DirNorth},
	'J': {-1, DirNorth, DirWest, -1},
	'7': {DirWest, DirSouth, -1, -1},
	'F': {DirEast, -1, -1, DirSouth},
}

func getStartNeighbor(grid [][]byte, pos Coord) (Coord, Dir, bool) {
	h := len(grid)
	w := len(grid[0])
	if pos.x > 0 {
		// has left neighbor
		switch grid[pos.y][pos.x-1] {
		case '-', 'L', 'F':
			return Coord{
				x: pos.x - 1,
				y: pos.y,
			}, DirWest, true
		}
	}
	if pos.y > 0 {
		// has top neighbor
		switch grid[pos.y-1][pos.x] {
		case '|', '7', 'F':
			return Coord{
				x: pos.x,
				y: pos.y - 1,
			}, DirNorth, true
		}
	}
	if pos.x+1 < w {
		// has right neighbor
		switch grid[pos.y][pos.x+1] {
		case '-', 'J', '7':
			return Coord{
				x: pos.x + 1,
				y: pos.y,
			}, DirEast, true
		}
	}
	if pos.y+1 < h {
		// has bot neighbor
		switch grid[pos.y+1][pos.x] {
		case '|', 'L', 'J':
			return Coord{
				x: pos.x,
				y: pos.y + 1,
			}, DirSouth, true
		}
	}
	return Coord{}, DirNorth, false
}
