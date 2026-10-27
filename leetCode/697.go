// https://leetcode.com/problems/degree-of-an-array/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(N)

package main

import "math"

func findShortestSubArray(nums []int) int {
	counts := make(map[int]int)
	firsts := make(map[int]int)
	lasts := make(map[int]int)
	var degree int
	for ni, n := range nums {
		counts[n]++
		if counts[n] == 1 {
			firsts[n] = ni
		}
		lasts[n] = ni
		degree = max(degree, counts[n])
	}

	ans := math.MaxInt
	for k, v := range counts {
		if v == degree {
			ans = min(ans, lasts[k]-firsts[k]+1)
		}
	}
	return ans
}
