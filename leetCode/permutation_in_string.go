// https://leetcode.com/problems/permutation-in-string/
// N is the length of s2
// time complexity is O(N)
// space complexity is O(1)

package main

// Very optimized over other solution.
// Count is used to compare a single value rather than compare the 26 values in an array.

func checkInclusion(s1 string, s2 string) bool {
	if len(s1) > len(s2) {
		return false
	}

	var s1Freq, s2Freq [26]int
	for i := 0; i < len(s1); i++ {
		s1Freq[s1[i]-'a']++
		s2Freq[s2[i]-'a']++
	}

	var count int
	for c := 'a'; c <= 'z'; c++ {
		if s1Freq[c-'a'] == s2Freq[c-'a'] {
			count++
		}
	}

	for i := 0; i < len(s2)-len(s1); i++ {
		if count == 26 {
			return true
		}

		r := s2[i+len(s1)] - 'a'
		l := s2[i] - 'a'
		s2Freq[r]++
		if s1Freq[r] == s2Freq[r] {
			count++
		} else if s1Freq[r] == s2Freq[r]-1 {
			count--
		}

		s2Freq[l]--
		if s1Freq[l] == s2Freq[l] {
			count++
		} else if s1Freq[l] == s2Freq[l]+1 {
			count--
		}
	}

	return count == 26
}

/*

less optimized, but same complexities and easier to read

func checkInclusion(s1 string, s2 string) bool {
	var s1Freq, s2Freq [26]int
	for _, r := range s1 {
		s1Freq[r-'a']++
	}

	var l int
	for r := 0; r < len(s2); r++ {
		s2Freq[s2[r]-'a']++
		if r-l == len(s1) {
			s2Freq[s2[l]-'a']--
			l++
		}
		if s1Freq == s2Freq {
			return true
		}
	}

	return false
}

*/
