package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
)

const (
	puzzleInput = "input.txt"
)

var (
	digitOnlyRegex = regexp.MustCompile(`\d`)
	digitRegex     = regexp.MustCompile(`\d|one|two|three|four|five|six|seven|eight|nine`)
	revDigitRegex  = regexp.MustCompile(`\d|enin|thgie|neves|xis|evif|ruof|eerht|owt|eno`)
	words          = map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
		"enin":  9,
		"thgie": 8,
		"neves": 7,
		"xis":   6,
		"evif":  5,
		"ruof":  4,
		"eerht": 3,
		"owt":   2,
		"eno":   1,
	}
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

	sum1 := 0
	sum2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Bytes()
		match := digitOnlyRegex.Find(s)
		if len(match) == 0 {
			log.Fatalln("not enough digits")
		}
		firstDigit, err := parseValue(string(match))
		if err != nil {
			log.Fatalln(err)
		}
		match = digitRegex.Find(s)
		if len(match) == 0 {
			log.Fatalln("not enough digits")
		}
		first, err := parseValue(string(match))
		if err != nil {
			log.Fatalln(err)
		}
		slices.Reverse(s)
		match = digitOnlyRegex.Find(s)
		if len(match) == 0 {
			log.Fatalln("not enough digits")
		}
		lastDigit, err := parseValue(string(match))
		if err != nil {
			log.Fatalln(err)
		}
		match = revDigitRegex.Find(s)
		if len(match) == 0 {
			log.Fatalln("not enough reverse digits")
		}
		last, err := parseValue(string(match))
		if err != nil {
			log.Fatalln(err)
		}
		sum1 += firstDigit*10 + lastDigit
		sum2 += first*10 + last
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", sum1)
	fmt.Println("Part 2:", sum2)
}

func parseValue(s string) (int, error) {
	if v, ok := words[s]; ok {
		return v, nil
	}
	return strconv.Atoi(s)
}
