// https://leetcode.com/problems/letter-combinations-of-a-phone-number/
// N is the length of digits
// time complexity is O(N*4^N)
// N from deep copy comb, 4^N because 7 and 9 digits have four letters.
// space complexity is O(N)

package main

func letterCombinations(digits string) []string {
	digitsToLetters := map[byte][]byte{
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
	var comb []byte
	backtrack17(&ans, comb, digits, digitsToLetters)
	return ans
}

func backtrack17(ans *[]string, comb []byte, digits string, digitsToLetters map[byte][]byte) {
	if len(comb) == len(digits) {
		*ans = append(*ans, string(comb))
		return
	}

	d := digits[len(comb)]
	for _, l := range digitsToLetters[d] {
		comb = append(comb, l)
		backtrack17(ans, comb, digits, digitsToLetters)
		comb = comb[:len(comb)-1]
	}
}
