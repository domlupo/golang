// https://leetcode.com/problems/merge-intervals
// n is the length of intervals
// time complexity is O(n log(n))
// space complexity is O(1)

package main

import "slices"

func merge56(intervals [][]int) [][]int {
	var ans [][]int
	slices.SortFunc(intervals, func(a, b []int) int {
		if a[0] < b[0] {
			return -1
		} else {
			return 1
		}
	})

	merged := intervals[0]
	for i := 1; i < len(intervals); i++ {
		if merged[1] >= intervals[i][0] {
			r := max(merged[1], intervals[i][1])
			merged = []int{merged[0], r}
		} else {
			ans = append(ans, merged)
			merged = intervals[i]
		}
	}

	return append(ans, merged)
}
