// https://leetcode.com/problems/surrounded-regions/
// R is the number of rows and C is the number of cols
// time complexity is O(R * C)
// space complexity is O(R * C)

package main

func solve(board [][]byte) {
	rows := len(board)
	cols := len(board[0])
	var seen [][]bool
	for range rows {
		seen = append(seen, make([]bool, cols))
	}

	for r := range rows {
		travel130(board, seen, r, 0)
		travel130(board, seen, r, cols-1)
	}
	for c := range cols {
		travel130(board, seen, 0, c)
		travel130(board, seen, rows-1, c)
	}

	for r := range rows {
		for c := range cols {
			switch board[r][c] {
			case 'O':
				board[r][c] = 'X'
			case 'F':
				board[r][c] = 'O'
			}
		}
	}
}

func travel130(board [][]byte, seen [][]bool, r, c int) {
	if seen[r][c] || board[r][c] == 'X' {
		return
	}

	seen[r][c] = true
	board[r][c] = 'F'

	if r-1 >= 0 {
		travel130(board, seen, r-1, c)
	}
	if r+1 < len(board) {
		travel130(board, seen, r+1, c)
	}
	if c-1 >= 0 {
		travel130(board, seen, r, c-1)
	}
	if c+1 < len(board[0]) {
		travel130(board, seen, r, c+1)
	}
}
