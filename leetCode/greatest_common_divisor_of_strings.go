// https://leetcode.com/problems/greatest-common-divisor-of-strings/

package main

func gcdOfStrings(str1 string, str2 string) string {
	if str1+str2 != str2+str1 {
		return ""
	}

	gcd, rem := len(str1), len(str2)
	for rem != 0 {
		gcd, rem = rem, gcd%rem
	}

	return str1[0:gcd]
}
