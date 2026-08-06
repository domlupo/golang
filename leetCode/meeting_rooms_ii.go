// https://leetcode.com/problems/meeting-rooms-ii/
// N is the length of intervals
// time complexity is O(N*log(N))
// space complexity is O(N)

package main

import "slices"

func minMeetingRooms(intervals [][]int) int {
	starts := make([]int, len(intervals))
	ends := make([]int, len(intervals))
	for i, interval := range intervals {
		starts[i] = interval[0]
		ends[i] = interval[1]
	}

	slices.Sort(starts)
	slices.Sort(ends)

	var ans, sPtr, ePtr int
	for sPtr < len(starts) {
		start := starts[sPtr]
		end := ends[ePtr]
		if end <= start {
			ePtr++
		} else {
			sPtr++
		}
		ans = max(ans, sPtr-ePtr)
	}

	return ans
}
