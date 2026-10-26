// https://leetcode.com/problems/add-strings/
// N is the max length of num1 and num2
// time complexity is O(N)
// space complexity is O(1)

package main

import "strconv"

func addStrings(num1 string, num2 string) string {
	var ans string
	var carry int
	p1, p2 := len(num1)-1, len(num2)-1
	for p1 >= 0 || p2 >= 0 || carry == 1 {
		var sum int
		if p1 >= 0 {
			sum += int(num1[p1] - '0')
			p1--
		}
		if p2 >= 0 {
			sum += int(num2[p2] - '0')
			p2--
		}
		sum += carry
		carry = sum / 10
		ans = strconv.Itoa(sum%10) + ans
	}
	return ans
}
