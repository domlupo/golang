// https://leetcode.com/problems/meeting-rooms/

package main

func canAttendMeetings(intervals [][]int) bool {
	for i := 0; i < len(intervals)-1; i++ {
		if i >= 0 && intervals[i][0] > intervals[i+1][0] {
			intervals[i], intervals[i+1] = intervals[i+1], intervals[i]
			i -= 2
		}
	}

	for i := 0; i < len(intervals)-1; i++ {
		if intervals[i][1] > intervals[i+1][0] {
			return false
		}
	}

	return true
}
