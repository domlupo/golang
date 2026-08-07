// https://leetcode.com/problems/valid-sudoku/
// N is the number of rows and cols
// time complexity is O(N*N)
// space complexity is O(N)

package main

func isValidSudoku(board [][]byte) bool {
	var rows, cols, boxes [9]int
	for r := range 9 {
		for c := range 9 {
			if board[r][c] == '.' {
				continue
			}

			v := board[r][c]
			bit := 1 << (v - 1)
			i := (r/3)*3 + c/3

			if rows[r]&bit != 0 || cols[c]&bit != 0 || boxes[i]&bit != 0 {
				return false
			}

			rows[r] |= bit
			cols[c] |= bit
			boxes[i] |= bit
		}
	}

	return true
}
