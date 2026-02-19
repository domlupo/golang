// https://leetcode.com/problems/backspace-string-compare

package main

func backspaceCompare(s string, t string) bool {
	si := len(s) - 1
	ti := len(t) - 1
	for si >= 0 || ti >= 0 {
		if processDeletes(s, &si) {
			continue
		}
		if processDeletes(t, &ti) {
			continue
		}
		if si >= 0 && ti >= 0 && s[si] == t[ti] {
			si--
			ti--
		} else {
			return false
		}
	}
	return true
}

func processDeletes(str string, ptr *int) bool {
	ptrStart := *ptr
	var d int
	for *ptr >= 0 {
		if str[*ptr] == '#' {
			d++
			*ptr--
			continue
		}
		if d > 0 {
			d--
			*ptr--
			continue
		}
		break
	}
	return ptrStart != *ptr
}
