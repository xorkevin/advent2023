package main

import (
	"fmt"
	"math"
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
	axis := 0.5 * float64(race.Time)
	if maxPossible := int(math.Floor(math.Pow(axis, 2))); maxPossible <= race.Dist {
		return 0
	}

	disc := math.Sqrt(math.Pow(float64(race.Time), 2)-(4.0*float64(race.Dist))) * 0.5
	start := int(math.Ceil(axis - disc))
	end := int(math.Floor(axis + disc))
	return end - start + 1
}
