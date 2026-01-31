// https://leetcode.com/problems/length-of-last-word/

package main

func lengthOfLastWord(s string) int {
	length := 0
	i := len(s) - 1
	for i != -1 && s[i] == ' ' {
		i--
	}
	for i != -1 && s[i] != ' ' {
		length++
		i--
	}
	return length
}
