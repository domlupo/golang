// https://leetcode.com/problems/unique-morse-code-words/

package main

func uniqueMorseRepresentations(words []string) int {
	unique := make(map[string]bool)

	code := []string{".-", "-...", "-.-.", "-..", ".", "..-.", "--.", "....", "..", ".---",
		"-.-", ".-..", "--", "-.", "---", ".--.", "--.-", ".-.", "...", "-", "..-", "...-", ".--",
		"-..-", "-.--", "--.."}

	for _, word := range words {
		var s string
		for _, c := range word {
			s = s + code[c-'a']
		}
		unique[s] = true
	}

	return len(unique)
}
