// https://leetcode.com/problems/long-pressed-name/

package main

func isLongPressedName(name string, typed string) bool {
	var ni, ti int
	var longPress byte
	for ni < len(name) && ti < len(typed) {
		if name[ni] == typed[ti] {
			longPress = typed[ti]
			ni++
			ti++
		} else if typed[ti] == longPress {
			ti++
		} else {
			return false
		}
	}
	for ti < len(typed) && typed[ti] == longPress {
		ti++
	}
	return ni == len(name) && ti == len(typed)
}
