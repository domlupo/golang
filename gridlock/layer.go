// A layer is a 2 dimensional plane which units occupy.

package main

import (
	"fmt"
	"log/slog"
	"strings"
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

func (l layer) string() string {
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
		panic(fmt.Sprintf("Cannot string layer: %d", l))
	}
}

func (l layer) log() slog.Attr {
	return slog.String("layer", l.string())
}

type layers []layer

func (layers layers) string() string {
	var s []string
	for i, l := range layers {
		s = append(s, l.string())
		if i != len(layers)-1 {
			s = append(s, ",")
		}
	}
	return strings.Join(s, "")
}
