// https://leetcode.com/problems/valid-sudoku/

package main

func isValidSudoku(board [][]byte) bool {
	var rows [9]int
	var cols [9]int
	var boxes [9]int

	for r := range 9 {
		for c := range 9 {
			if board[r][c] == '.' {
				continue
			}

			v := int(board[r][c])
			bit := 1 << (v - 1)

			// row check
			if rows[r]&bit != 0 {
				return false
			}
			rows[r] |= bit

			// col check
			if cols[c]&bit != 0 {
				return false
			}
			cols[c] |= bit

			// box check
			i := (r/3)*3 + c/3
			if boxes[i]&bit != 0 {
				return false
			}
			boxes[i] |= bit
		}
	}

	return true
}
