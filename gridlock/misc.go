// Misc code which is used in several other files but not large enough to need its own file.

package main

import (
	"fmt"
	"log/slog"
)

const maxPercent = 100
const minPercent = 0

type player int

const (
	pNone player = iota
	pNeutral
	pOne
	pTwo
)

func (p player) String() string {
	switch p {
	case pNone:
		return "pNone"
	case pNeutral:
		return "pNeutral"
	case pOne:
		return "pOne"
	case pTwo:
		return "pTwo"
	default:
		panic(fmt.Sprintf("Cannot convert player to string: %d", p))
	}
}

type windowConfig struct {
	width  int
	height int
}

func newWindowConfig() windowConfig {
	wc := windowConfig{
		width:  1280,
		height: 720,
	}
	return wc
}

func (wc windowConfig) maxPixels() int {
	return wc.height * wc.width
}

func (wc windowConfig) inbounds(p PixelPosition) bool {
	return !(p.X < 0 || p.X > wc.width-1 || p.Y < 0 || p.Y > wc.height-1)
}

// Stores a callback to invoke when a step is reached.
type event struct {
	step uint64
	do   func()
}

// Stores a callback to invoke when a position is reached.
type action struct {
	pos position
	do  func()
}

// Stores a callback to invoke when a move is confirmed.
type Move struct {
	position
	confirm func()
}

func (action action) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("pos", action.pos.log()))
	return a
}
