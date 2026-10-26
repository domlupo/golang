// https://leetcode.com/problems/island-perimeter/
// R is the number of rows, C is the number of cols
// time complexity is O(R * C)
// space complexity is O(1)

package main

func islandPerimeter(grid [][]int) int {
	var ans int
	for r := range grid {
		for c := range grid[0] {
			if isWater(grid, r, c) {
				continue
			}
			ans += perimeter(grid, r-1, c)
			ans += perimeter(grid, r+1, c)
			ans += perimeter(grid, r, c-1)
			ans += perimeter(grid, r, c+1)
		}
	}
	return ans
}

func isWater(grid [][]int, r, c int) bool {
	return grid[r][c] == 0
}

func perimeter(grid [][]int, r, c int) int {
	inBounds := r >= 0 && r < len(grid) && c >= 0 && c < len(grid[0])
	if !inBounds || (inBounds && isWater(grid, r, c)) {
		return 1
	} else {
		return 0
	}
}
