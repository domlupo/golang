// https://leetcode.com/problems/two-sum/

package main

func addBinary(a string, b string) string {
	var reversed []rune
	var carry bool

	for i := 0; i < len(a) || i < len(b); i++ {
		total := 0
		if i < len(a) && a[len(a)-1-i] == '1' {
			total++
		}
		if i < len(b) && b[len(b)-1-i] == '1' {
			total++
		}
		if carry {
			total++
		}

		carry = false
		switch total {
		case 0:
			reversed = append(reversed, '0')
		case 1:
			reversed = append(reversed, '1')
		case 2:
			carry = true
			reversed = append(reversed, '0')
		case 3:
			carry = true
			reversed = append(reversed, '1')
		default:
			panic("addBinary: total cannot exceed 3")
		}
	}

	if carry {
		reversed = append(reversed, '1')
	}

	var sum []rune
	for i := len(reversed) - 1; i >= 0; i-- {
		sum = append(sum, reversed[i])
	}
	return string(sum)
}
