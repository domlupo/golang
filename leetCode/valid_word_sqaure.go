// https://leetcode.com/problems/valid-word-square/

package main

func validWordSquare(words []string) bool {
	for r := range words {
		for c := range words[r] {

			if c >= len(words) {
				return false
			}

			if r >= len(words[c]) {
				return false
			}

			if words[r][c] != words[c][r] {
				return false
			}
		}
	}

	return true
}
