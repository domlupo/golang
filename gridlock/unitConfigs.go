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

var unitConfigs map[name]unitConfig = map[name]unitConfig{
	infantry:             infantryUnitConfig,
	fighter:              fighterUnitConfig,
	bomber:               bomberUnitConfig,
	bombGravity:          bombGravityUnitConfig,
	missileCruise:        missileCruiseUnitConfig,
	missileHeatSeeking:   missileHeatSeekingUnitConfig,
	missileRadarActive:   missileRadarActiveUnitConfig,
	missileBallistic:     missileBallisticUnitConfig,
	missileAntiBallistic: missileAntiBallisticUnitConfig,
	flare:                flareUnitConfig,
	chaff:                chaffUnitConfig,
}

type consumable struct {
	count         int
	countMax      int
	fireDelay     step
	firesAt       func(uc unitConfig) bool
	firesAtLayers layers
}

func newConsumable(name name, countMax int) consumable {
	c := consumable{
		countMax: countMax,
		count:    countMax,
	}
	c.setFireDelay(name)
	c.setFiresAt(name)
	c.setFiresAtLayers()
	return c
}

func (c *consumable) setFireDelay(name name) {
	var fireDelay step
	switch name {
	case missileHeatSeeking:
		fireDelay = 50
	case missileRadarActive:
		fireDelay = 50
	case flare:
		fireDelay = 50
	case chaff:
		fireDelay = 50
	case bombGravity:
	default:
		panic(fmt.Sprintf("setFireDelay: cannot set fire delay for unknown consumable %s", name.string()))
	}
	c.fireDelay = fireDelay
}

func (c *consumable) setFiresAt(name name) {
	firesAt := func(uc unitConfig) bool { return false }
	switch name {
	case missileHeatSeeking:
		firesAt = func(uc unitConfig) bool { return uc.emitsHeat && uc.name != flare }
	case missileRadarActive:
		firesAt = func(uc unitConfig) bool { return uc.emitsRadar && uc.name != chaff }
	case bombGravity, flare, chaff:
	default:
		panic(fmt.Sprintf("setFiresAt: cannot set fires at for unknown consumable %s", name.string()))
	}
	c.firesAt = firesAt
}

func (c *consumable) setFiresAtLayers() {
	layers := []layer{}
	for _, other := range unitConfigs {
		if c.firesAt(other) {
			layers = append(layers, other.produceLayer)
		}
	}
	c.firesAtLayers = layers
}

func (c consumable) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("count", c.count))
	a = append(a, slog.Int("countMax", c.countMax))
	a = append(a, slog.Uint64("fireDelay", c.fireDelay))
	a = append(a, slog.String("firesAt", c.firesAtLayers.string()))
	return a
}

type unitConfig struct {
	name                 name
	radius               int
	vision               int
	maxHealth            int
	attack               int
	reach                int
	produceLayer         layer
	staticTarget         bool
	tracks               func(t unitConfig) bool
	trackLayers          layers
	collisionLayers      layers
	explodeRadius        int
	terminalWhenMoveDone bool
	terminalStepDelay    step
	terminalStep         step
	explodeWhenTerminal  bool
	despawnWhenTerminal  bool
	reloadPayloadDelay   uint64
	maxAlivePerNode      int
	speedFactor          float64
	reactions            []reaction
	reactive             bool
	requests             []request
	requester            bool
	consumables          map[name]consumable
	emitsHeat            bool
	emitsRadar           bool
	heatSeeking          bool
	radarGuided          bool
	ignore               func(u, t *unit) bool
	ignoreTarget         func(t *unit) bool
	ignoreCollision      func(other *unit) bool
	survivesCollision    bool
	// This flag should rarely be used. If it is used often then write a function to avoid bugs.
	// This is safe to enable iff the unit only collides with units that have radiuses less
	// than state.spawnUnitRadiusLimit. For example, chaff only collides with radarGuided units
	// which are all less than the 20 pixels.
	skipSpawnRadiusCheck bool
}

func newUnitConfig(name name) unitConfig {
	uc, ok := unitConfigs[name]
	if !ok {
		panic(fmt.Sprintf("newUnitConfig: unknown unit name %s", name.string()))
	}
	uc.setTracks()
	uc.setTracksLayers()
	uc.setConsumables()
	return uc
}

