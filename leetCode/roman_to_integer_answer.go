// https://leetcode.com/problems/roman-to-integer/description/

package main

import "math"

func romanToInt(s string) int {
	symbolVals := map[rune]int{'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000}

	var romanInt int
	prev := math.MaxInt
	for _, r := range s {
		current := symbolVals[r]
		if prev < current {
			romanInt -= (prev * 2)
		}
		romanInt += current
		prev = current
	}

	return romanInt
}
