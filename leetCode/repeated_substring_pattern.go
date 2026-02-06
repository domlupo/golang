// https://leetcode.com/problems/repeated-substring-pattern/

package main

func repeatedSubstringPattern(s string) bool {
	str := s[1:] + s[:len(s)-1]

	for i := 0; i < len(str)-len(s)+1; i++ {
		if s == str[i:len(s)+i] {
			return true
		}
	}

	return false
}
