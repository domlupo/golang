// https://leetcode.com/problems/word-search/
// M is number of board rows, N is number of board cols, L is word length
// time complexity is O(M * N * 4^L)
// space complexity is O(L) due to call stack if board can be modified.
// if board cannot be modified it is O(M*N) for using a visited slice.

package main

func exist(board [][]byte, word string) bool {
	for r := range board {
		for c := range board[0] {
			if explore(board, r, c, word) {
				return true
			}
		}
	}
	return false
}

func explore(board [][]byte, r, c int, word string) bool {
	if len(word) == 0 {
		return true
	} else if r < 0 || r == len(board) || c < 0 || c == len(board[0]) {
		return false
	} else if board[r][c] != word[0] {
		return false
	}

	save := board[r][c]
	word = word[1:]
	board[r][c] = '0'
	left := explore(board, r-1, c, word)
	right := explore(board, r+1, c, word)
	up := explore(board, r, c-1, word)
	down := explore(board, r, c+1, word)
	board[r][c] = save
	return left || right || up || down
}

/*

// solution that does not modify board

func exist(board [][]byte, word string) bool {
	visited := make([][]bool, len(board))
	for r := range board {
		visited[r] = make([]bool, len(board[0]))
	}

	for r := range board {
		for c := range board[0] {
			if exists(board, r, c, word, visited) {
				return true
			}
		}
	}
	return false
}

func exists(board [][]byte, r, c int, word string, visited [][]bool) bool {
	if len(word) == 0 {
		return true
	} else if r < 0 || r == len(board) || c < 0 || c == len(board[0]) {
		return false
	} else if board[r][c] != word[0] {
		return false
	} else if visited[r][c] {
		return false
	}

	visited[r][c] = true
	word = word[1:]
	left := exists(board, r-1, c, word, visited)
	right := exists(board, r+1, c, word, visited)
	up := exists(board, r, c-1, word, visited)
	down := exists(board, r, c+1, word, visited)
	visited[r][c] = false
	return left || right || up || down
}

*/
