// https://leetcode.com/problems/confusing-number/

package main

func confusingNumber(n int) bool {
	remainder := n
	var new int

	for remainder > 0 {
		d := remainder % 10
		switch d {
		case 0, 1, 8:
			new = new*10 + d
		case 6:
			new = new*10 + 9
		case 9:
			new = new*10 + 6
		default:
			return false
		}
		remainder /= 10
	}

	return new != n
}
