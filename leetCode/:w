// https://leetcode.com/problems/word-search/
// M is number of board rows, N is number of board cols, L is word length
// time complexity is O(M * N * 4^L)
// space complexity is O(L)

package main

func exist(board [][]byte, word string) bool {
	for r := range len(board) {
		for c := range len(board[0]) {
			if explore(0, r, c, word, board) {
				return true
			}
		}
	}

	return false
}

func explore(i, r, c int, word string, board [][]byte) bool {
	if i == len(word) {
		return true
	} else if !inbounds(r, c, board) {
		return false
	} else if board[r][c] != word[i] {
		return false
	}

	// if the board cannot be modified see below for map solution
	save := board[r][c]
	board[r][c] = '*'
	i++
	down := explore(i, r+1, c, word, board)
	up := explore(i, r-1, c, word, board)
	right := explore(i, r, c+1, word, board)
	left := explore(i, r, c-1, word, board)
	board[r][c] = save
	return down || up || right || left
}

func inbounds(r, c int, board [][]byte) bool {
	return r >= 0 && r < len(board) && c >= 0 && c < len(board[0])
}

// solution that uses a map

/*

type pos struct {
	r int
	c int
}

func exist(board [][]byte, word string) bool {
	for row := range len(board) {
		for col := range len(board[0]) {
			explored := make(map[pos]bool)
			p := pos{r: row, c: col}
			if explore(0, p, word, board, explored) {
				return true
			}
		}
	}

	return false
}

func explore(i int, p pos, word string, board [][]byte, explored map[pos]bool) bool {
	if i == len(word) {
		return true
	} else if !inbounds(p, board) {
		return false
	} else if board[p.r][p.c] != word[i] {
		return false
	} else if explored[p] {
		return false
	}

	explored[p] = true
	i++
	d := explore(i, pos{r: p.r + 1, c: p.c}, word, board, explored)
	u := explore(i, pos{r: p.r - 1, c: p.c}, word, board, explored)
	r := explore(i, pos{r: p.r, c: p.c + 1}, word, board, explored)
	l := explore(i, pos{r: p.r, c: p.c - 1}, word, board, explored)
	explored[p] = false
	return d || u || r || l
}

func inbounds(p pos, board [][]byte) bool {
	return p.r >= 0 && p.r < len(board) && p.c >= 0 && p.c < len(board[0])
}

*/
