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
		hands = append(hands, CardHand{
			Kind:  handKind([]byte(a)),
			KindJ: handKindJ([]byte(a)),
			Bid:   num,
			Str:   a,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	slices.SortFunc(hands, func(a, b CardHand) int {
		if a.Kind < b.Kind {
			return -1
		}
		if a.Kind > b.Kind {
			return 1
		}
		for i := range a.Str {
			x := cardToScore(a.Str[i])
			y := cardToScore(b.Str[i])
			if x < y {
				return -1
			}
			if x > y {
				return 1
			}
		}
		return 0
	})

	sum := 0
	for n, i := range hands {
		sum += (n + 1) * i.Bid
	}
	fmt.Println("Part 1:", sum)

	slices.SortFunc(hands, func(a, b CardHand) int {
		if a.KindJ < b.KindJ {
			return -1
		}
		if a.KindJ > b.KindJ {
			return 1
		}
		for i := range a.Str {
			x := cardToScoreJ(a.Str[i])
			y := cardToScoreJ(b.Str[i])
			if x < y {
				return -1
			}
			if x > y {
				return 1
			}
		}
		return 0
	})

	sum = 0
	for n, i := range hands {
		sum += (n + 1) * i.Bid
	}
	fmt.Println("Part 2:", sum)
}

func cardToScore(b byte) int {
	if b >= '2' && b <= '9' {
		return int(b - '2' + 2)
	}
	switch b {
	case 'T':
		return 10
	case 'J':
		return 11
	case 'Q':
		return 12
	case 'K':
		return 13
	case 'A':
		return 14
	}
	return 0
}

func cardToScoreJ(b byte) int {
	if b == 'J' {
		return 1
	}
	return cardToScore(b)
}

type (
	CardHand struct {
		Kind  int
		KindJ int
		Bid   int
		Str   string
	}

	CardCount struct {
		Card int
		Num  int
	}
)

func handKind(b []byte) int {
	cardCountsMap := map[int]int{}
	for _, i := range b {
		id := cardToScore(i)
		cardCountsMap[id] = cardCountsMap[id] + 1
	}
	cardCounts := make([]CardCount, 0, len(cardCountsMap))
	for k, v := range cardCountsMap {
		cardCounts = append(cardCounts, CardCount{
			Card: k,
			Num:  v,
		})
	}
	slices.SortFunc(cardCounts, func(a, b CardCount) int {
		if a.Num > b.Num {
			return -1
		}
		if a.Num < b.Num {
			return 1
		}
		if a.Card > b.Card {
			return -1
		}
		if a.Card < b.Card {
			return 1
		}
		return 0
	})
	return sortedCardCountKind(cardCounts)
}

func handKindJ(b []byte) int {
	jokers := 0
	cardCountsMap := map[int]int{}
	for _, i := range b {
		if i == 'J' {
			jokers++
			continue
		}
		id := cardToScore(i)
		cardCountsMap[id] = cardCountsMap[id] + 1
	}
	cardCounts := make([]CardCount, 0, len(cardCountsMap))
	for k, v := range cardCountsMap {
		cardCounts = append(cardCounts, CardCount{
			Card: k,
			Num:  v,
		})
	}
	slices.SortFunc(cardCounts, func(a, b CardCount) int {
		if a.Num > b.Num {
			return -1
		}
		if a.Num < b.Num {
			return 1
		}
		if a.Card > b.Card {
			return -1
		}
		if a.Card < b.Card {
			return 1
		}
		return 0
	})
	if len(cardCounts) == 0 {
		return 6
	}
	cardCounts[0] = CardCount{
		Card: cardCounts[0].Card,
		Num:  cardCounts[0].Num + jokers,
	}
	return sortedCardCountKind(cardCounts)
}

func sortedCardCountKind(counts []CardCount) int {
	if counts[0].Num == 5 {
		// 5 of kind
		return 6
	}
	if counts[0].Num == 4 {
		// 4 of kind
		return 5
	}
	if counts[0].Num == 3 && counts[1].Num == 2 {
		// full house
		return 4
	}
	if counts[0].Num == 3 && len(counts) == 3 {
		// 3 of kind
		return 3
	}
	if counts[0].Num == 2 && counts[1].Num == 2 {
		// two pair
		return 2
	}
	if counts[0].Num == 2 {
		// one pair
		return 1
	}
	return 0
}
