// https://leetcode.com/problems/two-sum/

package main

func addBinary(a string, b string) string {
	var reversedSum []rune
	var carry bool

	for i := 0; ; i++ {
		if i >= len(a) && i >= len(b) && !carry {
			break
		}

		total := 0
		if i < len(a) && a[len(a)-i-1] == '1' {
			total++
		}
		if i < len(b) && b[len(b)-i-1] == '1' {
			total++
		}
		if carry {
			total++
		}

		carry = false
		switch total {
		case 0:
			reversedSum = append(reversedSum, '0')
		case 1:
			reversedSum = append(reversedSum, '1')
		case 2:
			carry = true
			reversedSum = append(reversedSum, '0')
		case 3:
			carry = true
			reversedSum = append(reversedSum, '1')
		default:
			panic("total cannot exceed 3")
		}
	}

	var sum []rune
	for i := len(reversedSum) - 1; i >= 0; i-- {
		sum = append(sum, reversedSum[i])
	}
	return string(sum)
}
