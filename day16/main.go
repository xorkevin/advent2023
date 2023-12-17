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

	fmt.Println("Part 1:", simulate(Beam{
		x:   0,
		y:   0,
		dir: DirEast,
	}, width, height, grid))

	fmt.Println("Part 2:", findLargest(width, height, grid))
}

func findLargest(width int, height int, grid [][]byte) int {
	greatest := 0
	for x := 0; x < width; x++ {
		m := simulate(Beam{
			x:   x,
			y:   0,
			dir: DirSouth,
		}, width, height, grid)
		greatest = max(greatest, m)
		m = simulate(Beam{
			x:   x,
			y:   height - 1,
			dir: DirNorth,
		}, width, height, grid)
		greatest = max(greatest, m)
	}
	for y := 0; y < height; y++ {
		m := simulate(Beam{
			x:   0,
			y:   y,
			dir: DirEast,
		}, width, height, grid)
		greatest = max(greatest, m)
		m = simulate(Beam{
			x:   width - 1,
			y:   y,
			dir: DirWest,
		}, width, height, grid)
		greatest = max(greatest, m)
	}
	return greatest
}

type (
	Dir int

	Beam struct {
		x, y int
		dir  Dir
	}

	Beams struct {
		beams []Beam
	}
)

func (b *Beams) Push(v Beam) {
	b.beams = append(b.beams, v)
}

func (b *Beams) Pop() (Beam, bool) {
	if len(b.beams) == 0 {
		return Beam{}, false
	}
	l := len(b.beams) - 1
	v := b.beams[l]
	b.beams = b.beams[:l]
	return v, true
}

const (
	DirNorth Dir = 0
	DirEast  Dir = 1
	DirSouth Dir = 2
	DirWest  Dir = 3
)

func simulate(start Beam, w, h int, grid [][]byte) int {
	hist := make([]bool, w*h)
	beamHist := make([]bool, w*h*4)

	sum := 0
	beams := &Beams{}
	beams.Push(start)
	for len(beams.beams) > 0 {
		beam, _ := beams.Pop()
		sum += stepBeam(w, h, beam, beams, grid, hist, beamHist)
	}

	return sum
}

func stepBeam(w, h int, beam Beam, beams *Beams, grid [][]byte, hist, beamHist []bool) int {
	hkey := beam.y*w + beam.x
	key := hkey*4 + int(beam.dir)
	if beamHist[key] {
		return 0
	}
	switch grid[beam.y][beam.x] {
	case '/':
		{
			next := beam
			switch next.dir {
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
			if isInBounds(next, w, h) {
				beams.Push(next)
			}
		}
	case '\\':
		{
			next := beam
			switch next.dir {
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
			if isInBounds(next, w, h) {
				beams.Push(next)
			}
		}
	case '|':
		{
			switch beam.dir {
			case DirNorth:
				{
					next := beam
					next.y--
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
				}
			case DirEast:
				{
					next := Beam{
						x:   beam.x,
						y:   beam.y - 1,
						dir: DirNorth,
					}
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
					next = Beam{
						x:   beam.x,
						y:   beam.y + 1,
						dir: DirSouth,
					}
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
				}
			case DirSouth:
				{
					next := beam
					next.y++
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
				}
			case DirWest:
				{
					next := Beam{
						x:   beam.x,
						y:   beam.y - 1,
						dir: DirNorth,
					}
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
					next = Beam{
						x:   beam.x,
						y:   beam.y + 1,
						dir: DirSouth,
					}
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
				}
			default:
				log.Fatalln("Invalid dir")
			}
		}
	case '-':
		{
			switch beam.dir {
			case DirNorth:
				{
					next := Beam{
						x:   beam.x - 1,
						y:   beam.y,
						dir: DirWest,
					}
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
					next = Beam{
						x:   beam.x + 1,
						y:   beam.y,
						dir: DirEast,
					}
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
				}
			case DirEast:
				{
					next := beam
					next.x++
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
				}
			case DirSouth:
				{
					next := Beam{
						x:   beam.x - 1,
						y:   beam.y,
						dir: DirWest,
					}
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
					next = Beam{
						x:   beam.x + 1,
						y:   beam.y,
						dir: DirEast,
					}
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
				}
			case DirWest:
				{
					next := beam
					next.x--
					if isInBounds(next, w, h) {
						beams.Push(next)
					}
				}
			default:
				log.Fatalln("Invalid dir")
			}
		}
	default:
		{
			next := beam
			switch next.dir {
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
			if isInBounds(next, w, h) {
				beams.Push(next)
			}
		}
	}
	beamHist[key] = true
	if hist[hkey] {
		return 0
	}
	hist[hkey] = true
	return 1
}

func isInBounds(beam Beam, w, h int) bool {
	return beam.x >= 0 && beam.y >= 0 && beam.x < w && beam.y < h
}
