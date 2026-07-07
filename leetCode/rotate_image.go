// https://leetcode.com/problems/rotate-image/

package main

func rotate(matrix [][]int) {
	// transpose
	for r := 0; r < len(matrix); r++ {
		for c := r + 1; c < len(matrix[r]); c++ {
			matrix[r][c], matrix[c][r] = matrix[c][r], matrix[r][c]
		}
	}

	// reverse rows
	for row := 0; row < len(matrix); row++ {
		var l int
		r := len(matrix[row]) - 1
		for l < r {
			matrix[row][l], matrix[row][r] = matrix[row][r], matrix[row][l]
			l++
			r--
		}
	}
}
