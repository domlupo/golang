// https://leetcode.com/problems/valid-sudoku/
// N is the number of rows and cols
// time complexity is O(N*N)
// space complexity is O(N)

package main

func isValidSudoku(board [][]byte) bool {
	var rows, cols, squares [9]int

	for r := range board {
		for c := range board[0] {
			if board[r][c] == '.' {
				continue
			}

			b := 1 << int(board[r][c]-1)
			s := ((r / 3) * 3) + c/3
			if rows[r]&b != 0 || cols[c]&b != 0 || squares[s]&b != 0 {
				return false
			}

			rows[r] |= b
			cols[c] |= b
			squares[s] |= b
		}
	}

	return true
}
