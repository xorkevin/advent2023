package main

import (
	"bufio"
	"crypto/sha256"
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
	other := make([][]byte, width)
	for i := range other {
		other[i] = make([]byte, height)
	}

	dropRocks(grid, other)

	fmt.Println("Part 1:", scoreRocks(grid))

	remaining := 0
	const p2Iterations = 1000000000
	cache := map[string]int{}
	for i := 0; i < p2Iterations; i++ {
		cycle(grid, other)
		s := getState(grid)
		if n, ok := cache[s]; ok {
			p := i - n
			remaining = (p2Iterations - i - 1) % p
			break
		}
		cache[s] = i
	}
	for i := 0; i < remaining; i++ {
		cycle(grid, other)
	}

	fmt.Println("Part 2:", scoreRocks(grid))
}

func getState(grid [][]byte) string {
	h := sha256.New()
	for _, i := range grid {
		h.Write(i)
	}
	return string(h.Sum(nil))
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

func dropRocks(grid [][]byte, other [][]byte) {
	height := len(grid)
	for r, i := range grid {
		for c, j := range i {
			if j == 'O' {
				dropRock(grid, other, r, c)
			} else {
				other[c][height-r-1] = j
			}
		}
	}
}

func dropRock(grid [][]byte, other [][]byte, r, c int) {
	height := len(grid)
	rest := r
	for rest >= 1 {
		if grid[rest-1][c] != '.' {
			break
		}
		rest--
	}
	grid[r][c] = '.'
	grid[rest][c] = 'O'
	other[c][height-r-1] = '.'
	other[c][height-rest-1] = 'O'
}

func cycle(grid [][]byte, other [][]byte) {
	// west
	dropRocks(grid, other)
	// south
	dropRocks(other, grid)
	// east
	dropRocks(grid, other)
	// north
	dropRocks(other, grid)
}
