// https://leetcode.com/problems/generate-parentheses/
// https://en.wikipedia.org/wiki/Catalan_number
// N is passed n
// time complexity is O((N * 4^N)/n^3/2)) -> O(4^N/sqrt(N))
// space complexity is O(N)

package main

func generateParenthesis(n int) []string {
	var ans []string
	backtrack22(&ans, []byte{}, n*2, 0, 0)
	return ans
}

func backtrack22(ans *[]string, comb []byte, size, l, r int) {
	if len(comb) == size {
		*ans = append(*ans, string(comb))
		return
	}

	if l+1 <= size/2 {
		backtrack22(ans, append(comb, '('), size, l+1, r)
	}
	if r+1 <= l {
		backtrack22(ans, append(comb, ')'), size, l, r+1)
	}
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
