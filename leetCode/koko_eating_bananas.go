// https://leetcode.com/problems/koko-eating-bananas/
// N is the number of piles, M is the maximum number within piles.
// time complexity is O(N * log(M))
// space complexity is O(1)

package main

import "math"

func minEatingSpeed(piles []int, h int) int {
	l := 1 // need to set to 1 for edge cases, don't set to piles[0]
	r := piles[0]
	for _, pile := range piles {
		l = min(l, pile)
		r = max(r, pile)
	}

	for l < r {
		m := (l + r) / 2
		var hours int
		for _, pile := range piles {
			hours += int(math.Ceil(float64(pile) / float64(m)))
		}
		if hours <= h {
			r = m
		} else {
			l = m + 1
		}
	}

	// can return l or r since for loop is l < r not l <= r
	return l
}
