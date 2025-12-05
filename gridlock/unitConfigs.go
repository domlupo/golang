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
	uc.setRadius()
	uc.setMaxHealth()
	uc.setAttack()
	uc.setVision()
	uc.setReach()
	uc.setExplodeRadius()
	uc.setMaxAlivePerNode()
	uc.setSpeedFactor()
	uc.setSurvivesCollision()
	uc.setTracks()
	uc.setTracksLayers()
	uc.setConsumables()
	return uc
}

func (uc *unitConfig) setRadius() {
	var radius int
	switch uc.name {
	case infantry, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare:
		radius = 3
	case bombGravity:
		radius = 4
	case fighter:
		radius = 8
	case bomber:
		radius = 20
	case chaff:
		radius = 100
	default:
		panic(fmt.Sprintf("setRadius: cannot set radius for unknown unit %s", uc.name.string()))
	}
	uc.radius = radius
}

func (uc *unitConfig) setMaxHealth() {
	var maxHealth int
	switch uc.name {
	case infantry:
		maxHealth = 50
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare,
		bombGravity, fighter, bomber, chaff:
		maxHealth = 1
	default:
		panic(fmt.Sprintf("setMaxHealth: cannot set max health for unknown unit %s", uc.name.string()))
	}
	uc.maxHealth = maxHealth
}

func (uc *unitConfig) setAttack() {
	var attack int
	switch uc.name {
	case infantry:
		attack = 1
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare,
		bombGravity, fighter, bomber, chaff:
	default:
		panic(fmt.Sprintf("setAttack: cannot set attack for unknown unit %s", uc.name.string()))
	}
	uc.attack = attack
}

func (uc *unitConfig) setVision() {
	var vision int
	switch uc.name {
	case infantry:
		vision = 20 - uc.radius
		if vision >= 20 {
			panic("setVision: temporary protection for partitioner.go infantry limitation. It needs to support higher order hexagon LUT")
		}
	case fighter:
		vision = 800
	case bomber:
		vision = 200
	case missileRadarActive, missileHeatSeeking, missileAntiBallistic:
		vision = 300
	case chaff:
		vision = 100
	case bombGravity, missileCruise, flare, missileBallistic:
		vision = uc.radius
	default:
		panic(fmt.Sprintf("setVision: cannot set vision for unknown unit %s", uc.name.string()))
	}
	uc.vision = vision
	if uc.vision < uc.radius {
		panic("setVision: unit radius and vision are both relative to center. vision < radius is invalid.")
	}
}

func (uc *unitConfig) setReach() {
	var reach int
	switch uc.name {
	case infantry:
		reach = 5
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare,
		bombGravity, fighter, bomber, chaff:
		reach = uc.radius
	default:
		panic(fmt.Sprintf("setReach: cannot set reach for unknown unit %s", uc.name.string()))
	}
	uc.reach = reach
	if uc.reach < uc.radius {
		panic("setReach: unit radius and reach are both relative to center. reach < radius is invalid.")
	}
	uc.reach = reach
}

func (uc *unitConfig) setExplodeRadius() {
	var explodeRadius int
	switch uc.name {
	case infantry, flare, chaff:
		explodeRadius = 0
	case fighter, bomber:
		explodeRadius = 100
	case bombGravity, missileCruise, missileBallistic, missileRadarActive, missileHeatSeeking, missileAntiBallistic:
		explodeRadius = 20
	default:
		panic(fmt.Sprintf("setExplodeRadius: cannot set explode radius unknown unit %s", uc.name.string()))
	}
	uc.explodeRadius = explodeRadius
}

func (uc *unitConfig) setMaxAlivePerNode() {
	var maxAlivePerNode int
	switch uc.name {
	case infantry, flare, chaff, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, bombGravity:
		maxAlivePerNode = math.MaxInt
	case fighter, bomber:
		maxAlivePerNode = 1
	default:
		panic(fmt.Sprintf("setMaxAlivePerNode: cannot set max alive per node for unknown unit %s", uc.name.string()))
	}
	uc.maxAlivePerNode = maxAlivePerNode
}

