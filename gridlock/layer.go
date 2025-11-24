// A layer is a 2 dimensional plane which units occupy.

package main

import (
	"fmt"
	"log/slog"
)

// A layer is a 2 dimensional plane which units occupy. Units occupy a unique position
// within a single layer. Units can interact with units in other layers. Units never switch layers.
type layer int

const (
	lNone layer = iota
	lLand
	lOrdnance
	lAircraft
	lBallistic
	lBombGravity
	lChaff
)

func (l layer) String() string {
	switch l {
	case lNone:
		return "noLayer"
	case lLand:
		return "land"
	case lOrdnance:
		return "ordnance"
	case lAircraft:
		return "aircraft"
	case lBallistic:
		return "ballistic"
	case lBombGravity:
		return "bombGravity"
	case lChaff:
		return "chaff"
	default:
		panic(fmt.Sprintf("Cannot String layer: %d", l))
	}
}

func (l layer) log() slog.Attr {
	return slog.String("layer", l.String())
}

type layers []layer

func (layers layers) String() string {
	var s string
	for i, l := range layers {
		s += l.String()
		if i != len(layers)-1 {
			s += ", "
		}
	}
	return s
}
