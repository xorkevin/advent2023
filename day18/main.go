package main

import (
	"bufio"
	"fmt"
	"log"
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

	current := Pos{x: 0, y: 0}
	current2 := Pos{x: 0, y: 0}
	area := 0
	perimeter := 0
	area2 := 0
	perimeter2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if len(line) != 3 {
			log.Fatalln("Invalid line")
		}
		{
			dir := line[0]
			num, err := strconv.Atoi(line[1])
			if err != nil {
				log.Fatalln(err)
			}
			current = move(current, dir, num)
			switch dir {
			case "L":
				area -= current.y * num
			case "R":
				area += current.y * num
			}
			perimeter += num
		}
		{
			line2 := strings.Trim(line[2], "(#)")
			if len(line2) != 6 {
				log.Fatalln("Invalid line2")
			}
			dir := line2[5]
			num64, err := strconv.ParseInt(string(line2[:5]), 16, 64)
			if err != nil {
				log.Fatalln(err)
			}
			num := int(num64)
			current2 = move2(current2, dir, num)
			switch dir {
			case '2':
				area2 -= current2.y * num
			case '0':
				area2 += current2.y * num
			}
			perimeter2 += num
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	if perimeter%2 != 0 {
		log.Fatalln("Perimeter not aligned to grid")
	}
	area = abs(area)
	halfPerimeter := perimeter / 2
	fmt.Println("Part 1:", area+halfPerimeter+1)

	if perimeter2%2 != 0 {
		log.Fatalln("Perimeter not aligned to grid")
	}
	area2 = abs(area2)
	halfPerimeter2 := perimeter2 / 2
	fmt.Println("Part 2:", area2+halfPerimeter2+1)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func move2(p Pos, dir byte, num int) Pos {
	switch dir {
	case '3':
		p.y -= num
	case '1':
		p.y += num
	case '2':
		p.x -= num
	case '0':
		p.x += num
	default:
		log.Fatalln("Invalid dir")
	}
	return p
}

func move(p Pos, dir string, num int) Pos {
	switch dir {
	case "U":
		p.y -= num
	case "D":
		p.y += num
	case "L":
		p.x -= num
	case "R":
		p.x += num
	default:
		log.Fatalln("Invalid dir")
	}
	return p
}

type (
	Pos struct {
		x, y int
	}
)

func getBounds(p, tl, br Pos) (Pos, Pos) {
	return Pos{
			x: min(tl.x, p.x),
			y: min(tl.y, p.y),
		}, Pos{
			x: max(br.x, p.x),
			y: max(br.y, p.y),
		}
}
