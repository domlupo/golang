// https://leetcode.com/problems/letter-combinations-of-a-phone-number/
// N is the length of digits
// time complexity is O(N*4^N)
// N from deep copy comb, 4^N because 7 and 9 digits have four letters.
// space complexity is O(N)

package main

func letterCombinations(digits string) []string {
	numToLetter := map[byte][]string{
		'2': {"a", "b", "c"},
		'3': {"d", "e", "f"},
		'4': {"g", "h", "i"},
		'5': {"j", "k", "l"},
		'6': {"m", "n", "o"},
		'7': {"p", "q", "r", "s"},
		'8': {"t", "u", "v"},
		'9': {"w", "x", "y", "z"},
	}

	var ans []string
	backtrack17(&ans, numToLetter, "", digits)
	return ans
}

func backtrack17(ans *[]string, numToLetter map[byte][]string, comb, digits string) {
	if len(digits) == 0 {
		*ans = append(*ans, comb)
		return
	}

	for _, s := range numToLetter[digits[0]] {
		comb = comb + s
		backtrack17(ans, numToLetter, comb, digits[1:])
		comb = comb[:len(comb)-1]
	}
}
