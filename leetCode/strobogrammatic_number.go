// https://leetcode.com/problems/strobogrammatic-number/

package main

func isStrobogrammatic(num string) bool {
	strobo := map[byte]byte{
		'1': '1',
		'8': '8',
		'0': '0',
		'6': '9',
		'9': '6',
	}

	for i := 0; i <= len(num)/2; i++ {
		s, ok := strobo[num[i]]
		if !ok {
			return false
		}
		if s != num[len(num)-1-i] {
			return false
		}
	}
	return true
}
