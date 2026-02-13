// https://leetcode.com/problems/shortest-completing-word/

package main

func shortestCompletingWord(licensePlate string, words []string) string {
	charCounts := make(map[rune]int)

	for _, c := range licensePlate {
		if c >= 'A' && c <= 'Z' {
			charCounts[c+32] += 1
		} else if c >= 'a' && c <= 'z' {
			charCounts[c] += 1
		}
	}

	o := "1234567890123456"
	for _, w := range words {
		if len(w) < len(o) {
			counts := make(map[rune]int)
			for _, c := range w {
				counts[c] += 1
			}

			completing := true
			for k, v := range charCounts {
				if counts[k] < v {
					completing = false
				}
			}

			if completing {
				o = w
			}
		}
	}

	return o
}
