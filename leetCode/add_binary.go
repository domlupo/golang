// https://leetcode.com/problems/two-sum/

package main

func addBinary(a string, b string) string {
	var reversedSum []rune
	var carry bool

	for i := 1; i <= len(a) || i <= len(b) || carry; i++ {
		var sum int
		if i <= len(a) && a[len(a)-i] == '1' {
			sum++
		}
		if i <= len(b) && b[len(b)-i] == '1' {
			sum++
		}
		if carry {
			sum++
		}

		switch sum {
		case 0:
			reversedSum = append(reversedSum, '0')
			carry = false
		case 1:
			reversedSum = append(reversedSum, '1')
			carry = false
		case 2:
			carry = true
			reversedSum = append(reversedSum, '0')
		case 3:
			carry = true
			reversedSum = append(reversedSum, '1')
		}
	}

	var sum []rune
	for i := len(reversedSum) - 1; i >= 0; i-- {
		sum = append(sum, reversedSum[i])
	}
	return string(sum)
}
