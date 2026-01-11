// Contains all toggles registered. Use to register all new toggles and check all toggle presses.

package main

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Contains all toggles registered. Use to register all new toggles and check all toggle presses.
type toggleRegister struct {
	enableToggles bool
	toggles       map[key]*toggle
	cooldown      time.Duration
}

// Only create a single instance per application. Interacts with global state via keyboard input.
func newToggleRegister(enableToggles bool) *toggleRegister {
	tr := toggleRegister{
		enableToggles: enableToggles,
		toggles:       make(map[key]*toggle),
		cooldown:      300 * time.Millisecond,
	}
	return &tr
}

// Create a toggle with the passed key and callback. Panics if the key is already used by another
// toggle even if the same callback is registered.
func (tr toggleRegister) register(key key, callback func()) {
	if _, ok := tr.toggles[key]; ok {
		panic(fmt.Sprintf("toggleRegister.register: key: %d is taken", key))
	}

	var k ebiten.Key
	switch key {
	case keyNoOp:
		// ebiten doesn't provide a no op key. It defaults to ebiten.KeyA, must be handled elsewhere.
	case keyA:
		k = ebiten.KeyA
	case keyB:
		k = ebiten.KeyB
	case keyC:
		k = ebiten.KeyC
	case keyD:
		k = ebiten.KeyD
	case keyH:
		k = ebiten.KeyH
	case keyI:
		k = ebiten.KeyI
	case keyK:
		k = ebiten.KeyK
	case keyL:
		k = ebiten.KeyL
	case keyM:
		k = ebiten.KeyM
	case keyN:
		k = ebiten.KeyN
	case keyO:
		k = ebiten.KeyO
	case keyP:
		k = ebiten.KeyP
	case keyQ:
		k = ebiten.KeyQ
	case keyS:
		k = ebiten.KeyS
	case keyT:
		k = ebiten.KeyT
	case keyU:
		k = ebiten.KeyU
	case keyV:
		k = ebiten.KeyV
	case keyW:
		k = ebiten.KeyW
	case keyX:
		k = ebiten.KeyX
	default:
		panic("toggleRegister.register: unknown key")
	}

	tr.toggles[key] = &toggle{
		key:      k,
		callback: callback,
	}
}

// Activate any enabled toggles which were recently pressed.
func (tr toggleRegister) checkAllTogglePresses() {
	currentTime := time.Now()
	if tr.enableToggles {
		for k, t := range tr.toggles {
			if k == keyNoOp {
				return
			}

			if currentTime.Sub(t.timer) > tr.cooldown && inpututil.IsKeyJustPressed(t.key) {
				t.callback()
				t.timer = currentTime
			}
		}
	}
}

// Represents a single toggle.
type toggle struct {
	key      ebiten.Key
	callback func()
	timer    time.Time
}

type key int

const (
	keyNoOp key = iota
	keyA
	keyB
	keyC
	keyD
	keyH
	keyI
	keyK
	keyL
	keyM
	keyN
	keyO
	keyP
	keyQ
	keyS
	keyT
	keyU
	keyV
	keyW
	keyX
)