func (uc *unitConfig) setTracks() {
	tracks := func(uc unitConfig) bool { return false }
	switch uc.name {
	case infantry:
		tracks = func(uc unitConfig) bool { return uc.name == infantry }
	case missileHeatSeeking:
		tracks = func(uc unitConfig) bool { return uc.emitsHeat }
	case missileRadarActive:
		tracks = func(uc unitConfig) bool { return uc.emitsRadar }
	case missileAntiBallistic:
		tracks = func(uc unitConfig) bool { return uc.name == missileBallistic }
	case fighter, bomber, bombGravity, missileCruise, missileBallistic, flare, chaff:
	default:
		panic(fmt.Sprintf("setTracks : cannot set tracks for unknown consumable %s", uc.name.string()))
	}
	uc.tracks = tracks
}

func (uc *unitConfig) setTracksLayers() {
	layers := []layer{}
	if uc.tracks != nil {
		for _, other := range unitConfigs {
			if uc.tracks(other) {
				layers = append(layers, other.produceLayer)
			}
		}
	}
	uc.trackLayers = layers
}

func (uc *unitConfig) setConsumables() {
	switch uc.name {
	case fighter:
		uc.consumables = map[name]consumable{
			bombGravity:        newConsumable(bombGravity, 1),
			flare:              newConsumable(flare, 1000),
			chaff:              newConsumable(chaff, 0),
			missileHeatSeeking: newConsumable(missileHeatSeeking, 10),
			missileRadarActive: newConsumable(missileRadarActive, 0),
		}
	case bomber:
		uc.consumables = map[name]consumable{
			bombGravity: newConsumable(bombGravity, 4),
			flare:       newConsumable(flare, 1000),
			chaff:       newConsumable(chaff, 1000),
		}
	case infantry, missileCruise, missileBallistic, missileAntiBallistic, missileHeatSeeking, missileRadarActive:
	default:
		panic(fmt.Sprintf("setConsumables: cannot set consumables for unknown unit %s", uc.name.string()))
	}
}

// Note: unit radius set to an odd number is better for testing due to division bugs
var infantryUnitConfig = unitConfig{
	name:              infantry,
	radius:            3,
	maxHealth:         50,
	attack:            1,
	vision:            17, // This must be 20 - infantry.radius until partitioner.go supports higher order hexagon LUT
	reach:             5,
	produceLayer:      lLand,
	collisionLayers:   []layer{lLand},
	explodeRadius:     0,
	maxAlivePerNode:   math.MaxInt,
	speedFactor:       2.0,
	survivesCollision: true,
}

var fighterUnitConfig = unitConfig{
	name:            fighter,
	radius:          8,
	maxHealth:       1,
	vision:          800,
	produceLayer:    lAircraft,
	collisionLayers: []layer{lAircraft},
	explodeRadius:   100,
	maxAlivePerNode: 1,
	speedFactor:     1.0,
	reactive:        true,
	requester:       true,
	emitsHeat:       true,
	emitsRadar:      true,
	ignoreTarget: func(t *unit) bool {
		return t.name == flare
	},
}

var bomberUnitConfig = unitConfig{
	name:            bomber,
	radius:          20,
	maxHealth:       1,
	vision:          200,
	produceLayer:    lAircraft,
	collisionLayers: []layer{lAircraft},
	explodeRadius:   100,
	maxAlivePerNode: 1,
	speedFactor:     1.0,
	emitsHeat:       true,
	emitsRadar:      true,
	reactive:        true,
	requester:       true,
	ignore:          func(u, t *unit) bool { return slices.Contains(u.skipTargets, t.id) },
}

var bombGravityUnitConfig = unitConfig{
	name:                bombGravity,
	radius:              4,
	vision:              4,
	maxHealth:           1,
	produceLayer:        lBombGravity,
	collisionLayers:     []layer{lBombGravity},
	staticTarget:        true,
	explodeRadius:       20,
	terminalStepDelay:   300,
	maxAlivePerNode:     math.MaxInt,
	explodeWhenTerminal: true,
}

var missileRadarActiveUnitConfig = unitConfig{
	name:                 missileRadarActive,
	radius:               3,
	maxHealth:            1,
	vision:               300,
	produceLayer:         lOrdnance,
	collisionLayers:      []layer{lOrdnance, lAircraft, lChaff},
	staticTarget:         true,
	despawnWhenTerminal:  true,
	speedFactor:          0.5,
	maxAlivePerNode:      math.MaxInt,
	ignoreTarget:         func(t *unit) bool { return t.name == flare },
	terminalWhenMoveDone: true,
	radarGuided:          true,
}

