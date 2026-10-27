// https://leetcode.com/problems/remove-all-adjacent-duplicates-in-string/
// N is the length of s
// time complexity is O(N)
// space complexity is O(N)

package main

func removeDuplicates1047(s string) string {
	var stack []rune
	for _, r := range s {
		if len(stack) > 0 && stack[len(stack)-1] == r {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, r)
		}
	}
	return string(stack)
}
