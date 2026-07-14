// https://leetcode.com/problems/edit-distance
// M is len(word1) and N is len(word2)
// time complexity is O(M * N)
// space complexity is O(M * N)

package main

func minDistance(word1 string, word2 string) int {
	var mem [][]int

	for i := range len(word1) + 1 {
		arr := make([]int, len(word2)+1)
		arr[0] = i
		mem = append(mem, arr)
	}

	for i := range len(word2) + 1 {
		mem[0][i] = i
	}

	for r := 1; r < len(word1)+1; r++ {
		for c := 1; c < len(word2)+1; c++ {
			if word1[r-1] == word2[c-1] {
				mem[r][c] = mem[r-1][c-1]
			} else {
				insert := mem[r][c-1]
				delete := mem[r-1][c]
				replace := mem[r-1][c-1]
				mem[r][c] = 1 + min72(insert, delete, replace)
			}
		}
	}

	return mem[len(word1)][len(word2)]
}

func min72(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= a && b <= c {
		return b
	}
	return c
}
