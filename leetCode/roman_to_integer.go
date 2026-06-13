// https://leetcode.com/problems/roman-to-integer/description/

package main

import "math"

func romanToInt(s string) int {
	symVals := map[rune]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}

	var ans int
	prev := math.MaxInt
	for _, c := range s {
		v := symVals[c]
		if prev < v {
			ans -= prev * 2
		}
		ans += v
		prev = v
	}

	return ans
}
