// https://leetcode.com/problems/unique-email-addresses/

package main

func numUniqueEmails(emails []string) int {
	uniqEmails := make(map[string]bool)

	for _, email := range emails {
		var str string
		chars := []byte(email)
		var afterPlus bool
		for ci, c := range chars {
			if c == '@' {
				str = str + string(chars[ci:])
				break
			} else if c == '+' {
				afterPlus = true
			} else if c == '.' || afterPlus {
				continue
			} else {
				str = str + string(c)
			}
		}
		uniqEmails[string(str)] = true
	}

	return len(uniqEmails)
}
