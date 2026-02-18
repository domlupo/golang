// https://leetcode.com/problems/rotate-string/

package main

func rotateString(s string, goal string) bool {
	if len(s) != len(goal) {
		return false
	}

	str := s + s
	for i := 0; i < len(str)-len(goal)+1; i++ {
		if str[i:len(goal)+i] == goal {
			return true
		}
	}

	return false
}
