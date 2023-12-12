package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
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

	sum := 0
	sum2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		a, b, ok := bytes.Cut(scanner.Bytes(), []byte{' '})
		if !ok {
			log.Fatalln("Invalid line")
		}
		var nums []int
		for _, i := range bytes.Split(b, []byte{','}) {
			num, err := strconv.Atoi(string(i))
			if err != nil {
				log.Fatalln(err)
			}
			nums = append(nums, num)
		}
		cache := make([]int, len(a)*len(nums))
		sum += getNumArrangements(a, nums, cache, len(a))
		aa := make([]byte, 0, len(a)*5+5)
		bb := make([]int, 0, len(nums)*5)
		for i := 0; i < 5; i++ {
			if len(aa) > 0 {
				aa = append(aa, '?')
			}
			aa = append(aa, a...)
			bb = append(bb, nums...)
		}
		cache = make([]int, len(aa)*len(bb))
		sum2 += getNumArrangements(aa, bb, cache, len(aa))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", sum)
	fmt.Println("Part 2:", sum2)
}

func getNumArrangements(b []byte, nums []int, cache []int, cacheRowWidth int) int {
	if len(nums) == 0 {
		if restNoGroup(b) {
			return 1
		} else {
			return 0
		}
	}
	if len(b) == 0 {
		return 0
	}

	key := cacheRowWidth*(len(nums)-1) + len(b) - 1
	if n := cache[key]; n > 0 {
		return n - 1
	}

	first := b[0]
	if first == '.' {
		count := getNumArrangements(b[1:], nums, cache, cacheRowWidth)
		cache[key] = count + 1
		return count
	}

	firstNum := nums[0]
	prefixMatches, isEnd := matchPrefix(b, firstNum)
	if isEnd {
		if len(nums) == 1 {
			cache[key] = 2
			return 1
		}
		cache[key] = 1
		return 0
	}

	count := 0
	if first == '?' {
		count += getNumArrangements(b[1:], nums, cache, cacheRowWidth)
	} else {
		if !prefixMatches {
			cache[key] = 1
			return 0
		}
	}
	if prefixMatches {
		count += getNumArrangements(b[firstNum+1:], nums[1:], cache, cacheRowWidth)
	}

	cache[key] = count + 1

	return count
}

func restNoGroup(b []byte) bool {
	for _, i := range b {
		if i == '#' {
			return false
		}
	}
	return true
}

func matchPrefix(b []byte, num int) (bool, bool) {
	if len(b) < num {
		return false, false
	}
	for i := 0; i < num; i++ {
		if b[i] == '.' {
			return false, false
		}
	}
	if len(b) == num {
		return true, true
	}
	if b[num] == '#' {
		return false, false
	}
	return true, false
}
