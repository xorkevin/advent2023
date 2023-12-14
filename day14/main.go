package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	var grid [][]byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	dropRocks(grid)

	fmt.Println("Part 1:", scoreRocks(grid))

	height := len(grid)
	width := len(grid[0])
	other := make([][]byte, width)
	for i := range other {
		other[i] = make([]byte, height)
	}

	cache := map[string]int{}

	const p2Iterations = 1000000000

	period := 0

	for i := 0; i < p2Iterations; i++ {
		cycle(grid, other)
		s := getState(grid)
		if n, ok := cache[s]; ok {
			p := i - n
			if period == 0 {
				period = p
			} else {
				if p != period {
					log.Fatalln("Inconsistent fixed point period")
				}
				if (p2Iterations-i-1)%p == 0 {
					break
				}
			}
			cache[s] = i
		} else {
			cache[s] = i
		}
	}

	fmt.Println("Part 2:", scoreRocks(grid))
}

func getState(grid [][]byte) string {
	var s strings.Builder
	for _, i := range grid {
		s.Write(i)
	}
	return s.String()
}

func scoreRocks(grid [][]byte) int {
	height := len(grid)
	sum := 0
	for r, i := range grid {
		for _, j := range i {
			if j == 'O' {
				sum += height - r
			}
		}
	}
	return sum
}

func dropRocks(grid [][]byte) {
	for r, i := range grid {
		for c, j := range i {
			if j == 'O' {
				dropRock(grid, r, c)
			}
		}
	}
}

func dropRock(grid [][]byte, r, c int) {
	rest := r
	for rest-1 >= 0 {
		if grid[rest-1][c] != '.' {
			break
		}
		rest--
	}
	grid[r][c] = '.'
	grid[rest][c] = 'O'
}

func cycle(grid [][]byte, other [][]byte) {
	dropRocks(grid)
	// west
	rotate(grid, other)
	grid, other = other, grid
	dropRocks(grid)
	// south
	rotate(grid, other)
	grid, other = other, grid
	dropRocks(grid)
	// east
	rotate(grid, other)
	grid, other = other, grid
	dropRocks(grid)
	// north
	rotate(grid, other)
	grid, other = other, grid
}

func rotate(grid [][]byte, other [][]byte) {
	height := len(grid)
	for r, i := range grid {
		for c, j := range i {
			other[c][height-r-1] = j
		}
	}
}
