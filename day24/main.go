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

	a := stones[0]
	b := stones[1]
	rest := stones[2:]
	ta, tb := 130621773037, 423178590960
	if candidate, ok := getCandidate(posAt(a, ta), posAt(b, tb), ta, tb); ok {
		fmt.Println(candidate)
		if willCollideAll(candidate, rest) {
			fmt.Println(ta, tb, 0)
			fmt.Println("Part 2:", candidate.pos[0]+candidate.pos[1]+candidate.pos[2])
			return
		}
	}
	if candidate, ok := getCandidate(posAt(b, ta), posAt(a, tb), ta, tb); ok {
		fmt.Println(candidate)
		if willCollideAll(candidate, rest) {
			fmt.Println(ta, tb, 1)
			fmt.Println("Part 2:", candidate.pos[0]+candidate.pos[1]+candidate.pos[2])
			return
		}
	}
}

func getCandidate(a, b [3]int, ta, tb int) (Stone, bool) {
	dt := tb - ta
	dp := [3]int{
		b[0] - a[0],
		b[1] - a[1],
		b[2] - a[2],
	}
	if dp[0]%dt != 0 {
		return Stone{}, false
	}
	if dp[1]%dt != 0 {
		return Stone{}, false
	}
	if dp[2]%dt != 0 {
		return Stone{}, false
	}
	vel := [3]int{
		dp[0] / dt,
		dp[1] / dt,
		dp[2] / dt,
	}
	pos := [3]int{
		a[0] - vel[0]*ta,
		a[1] - vel[1]*ta,
		a[2] - vel[2]*ta,
	}
	return Stone{
		pos: pos,
		vel: vel,
	}, true
}

func willCollideAll(candidate Stone, stones []Stone) bool {
	for n, i := range stones {
		if !willCollide(candidate, i) {
			fmt.Println("failed on n", n)
			return false
		}
	}
	return true
}

func willCollide(a, b Stone) bool {
	t := -1
	for n, i := range a.pos {
		dv := b.vel[n] - a.vel[n]
		dx := i - b.pos[n]
		if dv == 0 {
			if dx != 0 {
				return false
			}
			continue
		}
		if t < 0 {
			t = dx / dv
			if t < 0 {
				return false
			}
			continue
		}
		if i+a.vel[n]*t != b.pos[n]+b.vel[n]*t {
			return false
		}
	}
	return true
}

func posAt(stone Stone, t int) [3]int {
	return [3]int{
		stone.pos[0] + stone.vel[0]*t,
		stone.pos[1] + stone.vel[1]*t,
		stone.pos[2] + stone.vel[2]*t,
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
