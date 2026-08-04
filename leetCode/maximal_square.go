// https://leetcode.com/problems/maximal-square/
// R is the number of rows, C is the number of cols
// time complexity is O(R * C)
// space complexity is O(C)

package main

func maximalSquare(matrix [][]byte) int {
	dp := make([]int, len(matrix[0]))

	var ans int
	for r := range matrix {
		temp := make([]int, len(matrix[r]))
		for c := range matrix[r] {
			if matrix[r][c] == '0' {
				continue
			}

			if c == 0 {
				temp[c] = 1
			} else {
				temp[c] = 1 + min(temp[c-1], dp[c-1], dp[c])
			}
			ans = max(ans, temp[c])
		}
		dp = temp
	}

	return ans * ans
}
