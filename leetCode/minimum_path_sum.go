// https://leetcode.com/problems/minimum-path-sum/
// M is the numbers of rows, N is the number of cols.
// time complexity is O(M * N)
// space complexity is O(M * N)

package main

func minPathSum(grid [][]int) int {
	m := len(grid)
	n := len(grid[0])

	var mem [][]int
	for range m {
		arr := make([]int, n)
		mem = append(mem, arr)
	}

	for r := range m {
		var prior int
		if r != 0 {
			prior = mem[r-1][0]
		}
		mem[r][0] = prior + grid[r][0]
	}

	for c := range n {
		var prior int
		if c != 0 {
			prior = mem[0][c-1]
		}
		mem[0][c] = prior + grid[0][c]
	}

	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			mem[r][c] = grid[r][c] + min(mem[r-1][c], mem[r][c-1])
		}
	}

	return mem[m-1][n-1]
}
