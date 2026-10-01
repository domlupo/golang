// https://leetcode.com/problems/meeting-rooms/
// N is the length of intervals
// time complexity is O(N*log(N))
// space complexity is O(log(N))) Golang SortFunc uses pdqsort.

package main

import "slices"

func canAttendMeetings(intervals [][]int) bool {
	slices.SortFunc(intervals, func(a, b []int) int {
		if a[0] < b[0] {
			return -1
		} else {
			return 1
		}
	})

	for i := 1; i < len(intervals); i++ {
		l := intervals[i-1]
		r := intervals[i]
		if l[1] > r[0] {
			return false
		}
	}

	return true
}
