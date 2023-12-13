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

	sum := 0
	sum2 := 0

	var grid [][]byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			grid = append(grid, []byte(scanner.Text()))
			continue
		}

		s, s2, ok := findReflections(grid)
		if !ok {
			log.Fatalln("No mirror")
		}
		sum += s
		sum2 += s2
		grid = nil
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	s, s2, ok := findReflections(grid)
	if !ok {
		log.Fatalln("No mirror")
	}
	sum += s
	sum2 += s2

	fmt.Println("Part 1:", sum)
	fmt.Println("Part 2:", sum2)
}

func findReflections(grid [][]byte) (int, int, bool) {
	tgrid := transpose(grid)
	s, eh, ev, ok := findReflection(grid, tgrid)
	if !ok {
		return 0, 0, false
	}
	s2, ok := findSmudgeReflection(grid, tgrid, eh, ev)
	if !ok {
		return 0, 0, false
	}
	return s, s2, true
}

func findReflection(grid, transpose [][]byte) (int, int, int, bool) {
	if hm := findMirror(grid); hm > 0 {
		return hm * 100, hm, -1, true
	}
	if vm := findMirror(transpose); vm > 0 {
		return vm, -1, vm, true
	}
	return 0, -1, -1, false
}

func findSmudgeReflection(grid, transpose [][]byte, eh, ev int) (int, bool) {
	if hm := findSmudge(grid, eh); hm > 0 {
		return hm * 100, true
	}
	if vm := findSmudge(transpose, ev); vm > 0 {
		return vm, true
	}
	return 0, false
}

func findSmudge(grid [][]byte, except int) int {
	for i := 1; i < len(grid); i++ {
		if i == except {
			continue
		}
		if isAlmostMirroredAt(grid, i) {
			return i
		}
	}
	return -1
}

func isAlmostMirroredAt(grid [][]byte, r int) bool {
	height := len(grid)
	lim := min(height-r, r)
	count := 0
	for i := 0; i < lim; i++ {
		a := grid[r-i-1]
		b := grid[r+i]
		if string(a) != string(b) {
			if !isEditDistance1(a, b) {
				return false
			}
			count += 1
			if count > 1 {
				return false
			}
		}
	}
	return count == 1
}

func isEditDistance1(a, b []byte) bool {
	count := 0
	for i := range a {
		if a[i] != b[i] {
			count += 1
		}
		if count > 1 {
			return false
		}
	}
	return count == 1
}

func findMirror(grid [][]byte) int {
	for i := 1; i < len(grid); i++ {
		if isMirroredAt(grid, i) {
			return i
		}
	}
	return -1
}

func isMirroredAt(grid [][]byte, r int) bool {
	height := len(grid)
	lim := min(height-r, r)
	for i := 0; i < lim; i++ {
		if string(grid[r-i-1]) != string(grid[r+i]) {
			return false
		}
	}
	return true
}

func transpose(grid [][]byte) [][]byte {
	height := len(grid)
	width := len(grid[0])
	res := make([][]byte, width)
	for i := range res {
		res[i] = make([]byte, height)
		for j := range res[i] {
			res[i][j] = grid[j][i]
		}
	}
	return res
}
