// https://leetcode.com/problems/number-of-islands/
// R is the number of rows, C is the number of cols
// time complexity is O(R * C)
// space complexity is O(R * C)
// lower space complexity if allowed to modified passed array

package main

func numIslands(grid [][]byte) int {
	rows := len(grid)
	cols := len(grid[0])

	var seen [][]bool
	for range rows {
		seen = append(seen, make([]bool, cols))
	}

	var ans int
	for r := range rows {
		for c := range cols {
			if travel200(r, c, grid, seen) {
				ans++
			}
		}
	}
	return ans
}

func travel200(r, c int, grid [][]byte, seen [][]bool) bool {
	if r < 0 || r >= len(grid) {
		return false
	} else if c < 0 || c >= len(grid[0]) {
		return false
	} else if seen[r][c] {
		return false
	}

	seen[r][c] = true

	if grid[r][c] == '0' {
		return false
	}

	travel200(r-1, c, grid, seen)
	travel200(r+1, c, grid, seen)
	travel200(r, c-1, grid, seen)
	travel200(r, c+1, grid, seen)
	return true
}
