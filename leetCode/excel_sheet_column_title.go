// https://leetcode.com/problems/excel-sheet-column-title/

package main

func convertToTitle(columnNumber int) string {
	var result string
	for columnNumber > 0 {
		columnNumber--
		charCode := 'A' + columnNumber%26
		result = string(rune(charCode)) + result
		columnNumber /= 26
	}
	return result
}
