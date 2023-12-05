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

var digitRegex = regexp.MustCompile(`\d+`)

type (
	Range struct {
		Start int
		End   int
	}

	Range2 struct {
		Dest Range
		Src  Range
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

	var seeds []int
	var seeds2 []Range

	var rangeMaps [][]Range2

	var lastRangeMap []Range2

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "seeds:") {
			for _, i := range digitRegex.FindAllString(line, -1) {
				num, err := strconv.Atoi(i)
				if err != nil {
					log.Fatalln(err)
				}
				seeds = append(seeds, num)
			}
			for i := 1; i < len(seeds); i += 2 {
				seeds2 = append(seeds2, Range{
					Start: seeds[i-1],
					End:   seeds[i-1] + seeds[i],
				})
			}
			continue
		} else if strings.HasSuffix(line, "map:") {
			if len(lastRangeMap) != 0 {
				rangeMaps = append(rangeMaps, lastRangeMap)
				lastRangeMap = nil
			}
			continue
		} else if line == "" {
			continue
		}
		nums := digitRegex.FindAllString(line, -1)
		if len(nums) != 3 {
			log.Fatalln("Invalid range", line)
		}
		num1, err := strconv.Atoi(nums[0])
		if err != nil {
			log.Fatalln(err)
		}
		num2, err := strconv.Atoi(nums[1])
		if err != nil {
			log.Fatalln(err)
		}
		num3, err := strconv.Atoi(nums[2])
		if err != nil {
			log.Fatalln(err)
		}
		lastRangeMap = append(lastRangeMap, Range2{
			Dest: Range{
				Start: num1,
				End:   num1 + num3,
			},
			Src: Range{
				Start: num2,
				End:   num2 + num3,
			},
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	if len(lastRangeMap) != 0 {
		rangeMaps = append(rangeMaps, lastRangeMap)
		lastRangeMap = nil
	}

	for _, i := range rangeMaps {
		seeds = runRange(seeds, i)
		seeds2 = runRange2(seeds2, i)
	}

	minSeed := seeds[0]
	for _, i := range seeds {
		if i < minSeed {
			minSeed = i
		}
	}
	fmt.Println("Part 1:", minSeed)

	minSeed2 := seeds2[0].Start
	for _, i := range seeds2 {
		if i.Start < minSeed2 {
			minSeed2 = i.Start
		}
	}
	fmt.Println("Part 2:", minSeed2)
}

func runRange(seeds []int, rangeMap []Range2) []int {
	res := make([]int, 0, len(seeds))
	for _, i := range seeds {
		k := i
		for _, j := range rangeMap {
			if i >= j.Src.Start && i < j.Src.End {
				k = i - j.Src.Start + j.Dest.Start
				break
			}
		}
		res = append(res, k)
	}
	return res
}

func runRange2(seeds []Range, rangeMap []Range2) []Range {
	res := make([]Range, 0, len(seeds))
	otherRanges := make([]Range, 0, len(seeds))
	for _, i := range seeds {
		nextStart := i.Start
		nextEnd := i.End
		for _, j := range rangeMap {
			if i.Start >= j.Src.Start && i.Start < j.Src.End {
				nextStart = i.Start - j.Src.Start + j.Dest.Start
				if i.End <= j.Src.End {
					nextEnd = i.End - j.Src.Start + j.Dest.Start
				} else {
					nextEnd = j.Dest.End
					otherRanges = append(otherRanges, Range{
						Start: j.Src.End,
						End:   i.End,
					})
				}
				break
			}
			if i.End > j.Src.Start && i.End <= j.Src.End {
				nextEnd = i.End - j.Src.Start + j.Dest.Start
				nextStart = j.Dest.Start
				otherRanges = append(otherRanges, Range{
					Start: i.Start,
					End:   j.Src.Start,
				})
				break
			}
			if i.Start < j.Src.Start && i.End > j.Src.End {
				nextStart = j.Dest.Start
				nextEnd = j.Dest.End
				otherRanges = append(otherRanges,
					Range{
						Start: i.Start,
						End:   j.Src.Start,
					},
					Range{
						Start: j.Src.End,
						End:   i.End,
					},
				)
				break
			}
		}
		res = append(res, Range{
			Start: nextStart,
			End:   nextEnd,
		})
	}
	if len(otherRanges) != 0 {
		res = append(res, runRange2(otherRanges, rangeMap)...)
	}
	return res
}
