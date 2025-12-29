// Misc code which is used in several other files but not large enough to need its own file.

package main

import (
	"fmt"
	"log/slog"

	"github.com/gofrs/uuid/v5"
)

const maxPercent = 100
const minPercent = 0

// Stores a callback to invoke when a position is reached.
type action struct {
	position
	do func()
}

// Stores a callback to invoke when a move is confirmed.
type move struct {
	position
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
			errorLevel, noLogType, "generateId", slog.String("failed to generate UUID", err.Error()))
		panic("error when calling generateId. Game logic requires unique Ids")
	}
	return id{uuid}
}

func (action action) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, nest("pos", action.position.log()))
	return a
}

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
		panic(fmt.Sprintf("String name is unknown %d", n))
	}
}
