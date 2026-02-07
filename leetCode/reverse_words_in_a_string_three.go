// https://leetcode.com/problems/reverse-words-in-a-string-iii/

package main

func reverseWords(s string) string {
	b := []byte(s)
	b = append(b, ' ')

	var p1, p2 int
	for i := range b {
		if b[i] == ' ' {
			p2 = i - 1
			for p1 < p2 {
				b[p1], b[p2] = b[p2], b[p1]
				p1++
				p2--
			}
			p1 = i + 1
		}
	}

	return string(b[0 : len(b)-1])
}
