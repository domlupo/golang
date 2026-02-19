// https://leetcode.com/problems/transpose-matrix/

package main

func transpose(matrix [][]int) [][]int {
	o := [][]int{}

	for range len(matrix[0]) {
		o = append(o, make([]int, len(matrix)))
	}

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			o[j][i] = matrix[i][j]
		}
	}

	return o
}
