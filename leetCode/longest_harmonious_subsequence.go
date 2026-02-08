// https://leetcode.com/problems/longest-harmonious-subsequence/

package main

func findLHS(nums []int) int {
	counts := make(map[int]int)

	for _, n := range nums {
		counts[n] = counts[n] + 1
	}

	var m int
	for n, c := range counts {
		// can either do lower or upper, don't need both
		lowerCount, lowerOk := counts[n-1]
		if lowerOk {
			m = max(m, c+lowerCount)
		}
	}

	return m
}
