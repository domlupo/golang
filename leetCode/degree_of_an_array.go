// https://leetcode.com/problems/degree-of-an-array/

package main

import "math"

func findShortestSubArray(nums []int) int {
	first := make(map[int]int)
	last := make(map[int]int)
	counts := make(map[int]int)
	var m int

	for ni, n := range nums {
		if _, ok := first[n]; !ok {
			first[n] = ni
		}
		last[n] = ni

		counts[n] = counts[n] + 1
		m = max(m, counts[n])
	}

	smallest := math.MaxInt
	for n, c := range counts {
		if c == m {
			smallest = min(smallest, last[n]-first[n]+1)
		}
	}
	return smallest
}
