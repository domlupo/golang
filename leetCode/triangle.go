// https://leetcode.com/problems/two-sum/
// N is number of rows in triangle
// time complexity is O(N * N)
// space complexity is O(N)

package main

import "math"

func minimumTotal(triangle [][]int) int {
	lastRow := len(triangle) - 1
	dp := make([]int, len(triangle[lastRow]))

	for r := range len(triangle) {
		curr := make([]int, len(triangle[lastRow]))
		for c := range len(triangle[r]) {
			var minimum int
			if c == 0 {
				minimum = dp[c]
			} else if c == len(triangle[r])-1 {
				minimum = dp[c-1]
			} else {
				minimum = min(dp[c-1], dp[c])
			}
			curr[c] = triangle[r][c] + minimum
		}
		dp = curr
	}

	minimum := math.MaxInt
	for _, total := range dp {
		minimum = min(minimum, total)
	}
	return minimum
}
