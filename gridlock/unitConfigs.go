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
	attacksLand          bool
	attacksLowAir        bool
	attacksHighAir       bool
	collidesLand         bool
	collidesLowAir       bool
	collidesHighAir      bool
	despawnAtTerminalPos bool
	explodeAtTerminalPos bool
	explodeRadius        int
	explodeStepDelay     uint64
	maxAlivePerNode      int
	bombs                int
	speedFactor          float64
	reactions            []reaction
	reactive             bool
	requests             []request
	requester            bool
	missileCount         int
	missileCountMax      int
	missileFireDelay     uint64
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
	attacksLand:     true,
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
	maxHealth:            100,
	attack:               1,
	vision:               100,
	reach:                0,
	produceLayer:         lowAir,
	collidesLowAir:       true,
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
	maxHealth:            100,
	attack:               1000,
	vision:               300,
	reach:                0,
	produceLayer:         lowAir,
	collidesLowAir:       true,
	attacksLowAir:        true,
	attacksHighAir:       true,
	despawnAtTerminalPos: true,
	explodeRadius:        100,
	maxAlivePerNode:      math.MaxInt,
	speedFactor:          1,
}

var fighterConfig = unitConfig{
	name:          fighter,
	radius:        8,
	maxHealth:     100,
	attack:        0,
	vision:        400,
	reach:         0,
	produceLayer:  highAir,
	attacksLowAir: true,
	//attacksHighAir:   true,
	collidesHighAir:  true,
	explodeRadius:    100,
	maxAlivePerNode:  1,
	bombs:            1,
	speedFactor:      1.0,
	reactive:         true,
	requester:        true,
	missileCountMax:  100,
	missileFireDelay: 10,
}

var bomberConfig = unitConfig{
	name:            bomber,
	radius:          20,
	maxHealth:       100,
	attack:          0,
	vision:          40,
	reach:           0,
	produceLayer:    highAir,
	collidesHighAir: true,
	explodeRadius:   100,
	maxAlivePerNode: 1,
	bombs:           4,
	speedFactor:     1.0,
}

var bombConfig = unitConfig{
	name:             bomb,
	radius:           4,
	maxHealth:        100,
	attack:           0,
	vision:           0,
	reach:            0,
	produceLayer:     lowAir,
	collidesLowAir:   true,
	staticTarget:     true,
	attacksLand:      true,
	explodeRadius:    20,
	explodeStepDelay: 1000,
	maxAlivePerNode:  math.MaxInt,
	speedFactor:      1.0,
}

var missileConfig = unitConfig{
	name:                 missile,
	radius:               3,
	maxHealth:            100,
	attack:               0,
	vision:               0,
	reach:                0,
	produceLayer:         lowAir,
	collidesLowAir:       true,
	collidesHighAir:      false,
	staticTarget:         true,
	despawnAtTerminalPos: true,
	explodeRadius:        20,
	explodeStepDelay:     0,
	maxAlivePerNode:      math.MaxInt,
	speedFactor:          0.5,
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
	a = append(a, slog.Bool("attacksLand", uc.attacksLand))
	a = append(a, slog.Bool("attacksLowAir", uc.attacksLowAir))
	a = append(a, slog.Bool("attacksHighAir", uc.attacksHighAir))
	a = append(a, slog.Bool("explodeAtTerminalPos", uc.explodeAtTerminalPos))
	a = append(a, slog.Bool("despawnAtTerminalPos", uc.despawnAtTerminalPos))
	return a
}
