// https://leetcode.com/problems/strobogrammatic-number/

package main

func isStrobogrammatic(num string) bool {
	strobo := map[byte]byte{
		'0': '0',
		'1': '1',
		'6': '9',
		'8': '8',
		'9': '6',
	}

	var l int
	r := len(num) - 1
	for l <= r {
		rotated, ok := strobo[num[l]]
		if !ok || rotated != num[r] {
			return false
		}
		l++
		r--
	}

	return true
}
