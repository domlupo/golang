// https://leetcode.com/problems/merge-intervals
// n is the length of intervals
// time complexity is O(n log(n))
// space complexity is O(1)

package main

import "slices"

func merge56(intervals [][]int) [][]int {
	slices.SortFunc(intervals, func(a, b []int) int {
		if a[0] < b[0] {
			return -1
		} else {
			return 1
		}
	})

	var ans [][]int
	m := intervals[0]
	for _, inter := range intervals {
		if m[1] < inter[0] {
			ans = append(ans, m)
			m = inter
		} else {
			m[1] = max(inter[1], m[1])
		}
	}

	return append(ans, m)
}
