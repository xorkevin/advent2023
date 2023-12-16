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

	var grid [][]byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	height := len(grid)
	width := len(grid[0])

	fmt.Println("Part 1:", simulate([]Beam{
		{
			x:   0,
			y:   0,
			dir: DirEast,
		},
	}, width, height, grid))

	fmt.Println("Part 2:", findLargest(width, height, grid))
}

func findLargest(width int, height int, grid [][]byte) int {
	greatest := 0
	for x := 0; x < width; x++ {
		m := simulate([]Beam{
			{
				x:   x,
				y:   0,
				dir: DirSouth,
			},
		}, width, height, grid)
		greatest = max(greatest, m)
		m = simulate([]Beam{
			{
				x:   x,
				y:   height - 1,
				dir: DirNorth,
			},
		}, width, height, grid)
		greatest = max(greatest, m)
	}
	for y := 0; y < height; y++ {
		m := simulate([]Beam{
			{
				x:   0,
				y:   y,
				dir: DirEast,
			},
		}, width, height, grid)
		greatest = max(greatest, m)
		m = simulate([]Beam{
			{
				x:   width - 1,
				y:   y,
				dir: DirWest,
			},
		}, width, height, grid)
		greatest = max(greatest, m)
	}
	return greatest
}

func simulate(beams []Beam, w, h int, grid [][]byte) int {
	hist := make([]bool, w*h)

	beamHist := make([]bool, w*h*4)

	sum := 0
	for _, i := range beams {
		key := i.y*w + i.x
		if ok := hist[key]; !ok {
			hist[key] = true
			sum++
		}
	}

	for len(beams) > 0 {
		var s int
		beams, s = stepBeams(w, h, beams, grid, hist, beamHist)
		sum += s
	}

	return sum
}

type (
	Dir int

	Beam struct {
		x, y int
		dir  Dir
	}
)

const (
	DirNorth Dir = 0
	DirEast  Dir = 1
	DirSouth Dir = 2
	DirWest  Dir = 3
)

func stepBeams(w, h int, beams []Beam, grid [][]byte, hist []bool, beamHist []bool) ([]Beam, int) {
	sum := 0
	res := make([]Beam, 0, len(beams))
	for _, i := range beams {
		r := stepBeam(i, grid)
		for _, j := range r {
			if isInBounds(j.x, j.y, w, h) {
				key := ((j.y*w)+j.x)*4 + int(j.dir)
				if ok := beamHist[key]; !ok {
					hkey := j.y*w + j.x
					if ok := hist[hkey]; !ok {
						hist[hkey] = true
						sum++
					}
					res = append(res, j)
					beamHist[key] = true
				}
			}
		}
	}
	return res, sum
}

func isInBounds(x, y int, w, h int) bool {
	return x >= 0 && y >= 0 && x < w && y < h
}

func stepBeam(beam Beam, grid [][]byte) []Beam {
	b := grid[beam.y][beam.x]
	if b == '/' {
		next := beam
		switch beam.dir {
		case DirNorth:
			next.x++
			next.dir = DirEast
		case DirEast:
			next.y--
			next.dir = DirNorth
		case DirSouth:
			next.x--
			next.dir = DirWest
		case DirWest:
			next.y++
			next.dir = DirSouth
		default:
			log.Fatalln("Invalid dir")
		}
		return []Beam{next}
	}
	if b == '\\' {
		next := beam
		switch beam.dir {
		case DirNorth:
			next.x--
			next.dir = DirWest
		case DirEast:
			next.y++
			next.dir = DirSouth
		case DirSouth:
			next.x++
			next.dir = DirEast
		case DirWest:
			next.y--
			next.dir = DirNorth
		default:
			log.Fatalln("Invalid dir")
		}
		return []Beam{next}
	}
	if b == '|' {
		next := beam
		switch beam.dir {
		case DirNorth:
			next.y--
		case DirEast:
			return []Beam{
				{
					x:   beam.x,
					y:   beam.y - 1,
					dir: DirNorth,
				},
				{
					x:   beam.x,
					y:   beam.y + 1,
					dir: DirSouth,
				},
			}
		case DirSouth:
			next.y++
		case DirWest:
			return []Beam{
				{
					x:   beam.x,
					y:   beam.y - 1,
					dir: DirNorth,
				},
				{
					x:   beam.x,
					y:   beam.y + 1,
					dir: DirSouth,
				},
			}
		default:
			log.Fatalln("Invalid dir")
		}
		return []Beam{next}
	}
	if b == '-' {
		next := beam
		switch beam.dir {
		case DirNorth:
			return []Beam{
				{
					x:   beam.x - 1,
					y:   beam.y,
					dir: DirWest,
				},
				{
					x:   beam.x + 1,
					y:   beam.y,
					dir: DirEast,
				},
			}
		case DirEast:
			next.x++
		case DirSouth:
			return []Beam{
				{
					x:   beam.x - 1,
					y:   beam.y,
					dir: DirWest,
				},
				{
					x:   beam.x + 1,
					y:   beam.y,
					dir: DirEast,
				},
			}
		case DirWest:
			next.x--
		default:
			log.Fatalln("Invalid dir")
		}
		return []Beam{next}
	}
	if b == '.' {
		next := beam
		switch beam.dir {
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
		return []Beam{next}
	}
	log.Fatalln("Invalid char")
	return nil
}
