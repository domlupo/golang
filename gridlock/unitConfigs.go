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
	bomber
	bomb
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
	case bomber:
		return "bomber"
	case bomb:
		return "bomb"
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
	attacksGround        bool
	attacksLowAltitude   bool
	attacksHighAltitude  bool
	despawnAtTerminalPos bool
	explodeAtTerminalPos bool
	explodeRadius        int
	explodeStepDelay     uint64
	maxAlivePerNode      int
}

// Note: unit radius set to an odd number is better for testing due to division bugs
var infantryConfig = unitConfig{
	name:                 infantry,
	radius:               3,
	maxHealth:            50,
	attack:               1,
	vision:               30,
	reach:                5,
	produceLayer:         land,
	staticTarget:         false,
	attacksGround:        true,
	attacksLowAltitude:   false,
	attacksHighAltitude:  false,
	explodeAtTerminalPos: false,
	despawnAtTerminalPos: false,
	explodeRadius:        0,
	maxAlivePerNode:      math.MaxInt,
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
	staticTarget:         true,
	attacksGround:        false,
	attacksLowAltitude:   false,
	attacksHighAltitude:  false,
	explodeAtTerminalPos: true,
	despawnAtTerminalPos: true,
	explodeRadius:        20,
	maxAlivePerNode:      math.MaxInt,
}

var guidedAirMissileConfig = unitConfig{
	name:                 guidedAirMissile,
	radius:               5,
	maxHealth:            100,
	attack:               1000,
	vision:               40,
	reach:                0,
	produceLayer:         lowAir,
	staticTarget:         false,
	attacksGround:        false,
	attacksLowAltitude:   true,
	attacksHighAltitude:  true,
	explodeAtTerminalPos: false,
	despawnAtTerminalPos: true,
	explodeRadius:        100,
	maxAlivePerNode:      math.MaxInt,
}

var bomberConfig = unitConfig{
	name:                 bomber,
	radius:               20,
	maxHealth:            100,
	attack:               0,
	vision:               40,
	reach:                0,
	produceLayer:         highAir,
	staticTarget:         false,
	attacksGround:        false,
	attacksLowAltitude:   false,
	attacksHighAltitude:  false,
	explodeAtTerminalPos: false,
	despawnAtTerminalPos: false,
	explodeRadius:        100,
	maxAlivePerNode:      1,
}

var bombConfig = unitConfig{
	name:                 bomb,
	radius:               4,
	maxHealth:            100,
	attack:               0,
	vision:               0,
	reach:                0,
	produceLayer:         lowAir,
	staticTarget:         true,
	attacksGround:        true,
	attacksLowAltitude:   false,
	attacksHighAltitude:  false,
	explodeAtTerminalPos: false,
	despawnAtTerminalPos: false,
	explodeRadius:        20,
	explodeStepDelay:     300,
	maxAlivePerNode:      math.MaxInt,
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
	a = append(a, slog.Bool("attacksGround", uc.attacksGround))
	a = append(a, slog.Bool("attacksLowAltitude", uc.attacksLowAltitude))
	a = append(a, slog.Bool("attacksHighAltitude", uc.attacksHighAltitude))
	a = append(a, slog.Bool("explodeAtTerminalPos", uc.explodeAtTerminalPos))
	a = append(a, slog.Bool("despawnAtTerminalPos", uc.despawnAtTerminalPos))
	return a
}
