// https://leetcode.com/problems/convert-a-number-to-hexadecimal/

package main

import "math"

func toHex(num int) string {
	if num == 0 {
		return "0"
	}
	if num < 0 {
		num = num + int(math.Pow(2, 32))
	}

	intToHex := "0123456789abcdef"

	var b []byte
	for num > 0 {
		r := intToHex[num%16]
		b = append(b, r)
		num /= 16
	}

	var hex []byte
	for i := len(b) - 1; i >= 0; i-- {
		hex = append(hex, b[i])
	}

	return string(hex)
}
