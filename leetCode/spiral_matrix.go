// https://leetcode.com/problems/spiral-matrix/
// m and n are matrix dimensions
// time complexity is O(m * n)
// space complexity is O(1)

package main

func spiralOrder(matrix [][]int) []int {
	left, right := 0, len(matrix[0])-1
	top, bottom := 0, len(matrix)-1

	var ans []int
	var r, c, rDelta int
	cDelta := 1
	for len(matrix) < len(matrix)*len(matrix[0]) {
		if cDelta == 1 && c == right {
			rDelta = 1
			cDelta = 0
			top++
		} else if rDelta == 1 && r == bottom {
			rDelta = 0
			cDelta = -1
			right--
		} else if cDelta == -1 && c == left {
			rDelta = -1
			cDelta = 0
			bottom--
		} else if rDelta == -1 && r == top {
			rDelta = 0
			cDelta = 1
			left++
		}

		ans = append(ans, matrix[r][c])
		r += rDelta
		c += cDelta
	}

	return ans
}
