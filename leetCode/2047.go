// https://leetcode.com/problems/number-of-valid-words-in-a-sentence/
// N is the length of sentence
// time complexity is O(N)
// space complexity is O(1)

package main

func countValidWords(sentence string) int {
	var l, ans int
	for r := 0; r <= len(sentence); r++ {
		if r == len(sentence) || sentence[r] == ' ' {
			if validWord(sentence[l:r]) {
				ans++
			}
			l = r + 1
		}
	}
	return ans
}

func validWord(s string) bool {
	hyphen := -1
	for ci, c := range s {
		if c >= '0' && c <= '9' {
			return false
		} else if c == '-' && (ci == 0 || ci == len(s)-1) {
			return false
		} else if c == '-' && !(lowercaseChar(s[ci-1]) && lowercaseChar(s[ci+1])) {
			return false
		} else if c == '-' && hyphen != -1 {
			return false
		} else if c == '-' {
			hyphen = ci
		} else if ci < len(s)-1 && (c == '!' || c == '.' || c == ',') {
			return false
		}
	}
	return len(s) > 0
}

func lowercaseChar(b byte) bool {
	return b >= 'a' && b <= 'z'
}
