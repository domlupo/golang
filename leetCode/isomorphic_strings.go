// https://leetcode.com/problems/isomorphic-strings/

package main

func isIsomorphic(s string, t string) bool {
	sToT := make(map[byte]byte)
	tToS := make(map[byte]byte)

	for i := 0; i < len(s); i++ {
		sC := s[i]
		tC := t[i]

		sMap, sOk := sToT[sC]
		tMap, tOk := tToS[tC]

		if sOk && sMap != tC {
			return false
		}
		if tOk && tMap != sC {
			return false
		}

		sToT[sC] = tC
		tToS[tC] = sC
	}

	return true
}
