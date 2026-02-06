// https://leetcode.com/problems/keyboard-row/

package main

func findWords(words []string) []string {
	row := make(map[rune]int)

	topChars := "qwertyuiop"
	midChars := "asdfghjkl"
	botChars := "zxcvbnm"

	for _, c := range topChars {
		row[c] = 1
	}
	for _, c := range midChars {
		row[c] = 2
	}
	for _, c := range botChars {
		row[c] = 3
	}

	var o []string
	for _, word := range words {
		var p, r int
		for ; p < len(word); p++ {
			c := rune(word[p])
			if c >= 'A' && c <= 'Z' {
				c += 32
			}
			if r > 0 && row[c] != r {
				break
			}
			r = row[c]
		}
		if p == len(word) {
			o = append(o, word)
		}
	}
	return o
}
