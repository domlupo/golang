// https://leetcode.com/problems/generate-parentheses/

package main

func generateParenthesis(n int) []string {
	return backtrack("", 0, 0, n)
}

func backtrack(s string, lCount, rCount, pairCount int) []string {
	var ans []string

	if len(s) == 2*pairCount {
		ans = append(ans, s)
	}
	if lCount < pairCount {
		ans = append(ans, backtrack(s+"(", lCount+1, rCount, pairCount)...)
	}
	if rCount < lCount {
		ans = append(ans, backtrack(s+")", lCount, rCount+1, pairCount)...)
	}

	return ans
}

// n == 1
// ( -> ()

// n == 2
// (() -> (())
// ()( -> ()()

// n == 3
// ((( -> ((()))
// (() -> (())() and (()())
// ()( -> ()()() and ()(()
