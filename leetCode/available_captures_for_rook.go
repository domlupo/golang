// https://leetcode.com/problems/available-captures-for-rook/

package main

func numRookCaptures(board [][]byte) int {
	var row, col int
	for i := range board {
		for j := range board[i] {
			if board[i][j] == 'R' {
				row, col = i, j
				break
			}
		}
	}

	deltas := [][]int{[]int{1, 0}, []int{-1, 0}, []int{0, 1}, []int{0, -1}}
	var pawnsAttacked int
	for _, delta := range deltas {
		r, c := row, col
		for exists4(board, r, c) {
			if board[r][c] == 'p' {
				pawnsAttacked++
				break
			} else if board[r][c] == 'B' {
				break
			}
			r, c = r+delta[0], c+delta[1]
		}
	}

	return pawnsAttacked
}

func exists4(board [][]byte, r, c int) bool {
	return r >= 0 && r < len(board) && c >= 0 && c < len(board[r])
}
