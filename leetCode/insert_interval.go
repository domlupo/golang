// https://leetcode.com/problems/insert-interval/
// N is the length of intervals
// time complexity is O(N)
// space complexity is O(1)

package main

func insert(intervals [][]int, newInterval []int) [][]int {
	var ans [][]int
	i := 0
	n := len(intervals)

	for i < n && newInterval[0] > intervals[i][1] {
		ans = append(ans, intervals[i])
		i++
	}

	for i < n && newInterval[1] >= intervals[i][0] {
		newInterval[0] = min(newInterval[0], intervals[i][0])
		newInterval[1] = max(newInterval[1], intervals[i][1])
		i++
	}
	ans = append(ans, newInterval)

	for i < n {
		ans = append(ans, intervals[i])
		i++
	}

	return ans
}
