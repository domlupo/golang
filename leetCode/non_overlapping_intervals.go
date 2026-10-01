// https://leetcode.com/problems/non-overlapping-intervals/
// N is the length of intervals
// time complexity is O(N*log(N))
// space complexity is O(log(N)) (Golang SortFunc uses pdqsort)

package main

import "slices"

func eraseOverlapIntervals(intervals [][]int) int {
	var ans int

	slices.SortFunc(intervals, func(a, b []int) int {
		if a[0] < b[0] {
			return -1
		} else {
			return 1
		}
	})

	curr := intervals[0]
	for i := 1; i < len(intervals); i++ {
		other := intervals[i]
		if curr[1] > other[0] && curr[1] > other[1] {
			curr = other
			ans++
		} else if curr[1] > other[0] {
			ans++
		} else {
			curr = other
		}
	}

	return ans
}
