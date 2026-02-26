// https://leetcode.com/problems/matrix-cells-in-distance-order/

package main

func allCellsDistOrder(rows int, cols int, rCenter int, cCenter int) [][]int {
	ans := [][]int{[]int{rCenter, cCenter}}
	visited := make([][]bool, rows)
	for i := range len(visited) {
		visited[i] = make([]bool, cols)
	}

	deltas := [][]int{[]int{1, 0}, []int{-1, 0}, []int{0, 1}, []int{0, -1}}
	s := [][]int{[]int{rCenter, cCenter}}
	visited[rCenter][cCenter] = true
	for len(s) > 0 {
		curr := s[0]
		s = s[1:]
		for _, delta := range deltas {
			r := curr[0] + delta[0]
			c := curr[1] + delta[1]
			if exists5(rows, cols, r, c) && !visited[r][c] {
				ans = append(ans, []int{r, c})
				visited[r][c] = true
				s = append(s, []int{r, c})
			}
		}
	}

	return ans
}

func exists5(rMax, cMax, r, c int) bool {
	return r >= 0 && r < rMax && c >= 0 && c < cMax
}
