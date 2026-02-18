// https://leetcode.com/problems/toeplitz-matrix/

package main

func isToeplitzMatrix(matrix [][]int) bool {
	for i := range len(matrix) {
		for j := range len(matrix[i]) {
			if exists2(matrix, i-1, j-1) && matrix[i-1][j-1] != matrix[i][j] {
				return false
			}
		}
	}

	return true
}

func exists2(matrix [][]int, r, c int) bool {
	return r >= 0 && r < len(matrix) && c >= 0 && c < len(matrix[r])
}
