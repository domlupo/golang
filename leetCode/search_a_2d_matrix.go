// https://leetcode.com/problems/search-a-2d-matrix/
// M is the matrix rows and N is the matrix cols
// time complexity is O(log(M * N))
// space complexity is O(1)

package main

func searchMatrix(matrix [][]int, target int) bool {
	size := len(matrix) * len(matrix[0])
	var left int
	right := size - 1

	for left <= right {
		m := (left + right) / 2
		r := m / len(matrix[0])
		c := m % len(matrix[0])
		if matrix[r][c] == target {
			return true
		} else if matrix[r][c] < target {
			left = m + 1
		} else {
			right = m - 1
		}
	}

	return false
}
