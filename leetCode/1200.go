// https://leetcode.com/problems/minimum-absolute-difference/
// N is the length of nums
// time complexity is O(N)
// space complexity is O(N)

package main

import (
	"math"
	"slices"
)

func minimumAbsDifference(arr []int) [][]int {
	slices.Sort(arr)

	minAbsDiff := math.MaxInt
	for i := 1; i < len(arr); i++ {
		minAbsDiff = min(minAbsDiff, arr[i]-arr[i-1])
	}

	var ans [][]int
	for i := 1; i < len(arr); i++ {
		if arr[i]-arr[i-1] == minAbsDiff {
			ans = append(ans, []int{arr[i-1], arr[i]})
		}
	}
	return ans
}
