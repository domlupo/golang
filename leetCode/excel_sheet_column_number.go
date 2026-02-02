// https://leetcode.com/problems/excel-sheet-column-number/

package main

func titleToNumber(columnTitle string) int {
	var o int
	for _, r := range columnTitle {
		v := int(rune(r)-rune('A')) + 1
		o = o*26 + v
	}
	return o
}
