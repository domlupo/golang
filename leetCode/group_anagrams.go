// https://leetcode.com/problems/group-anagrams/
// N is length of strs
// K is max length of strs[i]
// time complexity is O(N*K log(K))
// space complexity is O(N*K)

package main

func groupAnagrams(strs []string) [][]string {
	anagramGroups := make(map[string][]string)

	for _, str := range strs {
		k := sort49(str)
		anagramGroups[k] = append(anagramGroups[k], str)
	}

	var ans [][]string
	for _, group := range anagramGroups {
		ans = append(ans, group)
	}
	return ans
}

func sort49(str string) string {
	s := []byte(str)
	for i := 0; i < len(s)-1; i++ {
		if i >= 0 && s[i] > s[i+1] {
			s[i], s[i+1] = s[i+1], s[i]
			i -= 2
		}
	}
	return string(s)
}
