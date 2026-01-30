// https://leetcode.com/problems/valid-parentheses/description/

package main

import (
	"testing"
)

func TestValidParentheses(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		s      string
		output bool
	}{
		{"()", true},
		{"()[]{}", true},
		{"(]", false},
		{"([])", true},
		{"([)]", false},
	}

	for _, tt := range tests {
		t.Run("validParentheses", func(t *testing.T) {
			if isValid(tt.s) != tt.output {
				t.Errorf("incorrect")
			}
		})
	}
}
