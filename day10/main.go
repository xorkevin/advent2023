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

	gridBounds := Coord{
		x: len(grid[0]),
		y: len(grid),
	}
	simpleGrid := make([][]byte, gridBounds.y)
	for i := range simpleGrid {
		simpleGrid[i] = make([]byte, gridBounds.x)
	}
	bounds := make([]Coord, gridBounds.y)
	for i := range bounds {
		bounds[i] = Coord{
			x: gridBounds.x,
			y: -1,
		}
	}

	neighbors := getStartNeighbors(grid, start)
	if len(neighbors) != 2 {
		log.Fatalln("Invalid neighbors")
	}

	startChar, ok := getStartChar(neighbors[0].dir, neighbors[1].dir)
	if !ok {
		log.Fatalln("Invalid neighbors")
	}
	simpleGrid[start.y][start.x] = startChar
	bounds[start.y] = Coord{
		x: start.x,
		y: start.x,
	}
	for _, i := range neighbors {
		simpleGrid[i.coord.y][i.coord.x] = grid[i.coord.y][i.coord.x]
		b := bounds[i.coord.y]
		bounds[i.coord.y] = Coord{
			x: min(b.x, i.coord.x),
			y: max(b.y, i.coord.x),
		}
	}

	steps := 1
	for allCoordsNotSame(neighbors) {
		for i := range neighbors {
			next, ok := nextStep(grid, neighbors[i])
			if !ok {
				log.Fatalln("Invalid next step")
			}
			neighbors[i] = next
			simpleGrid[next.coord.y][next.coord.x] = grid[next.coord.y][next.coord.x]
			b := bounds[next.coord.y]
			bounds[next.coord.y] = Coord{
				x: min(b.x, next.coord.x),
				y: max(b.y, next.coord.x),
			}
		}
		steps++
	}
	fmt.Println("Part 1:", steps)

	sum := 0
	for n, i := range simpleGrid {
		b := bounds[n]
		if b.y < 0 {
			continue
		}
		sum += getInsideCount(i[b.x : b.y+1])
	}
	fmt.Println("Part 2:", sum)
}

func getInsideCount(row []byte) int {
	count := 0
	inside := false
	var prev byte = 0
	for _, i := range row {
		switch i {
		case 0:
			if inside {
				count++
			}
		case '|':
			inside = !inside
			prev = 0
		case '-':
		case 'L':
			prev = 'L'
		case 'J':
			switch prev {
			case 'L':
			case 'F':
				inside = !inside
			}
			prev = 0
		case '7':
			switch prev {
			case 'L':
				inside = !inside
			case 'F':
			}
			prev = 0
		case 'F':
			prev = 'F'
		}
	}
	return count
}

func allCoordsNotSame(neighbors []CoordDir) bool {
	coord := neighbors[0].coord
	for _, i := range neighbors {
		if i.coord != coord {
			return true
		}
	}
	return false
}

func nextStep(grid [][]byte, posdir CoordDir) (CoordDir, bool) {
	transform, ok := tileDirMap[grid[posdir.coord.y][posdir.coord.x]]
	if !ok {
		return CoordDir{}, false
	}
	nextDir := transform[posdir.dir]
	switch nextDir {
	case DirNorth:
		return CoordDir{
			coord: Coord{
				x: posdir.coord.x,
				y: posdir.coord.y - 1,
			},
			dir: nextDir,
		}, true
	case DirEast:
		return CoordDir{
			coord: Coord{
				x: posdir.coord.x + 1,
				y: posdir.coord.y,
			},
			dir: nextDir,
		}, true
	case DirSouth:
		return CoordDir{
			coord: Coord{
				x: posdir.coord.x,
				y: posdir.coord.y + 1,
			},
			dir: nextDir,
		}, true
	case DirWest:
		return CoordDir{
			coord: Coord{
				x: posdir.coord.x - 1,
				y: posdir.coord.y,
			},
			dir: nextDir,
		}, true
	default:
		return CoordDir{}, false
	}
}

type (
	Coord struct {
		x, y int
	}

	CoordDir struct {
		coord Coord
		dir   Dir
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

func getStartChar(a, b Dir) (byte, bool) {
	if a > b {
		a, b = b, a
	}
	switch a {
	case DirNorth:
		switch b {
		case DirEast:
			return 'L', true
		case DirSouth:
			return '|', true
		case DirWest:
			return 'J', true
		}
	case DirEast:
		switch b {
		case DirSouth:
			return 'F', true
		case DirWest:
			return '-', true
		}
	case DirSouth:
		switch b {
		case DirWest:
			return '7', true
		}
	case DirWest:
	}
	return 0, false
}

func getStartNeighbors(grid [][]byte, pos Coord) []CoordDir {
	h := len(grid)
	w := len(grid[0])
	neighbors := make([]CoordDir, 0, 4)
	if pos.x > 0 {
		// has left neighbor
		switch grid[pos.y][pos.x-1] {
		case '-', 'L', 'F':
			neighbors = append(neighbors, CoordDir{
				coord: Coord{
					x: pos.x - 1,
					y: pos.y,
				},
				dir: DirWest,
			})
		}
	}
	if pos.y > 0 {
		// has top neighbor
		switch grid[pos.y-1][pos.x] {
		case '|', '7', 'F':
			neighbors = append(neighbors, CoordDir{
				coord: Coord{
					x: pos.x,
					y: pos.y - 1,
				},
				dir: DirNorth,
			})
		}
	}
	if pos.x+1 < w {
		// has right neighbor
		switch grid[pos.y][pos.x+1] {
		case '-', 'J', '7':
			neighbors = append(neighbors, CoordDir{
				coord: Coord{
					x: pos.x + 1,
					y: pos.y,
				},
				dir: DirEast,
			})
		}
	}
	if pos.y+1 < h {
		// has bot neighbor
		switch grid[pos.y+1][pos.x] {
		case '|', 'L', 'J':
			neighbors = append(neighbors, CoordDir{
				coord: Coord{
					x: pos.x,
					y: pos.y + 1,
				},
				dir: DirSouth,
			})
		}
	}
	return neighbors
}
