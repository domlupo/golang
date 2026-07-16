// https://leetcode.com/problems/letter-combinations-of-a-phone-number/
// N is the length of digits
// time complexity is O(N*4^N)
// N from deep copy comb, 4^N because 7 and 9 digits have four letters.
// space complexity is O(N)

package main

func letterCombinations(digits string) []string {
	numToLetters := map[byte][]byte{
		'2': {'a', 'b', 'c'},
		'3': {'d', 'e', 'f'},
		'4': {'g', 'h', 'i'},
		'5': {'j', 'k', 'l'},
		'6': {'m', 'n', 'o'},
		'7': {'p', 'q', 'r', 's'},
		'8': {'t', 'u', 'v'},
		'9': {'w', 'x', 'y', 'z'},
	}

	var ans []string
	comb := ""
	backtrack17(&ans, comb, digits, numToLetters)
	return ans
}

func backtrack17(ans *[]string, comb, digits string, numToLetters map[byte][]byte) {
	if len(digits) == 0 {
		*ans = append(*ans, comb)
		return
	}

	letters := numToLetters[digits[0]]
	for _, l := range letters {
		comb = comb + string(l)
		backtrack17(ans, comb, digits[1:], numToLetters)
		comb = comb[:len(comb)-1]
	}
}
