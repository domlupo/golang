// https://leetcode.com/problems/integer-to-roman/
// N is the length of nums
// time complexity is O(1) due to num upper range is 3999
// space complexity is O(1)

package main

type romanVal struct {
	str string
	val int
}

func intToRoman(num int) string {
	romanValsSorted := []romanVal{
		{str: "M", val: 1000},
		{str: "CM", val: 900},
		{str: "D", val: 500},
		{str: "CD", val: 400},
		{str: "C", val: 100},
		{str: "XC", val: 90},
		{str: "L", val: 50},
		{str: "XL", val: 40},
		{str: "X", val: 10},
		{str: "IX", val: 9},
		{str: "V", val: 5},
		{str: "IV", val: 4},
		{str: "I", val: 1},
	}
	var ans string
	rem := num
	for _, rv := range romanValsSorted {
		for rem >= rv.val {
			ans = ans + rv.str
			rem -= rv.val
		}
	}
	return ans
}
