// https://leetcode.com/problems/insert-interval/
// N is the length of intervals
// time complexity is O(N)
// space complexity is O(1)

package main

func insert(intervals [][]int, newInterval []int) [][]int {
	var i int
	var ans [][]int
	merged := newInterval

	for i < len(intervals) && merged[0] > intervals[i][1] {
		ans = append(ans, intervals[i])
		i++
	}
	for i < len(intervals) && intervals[i][0] <= merged[1] {
		l := min(merged[0], intervals[i][0])
		r := max(merged[1], intervals[i][1])
		merged = []int{l, r}
		i++
	}
	ans = append(ans, merged)
	for i < len(intervals) {
		ans = append(ans, intervals[i])
		i++
	}

	return ans
}
