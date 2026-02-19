// https://leetcode.com/problems/lemonade-change/

package main

func lemonadeChange(bills []int) bool {
	var fives, tens int

	for _, b := range bills {
		if b == 20 && tens > 0 {
			tens--
			fives--
		} else if b == 20 {
			fives -= 3
		} else if b == 10 {
			tens++
			fives--
		} else if b == 5 {
			fives++
		}

		if fives < 0 || tens < 0 {
			return false
		}
	}

	return true
}
