package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const (
	puzzleInput = "input.txt"
)

var digitsRegex = regexp.MustCompile(`\d+`)

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

	sum := 0

	gears := map[int][]int{}

	lim := Vec2{
		x: len(grid[0]),
		y: len(grid),
	}

	buf := make([]Vec2, lim.x*2+2)
	for y, i := range grid {
		matches := digitsRegex.FindAllIndex(i, -1)
		for _, j := range matches {
			left := Vec2{
				x: j[0],
				y: y,
			}
			right := Vec2{
				x: j[1] - 1,
				y: y,
			}
			num, err := strconv.Atoi(string(i[j[0]:j[1]]))
			if err != nil {
				log.Fatalln(err)
			}
			isSym := false
			n := getNeighbors(left, right, lim, buf)
			for _, k := range buf[:n] {
				if sym := grid[k.y][k.x]; isSymbol(sym) {
					isSym = true
					if sym == '*' {
						id := k.y*lim.x + k.x
						gears[id] = append(gears[id], num)
					}
				}
			}
			if isSym {
				sum += num
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", sum)

	sum2 := 0
	for _, v := range gears {
		if len(v) == 2 {
			sum2 += v[0] * v[1]
		}
	}

	fmt.Println("Part 2:", sum2)
}

type (
	Vec2 struct {
		x, y int
	}
)

func getNeighbors(p1, p2, lim Vec2, buf []Vec2) int {
	if p1.x > p2.x {
		p1, p2 = p2, p1
	}
	if p1.y > p2.y {
		p1.y, p2.y = p2.y, p1.y
	}
	idx := 0
	if y := p1.y - 1; y >= 0 {
		for i := p1.x; i <= p2.x; i++ {
			buf[idx] = Vec2{
				x: i,
				y: y,
			}
			idx++
		}
	}
	if y := p2.y + 1; y < lim.y {
		for i := p1.x; i <= p2.x; i++ {
			buf[idx] = Vec2{
				x: i,
				y: y,
			}
			idx++
		}
	}
	if x := p1.x - 1; x >= 0 {
		for i := p1.y; i <= p2.y; i++ {
			buf[idx] = Vec2{
				x: x,
				y: i,
			}
			idx++
		}
		if y := p1.y - 1; y >= 0 {
			buf[idx] = Vec2{
				x: x,
				y: y,
			}
			idx++
		}
		if y := p2.y + 1; y < lim.y {
			buf[idx] = Vec2{
				x: x,
				y: y,
			}
			idx++
		}
	}
	if x := p2.x + 1; x < lim.x {
		for i := p1.y; i <= p2.y; i++ {
			buf[idx] = Vec2{
				x: x,
				y: i,
			}
			idx++
		}
		if y := p1.y - 1; y >= 0 {
			buf[idx] = Vec2{
				x: x,
				y: y,
			}
			idx++
		}
		if y := p2.y + 1; y < lim.y {
			buf[idx] = Vec2{
				x: x,
				y: y,
			}
			idx++
		}
	}
	return idx
}

func isSymbol(b byte) bool {
	return b != '.' && (b < '0' || b > '9')
}
