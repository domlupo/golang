// https://leetcode.com/problems/ransom-note/

package main

func canConstruct(ransomNote string, magazine string) bool {
	letters := make(map[rune]int)

	for _, l := range magazine {
		letters[l] = letters[l] + 1
	}

	for _, l := range ransomNote {
		c, _ := letters[l]
		if c == 0 {
			return false
		}
		letters[l] = letters[l] - 1
	}

	return true
}
