// https://leetcode.com/problems/rotate-image/
// R is the number of rows, C is the number of cols
// time complexity is O(R*C)
// space complexity is O(1)

package main

func rotate(matrix [][]int) {
	// transpose
	for r := range matrix {
		for c := r + 1; c < len(matrix); c++ {
			matrix[r][c], matrix[c][r] = matrix[c][r], matrix[r][c]
		}
	}

	// reverse rows
	for r := range matrix {
		for lc := 0; lc < len(matrix)/2; lc++ {
			rc := len(matrix) - 1 - lc
			matrix[r][lc], matrix[r][rc] = matrix[r][rc], matrix[r][lc]
		}
	}
}
