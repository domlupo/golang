// Misc code which is used in several other files but not large enough to need its own file.

package main

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/gofrs/uuid/v5"
)

const maxPercent = 100
const minPercent = 0

// One liner to create int pointer. Avoids needing two lines: 1. set int 2. ref int.
func intPtr(i int) *int {
	return &i
}

// Returns true if floats are approximately equal.
func equalFloats(f1, f2 float64) bool {
	// arbitrarily chose an epsilon value. I suspect issues caused by too small of an epsilon are
	// quicker to debug than too large of an epsilon. So I went with 1e-9 instead of 1e-6.
	epsilon := 1e-9
	return math.Abs(f1-f2) < epsilon
}

// Stores a callback to invoke when a Pixel is reached.
type action struct {
	// TODO: I am not sure if this should be a Pixel or a Pos. I think this will be more clear in the
	// future once target is more structured so I will punt it. If i had to decide now, I think Pixel
	// is better because 1. Pixel uses ints which are simpler than float64s and 2. actions will
	// probably be user input which is at the Pixel level not the Pos level.
	Pixel
	do func()
}

// Stores a callback to invoke when a move is confirmed.
type move struct {
	Pos
	confirm func()
}

// Wraps a third party universally unique identifier.
type id struct {
	uuid.UUID
}

// Generate a third party universally unique identifier.
func generateId() id {
	uuid, err := uuid.NewV4()
	if err != nil {
		logger.log(
			errorLevel,
			noLogType,
			"generateId",
			slog.String("failed to generate UUID", err.Error()))
		panic("generateId: cannot create ID. Game logic requires unique IDs")
	}
	return id{uuid}
}

// Typical log method.
func (action action) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("Pixel", action.Pixel.log()))
	return a
}

type player int

const (
	pNone player = iota
	pNeutral
	pOne
	pTwo
)

// Typical string method.
func (p player) string() string {
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
		panic(fmt.Sprintf("player.string unknown player: %d", p))
	}
}

type name int

const (
	noName name = iota
	infantry
	fighter
	bomber
	bombGravity
	missileHeatSeeking
	missileRadarActive
	missileBallistic
	missileCruise
	missileAntiBallistic
	flare
	chaff
)

func names() []name {
	return []name{
		infantry,
		fighter,
		bomber,
		bombGravity,
		missileHeatSeeking,
		missileRadarActive,
		missileBallistic,
		missileCruise,
		missileAntiBallistic,
		flare,
		chaff}
}

// Typical string method.
func (n name) string() string {
	switch n {
	case noName:
		return "noName"
	case infantry:
		return "infantry"
	case fighter:
		return "fighter"
	case bomber:
		return "bomber"
	case bombGravity:
		return "bombGravity"
	case missileHeatSeeking:
		return "missileHeatSeeking"
	case missileRadarActive:
		return "missileRadarActive"
	case missileBallistic:
		return "missileBallistic"
	case missileCruise:
		return "missileCruise"
	case missileAntiBallistic:
		return "missileAntiBallistic"
	case flare:
		return "flare"
	case chaff:
		return "chaff"
	default:
		panic(fmt.Sprintf("name.string unknown name %d", n))
	}
}
