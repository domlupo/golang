// https://leetcode.com/problems/length-of-last-word/

package main

func lengthOfLastWord(s string) int {
	var length int
	i := len(s) - 1
	for s[i] == ' ' {
		i--
	}
	for i >= 0 && s[i] != ' ' {
		i--
		length++
	}
	return length
}
