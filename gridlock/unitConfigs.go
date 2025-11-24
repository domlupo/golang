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
	bombGravity
	missileHeatSeeking
	missileRadarActive
	missileBallistic
	missileCruise
	missileAntiBallistic
	flare
	chaff
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

type unitConfig struct {
	name                       name
	radius                     int
	maxHealth                  int
	attack                     int
	vision                     int
	reach                      int
	produceLayer               layer
	staticTarget               bool
	targetLayers               layers
	collisionLayers            layers
	despawnAtTerminalPos       bool
	explodeAtTerminalPos       bool
	explodeRadius              int
	eventStepDelay             uint64
	reloadPayloadDelay         uint64
	maxAlivePerNode            int
	speedFactor                float64
	reactions                  []reaction
	reactive                   bool
	requests                   []request
	requester                  bool
	bombGravityCount           int
	bombGravityCountMax        int
	flareCount                 int
	flareCountMax              int
	missileHeatSeekingCount    int
	missileHeatSeekingCountMax int
	missileRadarActiveCount    int
	missileRadarActiveCountMax int
	fireOrdnanceDelay          uint64
	heated                     bool
	heatSeeking                bool
	ignore                     func(u, t *unit) bool
	ignoreTarget               func(t *unit) bool
}

// Note: unit radius set to an odd number is better for testing due to division bugs
var infantryConfig = unitConfig{
	name:            infantry,
	radius:          3,
	maxHealth:       50,
	attack:          1,
	vision:          30,
	reach:           5,
	produceLayer:    lLand,
	targetLayers:    []layer{lLand},
	collisionLayers: []layer{lLand},
	explodeRadius:   0,
	maxAlivePerNode: math.MaxInt,
	speedFactor:     2.0,
}

var fighterConfig = unitConfig{
	name:                       fighter,
	radius:                     8,
	maxHealth:                  1,
	vision:                     800,
	produceLayer:               lAircraft,
	targetLayers:               []layer{lOrdnance, lAircraft},
	collisionLayers:            []layer{lAircraft},
	explodeRadius:              100,
	maxAlivePerNode:            1,
	bombGravityCountMax:        1,
	speedFactor:                1.0,
	reactive:                   true,
	requester:                  true,
	missileHeatSeekingCountMax: 0,
	missileRadarActiveCountMax: 1000,
	fireOrdnanceDelay:          50,
	heated:                     true,
	ignoreTarget: func(t *unit) bool {
		return t.name == flare
	},
}

var bomberConfig = unitConfig{
	name:                bomber,
	radius:              20,
	maxHealth:           1,
	vision:              200,
	produceLayer:        lAircraft,
	targetLayers:        []layer{lOrdnance},
	collisionLayers:     []layer{lAircraft},
	explodeRadius:       100,
	maxAlivePerNode:     1,
	bombGravityCountMax: 4,
	flareCountMax:       1000,
	fireOrdnanceDelay:   30,
	speedFactor:         1.0,
	heated:              true,
	reactive:            true,
	requester:           true,
	ignore:              func(u, t *unit) bool { return slices.Contains(u.skipTargets, t.id) },
}

var bombGravityConfig = unitConfig{
	name:            bombGravity,
	radius:          4,
	maxHealth:       1,
	produceLayer:    lBombGravity,
	targetLayers:    []layer{lLand},
	collisionLayers: []layer{lBombGravity},
	staticTarget:    true,
	explodeRadius:   20,
	eventStepDelay:  300,
	maxAlivePerNode: math.MaxInt,
}

var missileRadarActiveConfig = unitConfig{
	name:                 missileRadarActive,
	radius:               3,
	maxHealth:            1,
	vision:               300,
	produceLayer:         lOrdnance,
	targetLayers:         []layer{lOrdnance, lAircraft},
	collisionLayers:      []layer{lOrdnance, lAircraft},
	staticTarget:         true,
	despawnAtTerminalPos: true,
	speedFactor:          0.5,
	heated:               true,
	maxAlivePerNode:      math.MaxInt,
	ignoreTarget:         func(t *unit) bool { return t.name == flare },
}

var missileHeatSeekingConfig = unitConfig{
	name:                 missileHeatSeeking,
	radius:               3,
	maxHealth:            1,
	vision:               300,
	produceLayer:         lOrdnance,
	targetLayers:         []layer{lOrdnance, lAircraft},
	collisionLayers:      []layer{lOrdnance, lAircraft},
	staticTarget:         true,
	despawnAtTerminalPos: true,
	heatSeeking:          true,
	speedFactor:          0.5,
	heated:               true,
	maxAlivePerNode:      math.MaxInt,
	ignoreTarget:         func(t *unit) bool { return !t.heated },
}

var missileCruiseConfig = unitConfig{
	name:                 missileCruise,
	radius:               3,
	maxHealth:            1,
	vision:               1,
	produceLayer:         lOrdnance,
	collisionLayers:      []layer{lOrdnance},
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
	produceLayer:         lBallistic,
	collisionLayers:      []layer{lBallistic},
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
	produceLayer:         lBallistic,
	targetLayers:         []layer{lBallistic},
	collisionLayers:      []layer{lBallistic},
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
	produceLayer:         lOrdnance,
	collisionLayers:      []layer{lOrdnance},
	staticTarget:         true,
	despawnAtTerminalPos: true,
	speedFactor:          0.5,
	heated:               true,
}

var chaffConfig = unitConfig{
	name:            chaff,
	radius:          100,
	maxHealth:       1,
	produceLayer:    lChaff,
	collisionLayers: []layer{lChaff},
	eventStepDelay:  300,
	staticTarget:    true,
	maxAlivePerNode: math.MaxInt,
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
	a = append(a, slog.Int("missileHeatSeekingCount", uc.missileHeatSeekingCount))
	return a
}
