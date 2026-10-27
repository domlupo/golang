// https://leetcode.com/problems/rotate-string/
// N is the length of s
// time complexity is O(N*N)
// space complexity is O(1)
//
// You can scan the string s+s for goal but that is O(N*N) time and O(N)
// space. There is complex Knuth-Morris-Pratt solution for O(N) time and
// O(N) space.

package main

func rotateString(s string, goal string) bool {
	if len(s) != len(goal) {
		return false
	}

	var shift int
	for shift < len(s) {
		for i := 0; i < len(s); i++ {
			j := ((i + shift) % len(s))
			if s[i] != goal[j] {
				break
			} else if i == len(s)-1 {
				return true
			}
		}
		shift++
	}

	return false
}
