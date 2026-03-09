// https://leetcode.com/problems/index-pairs-of-a-string/

package main

func indexPairs(text string, words []string) [][]int {
	// assuming that it is okay to modify passed words
	for i := 0; i < len(words)-1; i++ {
		if i >= 0 && len(words[i+1]) < len(words[i]) {
			words[i], words[i+1] = words[i+1], words[i]
			i -= 2
		}
	}

	var ans [][]int
	for i := 0; i < len(text); i++ {
		for _, w := range words {
			if i+len(w) <= len(text) && w == text[i:i+len(w)] {
				ans = append(ans, []int{i, i + len(w) - 1})
			}
		}
	}

	return ans
}
