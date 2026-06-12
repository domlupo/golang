// https://leetcode.com/problems/excel-sheet-column-title/

package main

func convertToTitle(columnNumber int) string {
  var ans string

  for columnNumber > 0 {
    columnNumber--
    c := 'A' + columnNumber%26
    ans = string(c) + ans
    columnNumber /= 26
  }

  return ans

