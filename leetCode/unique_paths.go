// https://leetcode.com/problems/unique-paths/
// m is number of rows and n is numbers of columns
// time complexity is O(m*n)
// space complexity is O(m*n)

package main

func uniquePaths(m int, n int) int {
	var mem [][]int
	for range m {
		arr := make([]int, n)
		arr[0] = 1
		mem = append(mem, arr)
	}

	for i := range n {
		mem[0][i] = 1
	}

	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			mem[r][c] = mem[r-1][c] + mem[r][c-1]
		}
	}

	return mem[m-1][n-1]
}
