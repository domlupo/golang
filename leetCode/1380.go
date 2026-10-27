// https://leetcode.com/problems/lucky-numbers-in-a-matrix/
// R is the number of rows, C is the number of Cols
// time complexity is O(R * C)
// space complexity is O(1)

package main

import "math"

func luckyNumbers(matrix [][]int) []int {
	var ans []int
	for r := range matrix {
		minNum := math.MaxInt
		var minCol int
		for c := range matrix[0] {
			if matrix[r][c] < minNum {
				minNum = matrix[r][c]
				minCol = c
			}
		}

		var maxNum int
		for r := range matrix {
			maxNum = max(maxNum, matrix[r][minCol])
		}

		if minNum == maxNum {
			ans = append(ans, minNum)
		}
	}

	return ans
}
