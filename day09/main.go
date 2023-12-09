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

var digitRegex = regexp.MustCompile(`-?\d+`)

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
		numStrs := digitRegex.FindAllString(scanner.Text(), -1)
		nums := make([]int, 0, len(numStrs))
		for _, i := range numStrs {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatalln(err)
			}
			nums = append(nums, num)
		}
		sum += findNextSeq(nums)
		slices.Reverse(nums)
		sum2 += findNextSeq(nums)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", sum)
	fmt.Println("Part 2:", sum2)
}

func findNextSeq(nums []int) int {
	numsAllZero := true
	for _, i := range nums {
		if i != 0 {
			numsAllZero = false
			break
		}
	}
	if numsAllZero {
		return 0
	}
	next := make([]int, len(nums)-1)
	for i := 1; i < len(nums); i++ {
		next[i-1] = nums[i] - nums[i-1]
	}
	nextDiff := findNextSeq(next)
	return nums[len(nums)-1] + nextDiff
}
