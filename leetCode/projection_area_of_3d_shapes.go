// https://leetcode.com/problems/projection-area-of-3d-shapes/

package main

import "math"

func projectionArea(grid [][]int) int {
	var area int
	for i := range len(grid) {
		rowMax := math.MinInt
		colMax := math.MinInt
		for j := range len(grid[i]) {
			if grid[i][j] > 0 {
				area++
			}
			rowMax = max(rowMax, grid[i][j])
			colMax = max(colMax, grid[j][i])
		}
		area += rowMax
		area += colMax
	}
	return area
}
