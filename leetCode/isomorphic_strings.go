// https://leetcode.com/problems/isomorphic-strings/

package main

func isIsomorphic(s string, t string) bool {
	st := make(map[byte]byte)
	ts := make(map[byte]byte)

	for i := 0; i < len(s); i++ {
		tMap, stOk := st[s[i]]
		sMap, tsOk := ts[t[i]]

		if !stOk && !tsOk {
			st[s[i]] = t[i]
			ts[t[i]] = s[i]
		} else if tMap != t[i] || sMap != s[i] {
			return false
		}
	}

	return true
}
