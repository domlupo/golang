// https://leetcode.com/problems/search-a-2d-matrix-ii/
// R is the number of rows, C is the number of cols
// time complexity is O(R + C)
// space complexity is O(1)

package main

func searchMatrix240(matrix [][]int, target int) bool {
	rows := len(matrix)
	cols := len(matrix[0])
	r := rows - 1
	var c int

	for r >= 0 && c < cols {
		if matrix[r][c] == target {
			return true
		} else if matrix[r][c] > target {
			r--
		} else {
			c++
		}
	}

	return false
}