func (uc *unitConfig) setSpeedFactor() {
	var speedFactor float64
	switch uc.name {
	case infantry:
		speedFactor = 2.0
	case fighter, bomber, missileCruise, missileBallistic, missileAntiBallistic:
		speedFactor = 1.0
	case missileRadarActive, missileHeatSeeking, flare:
		speedFactor = 0.5
	case chaff, bombGravity:
		speedFactor = 0
	default:
		panic(fmt.Sprintf("setSpeedFactor: cannot set speed factor for unknown unit %s", uc.name.string()))
	}
	uc.speedFactor = speedFactor
}

func (uc *unitConfig) setSurvivesCollision() {
	var survivesCollision bool
	switch uc.name {
	case infantry, chaff:
		survivesCollision = true
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare,
		bombGravity, fighter, bomber:
	default:
		panic(fmt.Sprintf("setSurvivesCollision: cannot set survives collision for unknown unit %s", uc.name.string()))
	}
	uc.survivesCollision = survivesCollision
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
		panic(fmt.Sprintf("setTracks : cannot set tracks for unknown unit %s", uc.name.string()))
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
	name:            infantry,
	produceLayer:    lLand,
	collisionLayers: []layer{lLand},
}

var fighterUnitConfig = unitConfig{
	name:            fighter,
	produceLayer:    lAircraft,
	collisionLayers: []layer{lAircraft},
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
	produceLayer:    lAircraft,
	collisionLayers: []layer{lAircraft},
	emitsHeat:       true,
	emitsRadar:      true,
	reactive:        true,
	requester:       true,
	ignore:          func(u, t *unit) bool { return slices.Contains(u.skipTargets, t.id) },
}

var bombGravityUnitConfig = unitConfig{
	name:                bombGravity,
	produceLayer:        lBombGravity,
	collisionLayers:     []layer{lBombGravity},
	staticTarget:        true,
	terminalStepDelay:   300,
	explodeWhenTerminal: true,
}

var missileRadarActiveUnitConfig = unitConfig{
	name:                 missileRadarActive,
	produceLayer:         lOrdnance,
	collisionLayers:      []layer{lOrdnance, lAircraft, lChaff},
	staticTarget:         true,
	despawnWhenTerminal:  true,
	ignoreTarget:         func(t *unit) bool { return t.name == flare },
	terminalWhenMoveDone: true,
	radarGuided:          true,
}

var missileHeatSeekingUnitConfig = unitConfig{
	name:                 missileHeatSeeking,
	produceLayer:         lOrdnance,
	collisionLayers:      []layer{lOrdnance, lAircraft},
	staticTarget:         true,
	despawnWhenTerminal:  true,
	heatSeeking:          true,
	ignoreTarget:         func(t *unit) bool { return !t.emitsHeat },
	terminalWhenMoveDone: true,
}

var missileCruiseUnitConfig = unitConfig{
	name:                 missileCruise,
	produceLayer:         lOrdnance,
	collisionLayers:      layers{lOrdnance},
	staticTarget:         true,
	explodeWhenTerminal:  true,
	terminalWhenMoveDone: true,
}

var missileBallisticUnitConfig = unitConfig{
	name:                 missileBallistic,
	produceLayer:         lBallistic,
	collisionLayers:      []layer{lBallistic},
	staticTarget:         true,
	explodeWhenTerminal:  true,
	terminalWhenMoveDone: true,
}

var missileAntiBallisticUnitConfig = unitConfig{
	name:                 missileAntiBallistic,
	produceLayer:         lBallistic,
	collisionLayers:      []layer{lBallistic},
	staticTarget:         true,
	despawnWhenTerminal:  true,
	terminalWhenMoveDone: true,
}

var flareUnitConfig = unitConfig{
	name:                 flare,
	produceLayer:         lOrdnance,
	collisionLayers:      []layer{lOrdnance},
	staticTarget:         true,
	despawnWhenTerminal:  true,
	emitsHeat:            true,
	terminalWhenMoveDone: true,
}

var chaffUnitConfig = unitConfig{
	name:                 chaff,
	produceLayer:         lChaff,
	collisionLayers:      []layer{lOrdnance},
	staticTarget:         true,
	terminalStepDelay:    300,
	despawnWhenTerminal:  true,
	ignoreCollision:      func(other *unit) bool { return !other.radarGuided },
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
