// https://leetcode.com/problems/compare-version-numbers/
// N is the max length of version1 and version 2 strings
// time complexity is O(N)
// space complexity is O(N)

package main

import "strconv"

func compareVersion(version1 string, version2 string) int {
	var p1, p2 int
	for p1 < len(version1) || p2 < len(version2) {
		ver1 := parseVersion(version1, &p1)
		ver2 := parseVersion(version2, &p2)
		if ver1 < ver2 {
			return -1
		} else if ver2 < ver1 {
			return 1
		}
	}

	return 0
}

func parseVersion(version string, i *int) int {
	ver := "0"
	for *i < len(version) && version[*i] != '.' {
		ver = ver + string(version[*i])
		*i++
	}
	*i++
	v, _ := strconv.Atoi(ver)
	return v
}
