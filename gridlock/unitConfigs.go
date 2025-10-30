// Configurations for different unit types.

package main

import "log/slog"

type name int

const (
	noName name = iota
	infantry
	unguidedGroundMissile
	guidedAirMissile
	mobileSAM
)

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
	attacksPosition      bool
	explodeAtTerminalPos bool
	explodeRadius        int
}

// Note: unit radius set to an odd number is better for testing due to division bugs
var infantryConfig = unitConfig{
	name:                 infantry,
	radius:               3,
	maxHealth:            50,
	attack:               1,
	vision:               30,
	reach:                5,
	produceLayer:         ground,
	staticTarget:         false,
	attacksGround:        true,
	attacksLowAltitude:   false,
	attacksHighAltitude:  false,
	attacksPosition:      false,
	explodeAtTerminalPos: false,
	explodeRadius:        0,
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
	produceLayer:         lowAltitude,
	staticTarget:         true,
	attacksGround:        false,
	attacksLowAltitude:   false,
	attacksHighAltitude:  false,
	attacksPosition:      true,
	explodeAtTerminalPos: true,
	explodeRadius:        20,
}

var guidedAirMissileConfig = unitConfig{
	name:                 guidedAirMissile,
	radius:               5,
	maxHealth:            100,
	attack:               1000,
	vision:               40,
	reach:                0,
	produceLayer:         lowAltitude,
	staticTarget:         false,
	attacksGround:        false,
	attacksLowAltitude:   true,
	attacksHighAltitude:  true,
	attacksPosition:      false,
	explodeAtTerminalPos: true,
	explodeRadius:        100,
}

func (uc unitConfig) log() []slog.Attr {
	var a = []slog.Attr{}
	var name string
	switch uc.name {
	case noName:
		name = "no name"
	case infantry:
		name = "infantry"
	case unguidedGroundMissile:
		name = "unguidedGroundMissile"
	case guidedAirMissile:
		name = "guidedAirMissile"
	default:
		name = "unknown"
	}
	a = append(a, slog.String("name", name))
	a = append(a, slog.Int("radius", uc.radius))
	a = append(a, slog.Int("maxHealth", uc.maxHealth))
	a = append(a, slog.Int("attack", uc.attack))
	a = append(a, slog.Int("vision", uc.vision))
	a = append(a, slog.Int("reach", uc.reach))
	var layer string
	switch uc.produceLayer {
	case noLayer:
		layer = "noLayer"
	case ground:
		layer = "ground"
	case lowAltitude:
		layer = "lowAltitude"
	case highAltitude:
		layer = "highAltitude"
	}
	a = append(a, slog.String("produceLayer", layer))
	a = append(a, slog.Bool("staticTarget", uc.staticTarget))
	a = append(a, slog.Bool("attacksGround", uc.attacksGround))
	a = append(a, slog.Bool("attacksLowAltitude", uc.attacksLowAltitude))
	a = append(a, slog.Bool("attacksHighAltitude", uc.attacksHighAltitude))
	a = append(a, slog.Bool("attacksPosition", uc.attacksPosition))
	return a
}
