// https://leetcode.com/problems/integer-to-roman/

package main

import "strings"

func intToRoman(num int) string {
	vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	valSyms := map[int]string{
		1:    "I",
		4:    "IV",
		5:    "V",
		9:    "IX",
		10:   "X",
		40:   "XL",
		50:   "L",
		90:   "XC",
		100:  "C",
		400:  "CD",
		500:  "D",
		900:  "CM",
		1000: "M",
	}

	var ans []string
	for i := 0; i < len(vals) && num > 0; i++ {
		for vals[i] <= num {
			ans = append(ans, valSyms[vals[i]])
			num -= vals[i]
		}
	}

	return strings.Join(ans, "")
}
