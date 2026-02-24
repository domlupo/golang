// https://leetcode.com/problems/surface-area-of-3d-shapes/

package main

func surfaceArea(grid [][]int) int {
	var area int

	for i := range len(grid) {
		for j := range len(grid[i]) {
			if grid[i][j] > 0 {
				area += 2
			}

			var l, r, u, d int
			if exists3(grid, i-1, j) {
				l = grid[i-1][j]
			}
			if exists3(grid, i+1, j) {
				r = grid[i+1][j]
			}
			if exists3(grid, i, j-1) {
				u = grid[i][j-1]
			}
			if exists3(grid, i, j+1) {
				d = grid[i][j+1]
			}
			if grid[i][j] > l {
				area += grid[i][j] - l
			}
			if grid[i][j] > r {
				area += grid[i][j] - r
			}
			if grid[i][j] > u {
				area += grid[i][j] - u
			}
			if grid[i][j] > d {
				area += grid[i][j] - d
			}
		}
	}

	return area
}

func exists3(grid [][]int, r, c int) bool {
	return r >= 0 && r < len(grid) && c >= 0 && c < len(grid[r])
}
