// https://leetcode.com/problems/to-lower-case/

package main

func toLowerCase(s string) string {
	var str []rune
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			str = append(str, c+32)
		} else {
			str = append(str, c)
		}
	}
	return string(str)
}
