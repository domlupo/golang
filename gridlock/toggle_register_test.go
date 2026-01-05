package main

import (
	"testing"
)

func TestToggleRegister(t *testing.T) {
	t.Run("test double toggle register panics", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("want panic but did not panic")
			}
		}()

		tr := newToggleRegister(true)
		tr.register(keyNoOp, func() {})
		tr.register(keyNoOp, func() {})
	})
}
