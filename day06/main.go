package main

import (
	"fmt"
)

type (
	Race struct {
		Time int
		Dist int
	}
)

func main() {
	races := []Race{
		{
			Time: 47,
			Dist: 207,
		},
		{
			Time: 84,
			Dist: 1394,
		},
		{
			Time: 74,
			Dist: 1209,
		},
		{
			Time: 67,
			Dist: 1014,
		},
	}

	n := 1
	for _, race := range races {
		n *= simulate(race)
	}
	fmt.Println("Part 1:", n)

	fmt.Println("Part 2:", simulate(Race{
		Time: 47847467,
		Dist: 207139412091014,
	}))
}

func simulate(race Race) int {
	count := 0
	for i := 0; i <= race.Time; i++ {
		speed := i
		duration := race.Time - i
		dist := speed * duration
		if dist > race.Dist {
			count++
		}
	}
	return count
}
