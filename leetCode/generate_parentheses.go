// https://leetcode.com/problems/generate-parentheses/
// https://en.wikipedia.org/wiki/Catalan_number
// N is passed n
// time complexity is O((N * 4^N)/n^3/2)) -> O(4^N/sqrt(N))
// space complexity is O(N)

package main

func generateParenthesis(n int) []string {
	var ans []string
	backtrack(&ans, "", 0, 0, n*2)
	return ans
}

func backtrack(ans *[]string, str string, l, r, sum int) {
	if r > l || l > sum/2 {
		return
	} else if l+r == sum {
		*ans = append(*ans, str)
		return
	}

	backtrack(ans, str+"(", l+1, r, sum)
	backtrack(ans, str+")", l, r+1, sum)
}

// cancel if right is greater than left or left greater than half

// n = 2
// ( -> ok
// ) -> cancel

// () -> ok
// (( -> ok

// ()( -> ok
// ()) -> cancel
// (() -> ok
// ((( -> cancel

// ()() -> ok
// ()(( -> cancel
// (()) -> ok
// (()( -> cancel
