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

const numSlots = 10

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
	bonusCards := [numSlots]int{}

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

		slot := (cardNum + 4) % numSlots
		currentMultiplier := bonusCards[slot] + 1
		bonusCards[slot] = 0
		totalCards += currentMultiplier
		for i := 1; i <= count; i++ {
			k := (slot + i) % numSlots
			bonusCards[k] += currentMultiplier
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", sum)
	fmt.Println("Part 2:", totalCards)
}
