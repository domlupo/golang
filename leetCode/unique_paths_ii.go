// https://leetcode.com/problems/unique-paths-ii/
// M is the number of rows, N is the numbers of columns
// time complexity is O(M * N)
// space complexity is O(M * N)

package main

func uniquePathsWithObstacles(obstacleGrid [][]int) int {
	m := len(obstacleGrid)
	n := len(obstacleGrid[0])

	var mem [][]int
	for range m {
		arr := make([]int, n)
		mem = append(mem, arr)
	}

	for i := range m {
		if obstacleGrid[i][0] == 1 {
			break
		}
		mem[i][0] = 1
	}

	for i := range n {
		if obstacleGrid[0][i] == 1 {
			break
		}
		mem[0][i] = 1
	}

	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			if obstacleGrid[r][c] == 0 {
				mem[r][c] = mem[r-1][c] + mem[r][c-1]
			}
		}
	}

	return mem[m-1][n-1]
}
