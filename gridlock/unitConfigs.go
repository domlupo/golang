// Configurations for different unit types.

package main

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
)

type name int

const (
	noName name = iota
	infantry
	fighter
	bomber
	bomb
	missileHeatSeeking
	missileBallistic
	missileCruise
	missileAntiBallistic
	flare
)

func (n name) String() string {
	switch n {
	case noName:
		return "noName"
	case infantry:
		return "infantry"
	case fighter:
		return "fighter"
	case bomber:
		return "bomber"
	case bomb:
		return "bomb"
	case missileHeatSeeking:
		return "missileHeatSeeking"
	case missileBallistic:
		return "missileBallistic"
	case missileCruise:
		return "missileCruise"
	case missileAntiBallistic:
		return "missileAntiBallistic"
	case flare:
		return "flare"
	default:
		panic(fmt.Sprintf("String name is unknown %d", n))
	}
}

type unitConfig struct {
	name                 name
	radius               int
	maxHealth            int
	attack               int
	vision               int
	reach                int
	produceLayer         layer
	staticTarget         bool
	targetLayers         layers
	collisionLayers      layers
	despawnAtTerminalPos bool
	explodeAtTerminalPos bool
	explodeRadius        int
	explodeStepDelay     uint64
	reloadPayloadDelay   uint64
	maxAlivePerNode      int
	speedFactor          float64
	reactions            []reaction
	reactive             bool
	requests             []request
	requester            bool
	bombCount            int
	bombCountMax         int
	flareCount           int
	flareCountMax        int
	missileHeatCount     int
	missileCountMax      int
	missileHeatCountMax  int
	fireOrdnanceDelay    uint64
	heated               bool
	heatSeeking          bool
	ignoreTarget         func(u, t *unit) bool
}

// Note: unit radius set to an odd number is better for testing due to division bugs
var infantryConfig = unitConfig{
	name:            infantry,
	radius:          3,
	maxHealth:       50,
	attack:          1,
	vision:          30,
	reach:           5,
	produceLayer:    land,
	targetLayers:    []layer{land},
	collisionLayers: []layer{land},
	explodeRadius:   0,
	maxAlivePerNode: math.MaxInt,
	speedFactor:     2.0,
}

var fighterConfig = unitConfig{
	name:                fighter,
	radius:              8,
	maxHealth:           1,
	vision:              800,
	produceLayer:        aircraft,
	targetLayers:        []layer{ordnance, aircraft},
	collisionLayers:     []layer{aircraft},
	explodeRadius:       100,
	maxAlivePerNode:     1,
	bombCountMax:        1,
	speedFactor:         1.0,
	reactive:            true,
	requester:           true,
	missileCountMax:     10,
	missileHeatCountMax: 1000,
	fireOrdnanceDelay:   50,
	heated:              true,
	ignoreTarget: func(u, t *unit) bool {
		return t.name == bomb || t.name == flare
	},
}

var bomberConfig = unitConfig{
	name:              bomber,
	radius:            20,
	maxHealth:         1,
	vision:            200,
	produceLayer:      aircraft,
	targetLayers:      []layer{ordnance},
	collisionLayers:   []layer{aircraft},
	explodeRadius:     100,
	maxAlivePerNode:   1,
	bombCountMax:      4,
	flareCountMax:     1000,
	fireOrdnanceDelay: 30,
	speedFactor:       1.0,
	heated:            true,
	reactive:          true,
	requester:         true,
	ignoreTarget:      func(u, t *unit) bool { return slices.Contains(u.skipTargets, t.id) },
}

var bombConfig = unitConfig{
	name:             bomb,
	radius:           4,
	maxHealth:        1,
	produceLayer:     ordnance,
	targetLayers:     []layer{land},
	collisionLayers:  []layer{ordnance},
	staticTarget:     true,
	explodeRadius:    20,
	explodeStepDelay: 300,
	maxAlivePerNode:  math.MaxInt,
	speedFactor:      1.0,
}

var missileHeatSeekingConfig = unitConfig{
	name:                 missileHeatSeeking,
	radius:               3,
	maxHealth:            1,
	vision:               300,
	produceLayer:         ordnance,
	targetLayers:         []layer{ordnance, aircraft},
	collisionLayers:      []layer{ordnance, aircraft},
	staticTarget:         true,
	despawnAtTerminalPos: true,
	heatSeeking:          true,
	explodeRadius:        20,
	speedFactor:          0.5,
	heated:               true,
	maxAlivePerNode:      math.MaxInt,
	ignoreTarget:         func(u, t *unit) bool { return !t.heated },
}

var missileCruiseConfig = unitConfig{
	name:                 missileCruise,
	radius:               3,
	maxHealth:            1,
	vision:               1,
	produceLayer:         ordnance,
	collisionLayers:      []layer{ordnance},
	staticTarget:         true,
	explodeAtTerminalPos: true,
	explodeRadius:        20,
	speedFactor:          1,
	maxAlivePerNode:      math.MaxInt,
	heated:               true,
}

var missileBallisticConfig = unitConfig{
	name:                 missileBallistic,
	radius:               3,
	maxHealth:            1,
	vision:               1,
	produceLayer:         ballistic,
	collisionLayers:      []layer{ballistic},
	staticTarget:         true,
	explodeAtTerminalPos: true,
	explodeRadius:        20,
	speedFactor:          1,
	maxAlivePerNode:      math.MaxInt,
}

var missileAntiBallisticConfig = unitConfig{
	name:                 missileAntiBallistic,
	radius:               3,
	maxHealth:            1,
	vision:               300,
	produceLayer:         ballistic,
	targetLayers:         []layer{ballistic},
	collisionLayers:      []layer{ballistic},
	staticTarget:         true,
	despawnAtTerminalPos: true,
	explodeRadius:        20,
	speedFactor:          1,
	maxAlivePerNode:      math.MaxInt,
}

var flareConfig = unitConfig{
	name:                 flare,
	radius:               3,
	maxHealth:            1,
	produceLayer:         ordnance,
	collisionLayers:      []layer{ordnance},
	staticTarget:         true,
	despawnAtTerminalPos: true,
	speedFactor:          0.5,
	heated:               true,
}

func (uc unitConfig) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.String("name", uc.name.String()))
	a = append(a, slog.Int("radius", uc.radius))
	a = append(a, slog.Int("maxHealth", uc.maxHealth))
	a = append(a, slog.Int("attack", uc.attack))
	a = append(a, slog.Int("vision", uc.vision))
	a = append(a, slog.Int("reach", uc.reach))
	a = append(a, uc.produceLayer.log())
	a = append(a, slog.Bool("staticTarget", uc.staticTarget))
	a = append(a, slog.Bool("heatSeeking", uc.heatSeeking))
	a = append(a, slog.String("targetLayers", uc.targetLayers.String()))
	a = append(a, slog.String("collisionLayers", uc.collisionLayers.String()))
	a = append(a, slog.Bool("explodeAtTerminalPos", uc.explodeAtTerminalPos))
	a = append(a, slog.Bool("despawnAtTerminalPos", uc.despawnAtTerminalPos))
	a = append(a, slog.Int("missileHeatCount", uc.missileHeatCount))
	return a
}
