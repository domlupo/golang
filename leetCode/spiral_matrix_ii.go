// https://leetcode.com/problems/spiral-matrix-ii/
// N is the value of n
// time complexity is O(N * N)
// space complexity is O(1)

package main

func generateMatrix(n int) [][]int {
	var ans [][]int
	for range n {
		arr := make([]int, n)
		ans = append(ans, arr)
	}

	left, right := 0, n-1
	top, bottom := 0, n-1

	var r, c int
	rDelta, cDelta := 0, 1
	count := 1
	for count <= n*n {
		if cDelta == 1 && c == right {
			cDelta = 0
			rDelta = 1
			top++
		} else if rDelta == 1 && r == bottom {
			cDelta = -1
			rDelta = 0
			right--
		} else if cDelta == -1 && c == left {
			cDelta = 0
			rDelta = -1
			bottom--
		} else if rDelta == -1 && r == top {
			cDelta = 1
			rDelta = 0
			left++
		}
		ans[r][c] = count
		count++
		r += rDelta
		c += cDelta
	}

	return ans
}
