package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

var numRegex = regexp.MustCompile(`\d+`)

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

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		a, b, ok := strings.Cut(scanner.Text(), ": ")
		if !ok {
			log.Fatalln("Invalid line")
		}
		gameNum, err := strconv.Atoi(numRegex.FindString(a))
		if err != nil {
			log.Fatalln(err)
		}
		possible := true
		maxRed := 0
		maxGreen := 0
		maxBlue := 0
		rounds := strings.Split(b, "; ")
		for _, i := range rounds {
			cubes := strings.Split(i, ", ")
			for _, j := range cubes {
				a, b, ok := strings.Cut(j, " ")
				if !ok {
					log.Fatalln("Invalid cubes")
				}
				count, err := strconv.Atoi(a)
				if err != nil {
					log.Fatalln(err)
				}
				switch b {
				case "red":
					if count > 12 {
						possible = false
					}
					maxRed = max(maxRed, count)
				case "green":
					if count > 13 {
						possible = false
					}
					maxGreen = max(maxGreen, count)
				case "blue":
					if count > 14 {
						possible = false
					}
					maxBlue = max(maxBlue, count)
				default:
					log.Fatalln("Invalid color")
				}
			}
		}
		if possible {
			sum += gameNum
		}
		sum2 += maxRed * maxGreen * maxBlue
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", sum)
	fmt.Println("Part 2:", sum2)
}
