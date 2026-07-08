// https://leetcode.com/problems/merge-intervals
// n is the length of intervals
// time complexity is O(n log(n))
// space complexity is O(1)

package main

func merge56(intervals [][]int) [][]int {
	for i := 0; i < len(intervals)-1; i++ {
		if i >= 0 && intervals[i][0] > intervals[i+1][0] {
			intervals[i], intervals[i+1] = intervals[i+1], intervals[i]
			i -= 2
		}
	}

	var ans [][]int
	start := intervals[0][0]
	end := intervals[0][1]
	for _, i := range intervals {
		if end >= i[0] {
			end = max(end, i[1])
		} else {
			ans = append(ans, []int{start, end})
			start = i[0]
			end = i[1]
		}
	}
	ans = append(ans, []int{start, end})

	return ans
}
