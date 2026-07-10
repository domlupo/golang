// https://leetcode.com/problems/set-matrix-zeroes
// M is the matrix rows and N is the matrix cols
// time complexity is O(M * N)
// space complexity is O(1)

package main

func setZeroes(matrix [][]int) {
	// use [0][0] as flag for first row
	// use firstCol as flag for firstCol
	var firstCol bool

	// set a row and col flag to indicate which row or col needs to become all zeros
	for r := range len(matrix) {
		// must check this before modifying first col
		if matrix[r][0] == 0 {
			firstCol = true
		}

		for c := range len(matrix[0]) {
			if matrix[r][c] == 0 {
				matrix[0][c] = 0
				matrix[r][0] = 0
			}
		}
	}

	// turn row and col to all zeros
	for r := 1; r < len(matrix); r++ {
		for c := 1; c < len(matrix[r]); c++ {
			if matrix[r][0] == 0 || matrix[0][c] == 0 {
				matrix[r][c] = 0
			}
		}
	}

	// set first row to all zeros
	if matrix[0][0] == 0 {
		for c := range len(matrix[0]) {
			matrix[0][c] = 0
		}
	}

	// set first col to all zeros
	if firstCol {
		for r := range len(matrix) {
			matrix[r][0] = 0
		}
	}
}
