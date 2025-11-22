// Configurations for different unit types.

package main

import (
	"fmt"
	"log/slog"
	"math"
)

type name int

const (
	noName name = iota
	infantry
	unguidedGroundMissile
	guidedAirMissile
	fighter
	bomber
	bomb
	missile
	missileHeatSeeking
	flare
)

func (n name) String() string {
	switch n {
	case noName:
		return "noName"
	case infantry:
		return "infantry"
	case unguidedGroundMissile:
		return "unguidedGroundMissile"
	case guidedAirMissile:
		return "guidedAirMissile"
	case fighter:
		return "fighter"
	case bomber:
		return "bomber"
	case bomb:
		return "bomb"
	case missile:
		return "missile"
	case missileHeatSeeking:
		return "missileHeatSeeking"
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
	targetsLand          bool
	targetsOrdnance      bool
	targetsAircraft      bool
	collidesLand         bool
	collidesOrdnance     bool
	collidesAircraft     bool
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
	missileCount         int
	missileHeatCount     int
	missileCountMax      int
	missileHeatCountMax  int
	fireOrdnanceDelay    uint64
	heated               bool
	heatSeeking          bool
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
	targetsLand:     true,
	collidesLand:    true,
	explodeRadius:   0,
	maxAlivePerNode: math.MaxInt,
	speedFactor:     2.0,
}

// TODO: Figure out better naming for various missiles
// TODO: Figure out better attack conventions e.g. direct, vs explosion, etc
var unguidedGroundMissileConfig = unitConfig{
	name:                 unguidedGroundMissile,
	radius:               5,
	maxHealth:            1,
	attack:               1,
	vision:               100,
	produceLayer:         ordnance,
	collidesOrdnance:     true,
	staticTarget:         true,
	explodeAtTerminalPos: true,
	despawnAtTerminalPos: true,
	explodeRadius:        20,
	maxAlivePerNode:      math.MaxInt,
	speedFactor:          1.0,
}

var guidedAirMissileConfig = unitConfig{
	name:                 guidedAirMissile,
	radius:               5,
	maxHealth:            1,
	attack:               1000,
	vision:               300,
	produceLayer:         ordnance,
	collidesOrdnance:     true,
	targetsOrdnance:      true,
	targetsAircraft:      true,
	despawnAtTerminalPos: true,
	explodeRadius:        100,
	maxAlivePerNode:      math.MaxInt,
	speedFactor:          1,
}

var fighterConfig = unitConfig{
	name:                fighter,
	radius:              8,
	maxHealth:           1,
	vision:              800,
	produceLayer:        aircraft,
	targetsOrdnance:     true,
	targetsAircraft:     true,
	collidesAircraft:    true,
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
}

var bomberConfig = unitConfig{
	name:              bomber,
	radius:            20,
	maxHealth:         1,
	vision:            200,
	produceLayer:      aircraft,
	targetsOrdnance:   true,
	collidesAircraft:  true,
	explodeRadius:     100,
	maxAlivePerNode:   1,
	bombCountMax:      4,
	flareCountMax:     1000,
	fireOrdnanceDelay: 30,
	speedFactor:       1.0,
	heated:            true,
	reactive:          true,
	requester:         true,
}

var bombConfig = unitConfig{
	name:             bomb,
	radius:           4,
	maxHealth:        1,
	produceLayer:     ordnance,
	collidesOrdnance: true,
	staticTarget:     true,
	targetsLand:      true,
	explodeRadius:    20,
	explodeStepDelay: 300,
	maxAlivePerNode:  math.MaxInt,
	speedFactor:      1.0,
}

var missileConfig = unitConfig{
	name:                 missile,
	radius:               3,
	maxHealth:            1,
	vision:               800,
	produceLayer:         ordnance,
	collidesOrdnance:     true,
	staticTarget:         true,
	despawnAtTerminalPos: true,
	explodeRadius:        20,
	speedFactor:          0.5,
}

var missileHeatSeekingConfig = unitConfig{
	name:                 missileHeatSeeking,
	radius:               3,
	maxHealth:            1,
	vision:               800,
	produceLayer:         ordnance,
	collidesOrdnance:     true,
	targetsOrdnance:      true,
	staticTarget:         true,
	despawnAtTerminalPos: true,
	heatSeeking:          true,
	targetsAircraft:      true,
	collidesAircraft:     true,
	explodeRadius:        20,
	speedFactor:          0.5,
}

var flareConfig = unitConfig{
	name:                 flare,
	radius:               3,
	maxHealth:            1,
	produceLayer:         ordnance,
	collidesOrdnance:     true,
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
	a = append(a, slog.Bool("targetsLand", uc.targetsLand))
	a = append(a, slog.Bool("targetsOrdnance", uc.targetsOrdnance))
	a = append(a, slog.Bool("targetsAircraft", uc.targetsAircraft))
	a = append(a, slog.Bool("explodeAtTerminalPos", uc.explodeAtTerminalPos))
	a = append(a, slog.Bool("despawnAtTerminalPos", uc.despawnAtTerminalPos))
	a = append(a, slog.Int("missileCount", uc.missileCount))
	a = append(a, slog.Int("missileHeatCount", uc.missileHeatCount))
	return a
}
