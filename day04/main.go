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

	sum := 0

	totalCards := 0
	bonusCards := make([]int, 256)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		a, b, ok := strings.Cut(line, ": ")
		if !ok {
			log.Fatalln("Invalid card")
		}
		cardNumStr := digitsRegex.FindString(a)
		cardNum, err := strconv.Atoi(cardNumStr)
		if err != nil {
			log.Fatalln(err)
		}
		a, b, ok = strings.Cut(b, " | ")
		winning := map[int]struct{}{}
		for _, i := range digitsRegex.FindAllString(a, -1) {
			winNum, err := strconv.Atoi(i)
			if err != nil {
				log.Fatalln(err)
			}
			winning[winNum] = struct{}{}
		}
		count := 0
		points := 0
		for _, i := range digitsRegex.FindAllString(b, -1) {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatalln(err)
			}
			if _, ok := winning[num]; ok {
				count++
				if points == 0 {
					points = 1
				} else {
					points *= 2
				}
			}
		}
		sum += points

		currentMultiplier := bonusCards[cardNum] + 1
		totalCards += currentMultiplier
		for i := 0; i < count; i++ {
			k := cardNum + 1 + i
			bonusCards[k] += currentMultiplier
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", sum)
	fmt.Println("Part 2:", totalCards)
}
