// https://leetcode.com/problems/island-perimeter/

package main

func islandPerimeter(grid [][]int) int {
	var p int

	for i := range len(grid) {
		for j := range len(grid[i]) {
			if grid[i][j] == 0 {
				continue
			}

			if i-1 < 0 || grid[i-1][j] == 0 {
				p++
			}

			if j-1 < 0 || grid[i][j-1] == 0 {
				p++
			}

			if i+1 > len(grid)-1 || grid[i+1][j] == 0 {
				p++
			}

			if j+1 > len(grid[i])-1 || grid[i][j+1] == 0 {
				p++
			}

		}
	}

	return p
}