var missileHeatSeekingUnitConfig = unitConfig{
	name:                 missileHeatSeeking,
	radius:               3,
	maxHealth:            1,
	vision:               300,
	produceLayer:         lOrdnance,
	collisionLayers:      []layer{lOrdnance, lAircraft},
	staticTarget:         true,
	despawnWhenTerminal:  true,
	heatSeeking:          true,
	speedFactor:          0.5,
	maxAlivePerNode:      math.MaxInt,
	ignoreTarget:         func(t *unit) bool { return !t.emitsHeat },
	terminalWhenMoveDone: true,
}

var missileCruiseUnitConfig = unitConfig{
	name:                 missileCruise,
	radius:               3,
	vision:               3,
	maxHealth:            1,
	produceLayer:         lOrdnance,
	collisionLayers:      layers{lOrdnance},
	staticTarget:         true,
	explodeWhenTerminal:  true,
	explodeRadius:        20,
	speedFactor:          1,
	maxAlivePerNode:      math.MaxInt,
	terminalWhenMoveDone: true,
}

var missileBallisticUnitConfig = unitConfig{
	name:                 missileBallistic,
	radius:               3,
	vision:               3,
	maxHealth:            1,
	produceLayer:         lBallistic,
	collisionLayers:      []layer{lBallistic},
	staticTarget:         true,
	explodeWhenTerminal:  true,
	explodeRadius:        20,
	speedFactor:          1,
	maxAlivePerNode:      math.MaxInt,
	terminalWhenMoveDone: true,
}

var missileAntiBallisticUnitConfig = unitConfig{
	name:                 missileAntiBallistic,
	radius:               3,
	maxHealth:            1,
	vision:               300,
	produceLayer:         lBallistic,
	collisionLayers:      []layer{lBallistic},
	staticTarget:         true,
	despawnWhenTerminal:  true,
	explodeRadius:        20,
	speedFactor:          1,
	maxAlivePerNode:      math.MaxInt,
	terminalWhenMoveDone: true,
}

var flareUnitConfig = unitConfig{
	name:                 flare,
	radius:               3,
	vision:               3,
	maxHealth:            1,
	produceLayer:         lOrdnance,
	collisionLayers:      []layer{lOrdnance},
	staticTarget:         true,
	despawnWhenTerminal:  true,
	speedFactor:          0.5,
	emitsHeat:            true,
	terminalWhenMoveDone: true,
}

var chaffUnitConfig = unitConfig{
	name:                 chaff,
	radius:               100,
	vision:               100,
	maxHealth:            1,
	produceLayer:         lChaff,
	collisionLayers:      []layer{lOrdnance},
	staticTarget:         true,
	maxAlivePerNode:      math.MaxInt,
	terminalStepDelay:    300,
	despawnWhenTerminal:  true,
	ignoreCollision:      func(other *unit) bool { return !other.radarGuided },
	survivesCollision:    true,
	skipSpawnRadiusCheck: true,
}

// This log method is frequently used for debugging. Also it is typically logged at the end of unit.log().
// Log frequently used info at the end for quicker debugging. Less used info can be parsed via tooling e.g. grep
func (uc unitConfig) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("radius", uc.radius))
	a = append(a, slog.Int("maxHealth", uc.maxHealth))
	a = append(a, slog.Int("attack", uc.attack))
	a = append(a, slog.Int("vision", uc.vision))
	a = append(a, slog.Int("reach", uc.reach))
	a = append(a, slog.Bool("staticTarget", uc.staticTarget))
	a = append(a, slog.Bool("heatSeeking", uc.heatSeeking))
	a = append(a, slog.Bool("explodeWhenTerminal", uc.explodeWhenTerminal))
	a = append(a, slog.Bool("despawnWhenTerminal", uc.despawnWhenTerminal))
	a = append(a, slog.Uint64("terminalStepDelay", uc.terminalStepDelay))
	a = append(a, slog.Uint64("terminalStep", uc.terminalStep))
	a = append(a, slog.String("collisionLayers", uc.collisionLayers.string()))
	a = append(a, uc.produceLayer.log())
	for k, v := range uc.consumables {
		a = append(a, nest(k.string(), v.log()))
	}
	a = append(a, slog.String("name", uc.name.string()))
	return a
}
