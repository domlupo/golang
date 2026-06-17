// https://leetcode.com/problems/excel-sheet-column-number/

package main

func titleToNumber(columnTitle string) int {
	var o int
	for _, c := range columnTitle {
		o *= 26
		o += int(c-'A') + 1
	}
	return o
}
