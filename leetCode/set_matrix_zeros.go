// https://leetcode.com/problems/set-matrix-zeroes
// M is the matrix rows and N is the matrix cols
// time complexity is O(M * N)
// space complexity is O(1)

package main

func setZeroes(matrix [][]int) {
	var setFirstCol, setFirstRow bool

	for r := range matrix {
		for c := range matrix[0] {
			if r == 0 && matrix[r][c] == 0 {
				setFirstRow = true
			}
			if c == 0 && matrix[r][c] == 0 {
				setFirstCol = true
			}
			if r > 0 && c > 0 && matrix[r][c] == 0 {
				matrix[0][c] = 0
				matrix[r][0] = 0
			}
		}
	}

	for r := 1; r < len(matrix); r++ {
		for c := 1; c < len(matrix[0]); c++ {
			if matrix[0][c] == 0 || matrix[r][0] == 0 {
				matrix[r][c] = 0
			}
		}
	}

	if setFirstCol {
		for r := range matrix {
			matrix[r][0] = 0
		}
	}

	if setFirstRow {
		for c := range matrix[0] {
			matrix[0][c] = 0
		}
	}
}
