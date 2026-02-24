// https://leetcode.com/problems/uncommon-words-from-two-sentences/

package main

func uncommonFromSentences(s1 string, s2 string) []string {
	counts := make(map[string]int)
	countWords(counts, s1)
	countWords(counts, s2)

	var uncommon []string
	for k, v := range counts {
		if v == 1 {
			uncommon = append(uncommon, k)
		}
	}
	return uncommon
}

func countWords(counts map[string]int, s string) {
	var chars []byte
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ' ' {
			w := string(chars)
			chars = []byte{}
			counts[w]++
		} else {
			chars = append(chars, s[i])
		}
	}
}
