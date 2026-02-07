// https://leetcode.com/problems/reverse-string-ii/

package main

func reverseStr(s string, k int) string {
	b := []byte(s)

	var p1, p2 int
	for i := range b {
		if i%(2*k) == 0 {
			p1 = i
			p2 = i + k - 1
		}

		if p2 >= len(b) {
			p2 = len(b) - 1
		}

		if p1 < p2 {
			b[p1], b[p2] = b[p2], b[p1]
		}
		p1++
		p2--
	}

	return string(b)
}
