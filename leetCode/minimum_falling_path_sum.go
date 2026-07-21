// https://leetcode.com/problems/minimum-falling-path-sum/
// R is the number of rows of the matrix, C is the number of cols
// time complexity is O(R * C)
// space complexity is O(C)

package main

import "math"

func minFallingPathSum(matrix [][]int) int {
	dp := make([]int, len(matrix[0]))

	for r := range len(matrix) {
		curr := make([]int, len(matrix[0]))
		for c := 0; c < len(matrix[r]); c++ {
			left, right := math.MaxInt, math.MaxInt
			if c-1 >= 0 {
				left = dp[c-1]
			}
			if c+1 < len(dp) {
				right = dp[c+1]
			}
			down := dp[c]

			minimum := min(left, right, down)
			curr[c] = matrix[r][c] + minimum
		}
		dp = curr
	}

	minimum := math.MaxInt
	for _, sum := range dp {
		minimum = min(minimum, sum)
	}

	return minimum
}
