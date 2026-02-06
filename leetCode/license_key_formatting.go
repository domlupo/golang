package main

func licenseKeyFormatting(s string, k int) string {
	var r []rune
	var count int

	for i := len(s) - 1; i >= 0; i-- {
		c := rune(s[i])
		if c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		if count == k {
			r = append(r, '-', c)
			count = 0
		} else {
			r = append(r, c)
		}
		count++
	}

	var o []rune
	for i := len(r) - 1; i >= 0; i-- {
		o = append(o, r[i])
	}

	return string(o)
}
