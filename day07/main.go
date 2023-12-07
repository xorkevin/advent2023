package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

type (
	CardHand struct {
		Kind   int8
		KindJ  int8
		Score  int
		ScoreJ int
		Bid    int
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

	var hands []CardHand

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		a, b, ok := strings.Cut(scanner.Text(), " ")
		if !ok || len(a) != 5 {
			log.Fatalln("Invalid line")
		}
		num, err := strconv.Atoi(b)
		if err != nil {
			log.Fatalln(err)
		}
		ab := []byte(a)
		hands = append(hands, CardHand{
			Kind:   handKind(ab, false),
			KindJ:  handKind(ab, true),
			Score:  scoreHand(ab, false),
			ScoreJ: scoreHand(ab, true),
			Bid:    num,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	slices.SortFunc(hands, func(a, b CardHand) int {
		if a.Kind == b.Kind {
			return a.Score - b.Score
		}
		return int(a.Kind - b.Kind)
	})

	sum := 0
	for n, i := range hands {
		sum += (n + 1) * i.Bid
	}
	fmt.Println("Part 1:", sum)

	slices.SortFunc(hands, func(a, b CardHand) int {
		if a.KindJ == b.KindJ {
			return a.ScoreJ - b.ScoreJ
		}
		return int(a.KindJ - b.KindJ)
	})

	sum = 0
	for n, i := range hands {
		sum += (n + 1) * i.Bid
	}
	fmt.Println("Part 2:", sum)
}

func scoreCard(b byte, withJoker bool) int {
	if b >= '2' && b <= '9' {
		return int(b - '2' + 1)
	}
	switch b {
	case 'T':
		return 9
	case 'J':
		if withJoker {
			return 0
		} else {
			return 10
		}
	case 'Q':
		return 11
	case 'K':
		return 12
	case 'A':
		return 13
	}
	return 0
}

func scoreHand(b []byte, withJoker bool) int {
	sum := 0
	for _, i := range b {
		sum = sum*14 + scoreCard(i, withJoker)
	}
	return sum
}

func handKind(b []byte, withJoker bool) int8 {
	cardCountsMap := map[byte]int8{}
	var jokers int8 = 0
	for _, i := range b {
		if withJoker && i == 'J' {
			jokers++
			continue
		}
		cardCountsMap[i] = cardCountsMap[i] + 1
	}
	if len(cardCountsMap) == 0 {
		return 6
	}
	cardCounts := make([]int8, 0, len(cardCountsMap))
	for _, v := range cardCountsMap {
		cardCounts = append(cardCounts, v)
	}
	slices.Sort(cardCounts)
	slices.Reverse(cardCounts)
	cardCounts[0] += jokers
	if cardCounts[0] == 5 {
		// 5 of kind
		return 6
	}
	if cardCounts[0] == 4 {
		// 4 of kind
		return 5
	}
	if cardCounts[0] == 3 {
		if cardCounts[1] == 2 {
			// full house
			return 4
		} else {
			// 3 of kind
			return 3
		}
	}
	if cardCounts[0] == 2 {
		if cardCounts[1] == 2 {
			// two pair
			return 2
		} else {
			// one pair
			return 1
		}
	}
	return 0
}
