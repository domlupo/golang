// https://leetcode.com/problems/fraction-to-recurring-decimal/
// N is the number of repeating digits
// time complexity is O(N)
// space complexity is O(N)

package main

import "strconv"

func fractionToDecimal(numerator int, denominator int) string {
	n := numerator
	d := denominator
	if n == 0 {
		return "0"
	}

	var start string
	if (n < 0) != (d < 0) {
		start = start + "-"
	}
	if n < 0 {
		n = n * -1
	}
	if d < 0 {
		d = d * -1
	}

	num := n / d
	start = start + strconv.Itoa(num)
	n -= num * d
	if n == 0 {
		return start
	}
	start = start + "."
	n *= 10

	var decimal string
	var pos int
	seen := make(map[int]int)
	for n != 0 {
		if p, ok := seen[n]; ok {
			return start + decimal[:p] + "(" + decimal[p:] + ")"
		}
		seen[n] = pos
		pos++
		num := n / d
		decimal = decimal + strconv.Itoa(num)
		n -= (num * d)
		n *= 10
	}

	return start + decimal
}
