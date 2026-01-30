// https://leetcode.com/problems/roman-to-integer/description/

package main

import "math"

func romanToInt(s string) int {
	var romanValue int
	prior := math.MaxInt
	for i := 0; i < len(s); i++ {
		current := oneSymbolValue(rune(s[i]))
		if prior < current {
			romanValue -= (prior * 2)
		}
		romanValue += current
		prior = current
	}

	return romanValue
}

func oneSymbolValue(r rune) int {
	switch r {
	case 'I':
		return 1
	case 'V':
		return 5
	case 'X':
		return 10
	case 'L':
		return 50
	case 'C':
		return 100
	case 'D':
		return 500
	case 'M':
		return 1000
	default:
		panic("symbolValue: unknown rune")
	}
}
