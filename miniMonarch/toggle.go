// file contains a struct which performs a respective action for a key press
package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const defaultCooldown = 300 * time.Millisecond

type callBack func()

type toggle struct {
	key      ebiten.Key
	callBack callBack
	timer    time.Time
	cooldown time.Duration
}

func newToggle(key ebiten.Key, callback callBack) *toggle {
	t := toggle{
		key:      key,
		callBack: callback,
		timer:    time.Now(),
		cooldown: defaultCooldown,
	}
	return &t
}

func (t *toggle) checkToggle(currentTime time.Time) {
	if currentTime.Sub(t.timer) > t.cooldown && ebiten.IsKeyPressed(t.key) {
		t.callBack()
		t.timer = currentTime
	}
}
