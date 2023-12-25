package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
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

	var stones []Stone
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		posStr, velStr, ok := strings.Cut(scanner.Text(), " @ ")
		if !ok {
			log.Fatalln("Invalid line")
		}
		posNumStrs := strings.Split(posStr, ", ")
		if len(posNumStrs) != 3 {
			log.Fatalln("Invalid line")
		}
		var pos [3]int
		for n, i := range posNumStrs {
			var err error
			pos[n], err = strconv.Atoi(i)
			if err != nil {
				log.Fatalln(err)
			}
		}
		velNumStrs := strings.Split(velStr, ", ")
		if len(velNumStrs) != 3 {
			log.Fatalln("Invalid line")
		}
		var vel [3]int
		for n, i := range velNumStrs {
			var err error
			vel[n], err = strconv.Atoi(i)
			if err != nil {
				log.Fatalln(err)
			}
		}
		stones = append(stones, Stone{
			pos: pos,
			vel: vel,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	count := 0
	const boundA = 200000000000000.0
	const boundB = 400000000000000.0
	for n, i := range stones[:len(stones)-1] {
		for _, j := range stones[n+1:] {
			point, hasIntersection := findIntersection(i, j)
			if !hasIntersection || !inFuture(i, point) || !inFuture(j, point) {
				continue
			}
			if inBounds(point, boundA, boundB) {
				count++
			}
		}
	}
	fmt.Println("Part 1:", count)

	matrix := make([][]float64, 4)
	for i := range matrix {
		matrix[i] = matrixLine(stones[i], stones[i+1])
	}
	res, ok := gaussianElimination(matrix)
	if !ok {
		log.Fatalln("Unable to find exact solution")
	}
	x := int(math.Round(res[0]))
	y := int(math.Round(res[1]))
	vx := int(math.Round(res[2]))
	c0, t0, ok := matrixLine2(x, vx, stones[0])
	if !ok {
		log.Fatalln("Solution is non-integer")
	}
	c1, t1, ok := matrixLine2(x, vx, stones[1])
	if !ok {
		log.Fatalln("Solution is non-integer")
	}
	dc := c0 - c1
	dt := t0 - t1
	if dc%dt != 0 {
		log.Fatalln("Solution is non-integer")
	}
	vz := dc / dt
	z := (stones[0].vel[2]-vz)*t0 + stones[0].pos[2]
	fmt.Println("Part 2:", x+y+z)
}

func matrixLine2(x, vx int, a Stone) (int, int, bool) {
	dx := x - a.pos[0]
	dv := a.vel[0] - vx
	if dx%dv != 0 {
		return 0, 0, false
	}
	t := dx / dv
	return a.pos[2] + t*a.vel[2], t, true
}

func matrixLine(a, b Stone) []float64 {
	return []float64{
		float64(a.vel[1] - b.vel[1]),
		float64(b.vel[0] - a.vel[0]),
		float64(b.pos[1] - a.pos[1]),
		float64(a.pos[0] - b.pos[0]),
		float64(a.pos[0]*a.vel[1] - a.pos[1]*a.vel[0] - b.pos[0]*b.vel[1] + b.pos[1]*b.vel[0]),
	}
}

func inBounds(pos [2]float64, a, b float64) bool {
	return pos[0] >= a && pos[0] <= b && pos[1] >= a && pos[1] <= b
}

func inFuture(a Stone, b [2]float64) bool {
	dx := b[0] - float64(a.pos[0])
	dy := b[1] - float64(a.pos[1])
	if (dx < 0) != (a.vel[0] < 0) {
		return false
	}
	return (dx < 0) == (a.vel[0] < 0) && (dy < 0) == (a.vel[1] < 0)
}

func findIntersection(a, b Stone) ([2]float64, bool) {
	ma := float64(a.vel[1]) / float64(a.vel[0])
	mb := float64(b.vel[1]) / float64(b.vel[0])
	dm := mb - ma
	if math.Abs(dm) < 1e-6 {
		return [2]float64{}, false
	}
	x := float64(a.pos[1]-b.pos[1])/dm + (mb*float64(b.pos[0])-ma*float64(a.pos[0]))/dm
	y := (x-float64(a.pos[0]))*ma + float64(a.pos[1])
	return [2]float64{x, y}, true
}

type (
	Stone struct {
		pos [3]int
		vel [3]int
	}

	Vec3 struct {
		x, y, z int
	}
)

func gaussianElimination(matrix [][]float64) ([]float64, bool) {
	m := len(matrix)
	n := len(matrix[0])
	if n != m+1 {
		// invalid augmented matrix size
		return nil, false
	}
	if !forwardElim(matrix) {
		return nil, false
	}
	res := make([]float64, m)
	for i := m - 1; i >= 0; i-- {
		x := matrix[i][m]
		for j := i + 1; j < m; j++ {
			x -= matrix[i][j] * res[j]
		}
		res[i] = x / matrix[i][i]
	}
	return res, true
}

func forwardElim(matrix [][]float64) bool {
	m := len(matrix)
	n := len(matrix[0])
	for k := 0; k < m && k < n; k++ {
		imax := k
		amax := abs(matrix[imax][k])
		for i := k + 1; i < m; i++ {
			if v := abs(matrix[i][k]); v > amax {
				amax = v
				imax = i
			}
		}
		if amax == 0 {
			// singular matrix, so not guaranteed to be satisfiable or have unique
			// solutions
			return false
		}
		if imax != k {
			swapRows(matrix, k, n, k, imax)
		}
		first := matrix[k][k]
		for i := k + 1; i < m; i++ {
			f := matrix[i][k] / first
			matrix[i][k] = 0
			for j := k + 1; j < n; j++ {
				matrix[i][j] -= matrix[k][j] * f
			}
		}
	}
	return true
}

func swapRows(matrix [][]float64, k, n, a, b int) {
	for i := k; i < n; i++ {
		matrix[a][k], matrix[b][k] = matrix[b][k], matrix[a][k]
	}
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
