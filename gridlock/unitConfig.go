// Refer to unitConfigProvider.go file for documentation. Code in separate file for readability.

package main

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
)

type unitConfig struct {
	name                   name
	radius                 int
	vision                 int
	maxHealth              int
	attack                 int
	reach                  int
	produceLayer           layer
	tracks                 func(t unitConfig) bool
	trackLayers            layers
	collisionLayers        layers
	explodeRadius          int
	terminalStepDelay      step
	terminalStep           step
	reloadPayloadDelay     uint64
	maxAlivePerNode        int
	speedFactor            float64
	reactions              []reaction
	requests               []request
	equipments             []equipment
	consumables            consumables
	isReactive             bool
	isEmitsHeat            bool
	isEmitsRadar           bool
	isHeatSeeking          bool
	isRadarGuided          bool
	isRequester            bool
	isSurvivesCollision    bool
	isExplodeWhenTerminal  bool
	isDespawnWhenTerminal  bool
	isTerminalWhenMoveDone bool
	isStaticTarget         bool
	ignore                 func(u, t *unit) bool
	ignoreTarget           func(t *unit) bool
	ignoreCollision        func(other *unit) bool
	isSkipSpawnRadiusCheck bool
	// isSkipSpawnRadiusCheck should rarely be used. If it is used often then write a function to
	// avoid bugs. This is safe to enable iff the unit only collides with units that have radiuses
	// less than state.spawnUnitRadiusLimit. For example, chaff only collides with radarGuided units
	// which are all less than the 20 pixels.
}

func (uc *unitConfig) setUCFields() {
	// TODO: Once unit is more finalized add set methods for these fields.
	// I suspect these might change in the near future so I am punting them for a bit
	switch uc.name {
	case fighter:
		uc.ignoreTarget = func(t *unit) bool { return t.name == flare }
	case bomber:
		uc.ignore = func(u, t *unit) bool { return slices.Contains(u.skipTargets, t.id) }
	case missileRadarActive:
		uc.ignoreTarget = func(t *unit) bool { return t.name == flare }
	case missileHeatSeeking:
		uc.ignoreTarget = func(t *unit) bool { return !t.isEmitsHeat }
	case chaff:
		uc.ignoreCollision = func(other *unit) bool { return !other.isRadarGuided }
	}
	uc.setRadius()
	uc.setMaxHealth()
	uc.setAttack()
	uc.setVision()
	uc.setReach()
	uc.setExplodeRadius()
	uc.setMaxAlivePerNode()
	uc.setSpeedFactor()
	uc.setIsSurvivesCollision()
	uc.setIsEmitsHeat()
	uc.setIsEmitsRadar()
	uc.setIsRequester()
	uc.setIsReactive()
	uc.setTerminalStepDelay()
	uc.setIsExplodeWhenTerminal()
	uc.setIsStaticTarget()
	uc.setIsDespawnWhenTerminal()
	uc.setIsTerminalWhenMoveDone()
	uc.setIsRadarGuided()
	uc.setIsHeatSeeking()
	uc.setIsSkipSpawnRadiusCheck()
	uc.setTracks()
	uc.setProduceLayer()
	uc.setCollisionLayers()
	uc.setEquipments()
}

func (uc *unitConfig) setRadius() {
	var radius int
	switch uc.name {
	case infantry, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, flare:
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
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, flare,
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
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, flare, bombGravity, fighter, bomber, chaff:
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
			panic("setVision: temporary protection for partitioner.go infantry limitation. It needs to " +
				"support higher order hexagon LUT")
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
		panic(
			"setVision: unit radius and vision are both relative to center. vision < radius is invalid.")
	}
}

func (uc *unitConfig) setReach() {
	var reach int
	switch uc.name {
	case infantry:
		reach = 5
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, flare, bombGravity, fighter, bomber, chaff:
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
	case bombGravity, missileCruise, missileBallistic, missileRadarActive, missileHeatSeeking,
		missileAntiBallistic:
		explodeRadius = 20
	default:
		panic(
			fmt.Sprintf("setExplodeRadius: cannot set explode radius unknown unit %s", uc.name.string()))
	}
	uc.explodeRadius = explodeRadius
}

func (uc *unitConfig) setMaxAlivePerNode() {
	var maxAlivePerNode int
	switch uc.name {
	case infantry, flare, chaff, missileRadarActive, missileHeatSeeking, missileCruise,
		missileBallistic, missileAntiBallistic, bombGravity:
		maxAlivePerNode = math.MaxInt
	case fighter, bomber:
		maxAlivePerNode = 1
	default:
		panic(fmt.Sprintf("setMaxAlivePerNode: cannot set max alive per node for unknown unit %s",
			uc.name.string()))
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
		panic(fmt.Sprintf("setSpeedFactor: cannot set speed factor for unknown unit %s",
			uc.name.string()))
	}
	uc.speedFactor = speedFactor
}

func (uc *unitConfig) setIsSurvivesCollision() {
	var survivesCollision bool
	switch uc.name {
	case infantry, chaff:
		survivesCollision = true
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, flare, bombGravity, fighter, bomber:
	default:
		panic(fmt.Sprintf("setIsSurvivesCollision: cannot set isSurvivesCollision for unknown unit %s",
			uc.name.string()))
	}
	uc.isSurvivesCollision = survivesCollision
}

func (uc *unitConfig) setIsEmitsHeat() {
	var isEmitsHeat bool
	switch uc.name {
	case fighter, bomber, flare:
		isEmitsHeat = true
	case infantry, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, bombGravity, chaff:
	default:
		panic(fmt.Sprintf("setIsEmitsHeat: cannot set isEmitsHeats for unknown unit %s",
			uc.name.string()))
	}
	uc.isEmitsHeat = isEmitsHeat
}

func (uc *unitConfig) setIsEmitsRadar() {
	var isEmitsRadar bool
	switch uc.name {
	case fighter, bomber:
		isEmitsRadar = true
	case infantry, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, flare, bombGravity, chaff:
	default:
		panic(fmt.Sprintf("setIsEmitsRadar: cannot set isEmitsRadar for unknown unit %s",
			uc.name.string()))
	}
	uc.isEmitsRadar = isEmitsRadar
}

func (uc *unitConfig) setIsRequester() {
	var isRequester bool
	switch uc.name {
	case fighter, bomber:
		isRequester = true
	case infantry, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, flare, bombGravity, chaff:
	default:
		panic(fmt.Sprintf("setIsRequester: cannot set isRequester for unknown unit %s",
			uc.name.string()))
	}
	uc.isRequester = isRequester
}

func (uc *unitConfig) setIsReactive() {
	var isReactive bool
	switch uc.name {
	case fighter, bomber:
		isReactive = true
	case infantry, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, flare, bombGravity, chaff:
	default:
		panic(fmt.Sprintf("setIsReactive: cannot set isReactive for unknown unit %s",
			uc.name.string()))
	}
	uc.isReactive = isReactive
}

func (uc *unitConfig) setProduceLayer() {
	var produceLayer layer
	switch uc.name {
	case infantry:
		produceLayer = lLand
	case fighter, bomber:
		produceLayer = lAircraft
	case missileRadarActive, missileHeatSeeking, missileCruise, flare:
		produceLayer = lOrdnance
	case missileBallistic, missileAntiBallistic:
		produceLayer = lBallistic
	case chaff:
		produceLayer = lChaff
	case bombGravity:
		produceLayer = lBombGravity
	default:
		panic(fmt.Sprintf("setProduceLayer: cannot set produceLayer for unknown unit %s",
			uc.name.string()))
	}
	uc.produceLayer = produceLayer
}

func (uc *unitConfig) setTerminalStepDelay() {
	var terminalStepDelay step
	switch uc.name {
	case bombGravity, chaff:
		terminalStepDelay = 300
	case infantry, fighter, bomber, missileRadarActive, missileHeatSeeking, missileCruise,
		missileBallistic, missileAntiBallistic, flare:
	default:
		panic(fmt.Sprintf("setTerminalStepDelay: cannot set terminal step delay for unknown unit %s",
			uc.name.string()))
	}
	uc.terminalStepDelay = terminalStepDelay
}

func (uc *unitConfig) setIsExplodeWhenTerminal() {
	var isExplodeWhenTerminal bool
	switch uc.name {
	case bombGravity, missileCruise, missileBallistic:
		isExplodeWhenTerminal = true
	case infantry, fighter, bomber, missileRadarActive, missileHeatSeeking, missileAntiBallistic,
		flare, chaff:
	default:
		panic(
			fmt.Sprintf(
				"setIsExplodeWhenTerminal: cannot set is explode when terminal for unknown unit %s",
				uc.name.string()))
	}
	uc.isExplodeWhenTerminal = isExplodeWhenTerminal
}

func (uc *unitConfig) setIsStaticTarget() {
	var isStatictarget bool
	switch uc.name {
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, flare, bombGravity, chaff:
		isStatictarget = true
	case infantry, fighter, bomber:
	default:
		panic(fmt.Sprintf("setIsStaticTarget: cannot set isStaticTarget for unknown unit %s",
			uc.name.string()))
	}
	uc.isStaticTarget = isStatictarget
}

func (uc *unitConfig) setIsDespawnWhenTerminal() {
	var isDespawnWhenTerminal bool
	switch uc.name {
	case missileRadarActive, missileHeatSeeking, missileAntiBallistic, flare, chaff:
		isDespawnWhenTerminal = true
	case infantry, fighter, bomber, missileCruise, missileBallistic, bombGravity:
	default:
		panic(fmt.Sprintf(
			"setIsDespawnWhenTerminal: cannot set isDespawnWhenTerminal for unknown unit %s",
			uc.name.string()))
	}
	uc.isDespawnWhenTerminal = isDespawnWhenTerminal
}

func (uc *unitConfig) setIsTerminalWhenMoveDone() {
	var isTerminalWhenMoveDone bool
	switch uc.name {
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, flare:
		isTerminalWhenMoveDone = true
	case infantry, fighter, bomber, bombGravity, chaff:
	default:
		panic(fmt.Sprintf(
			"setIsTerminalWhenMoveDone: cannot set isTerminalWhenMoveDone for unknown unit %s",
			uc.name.string()))
	}
	uc.isTerminalWhenMoveDone = isTerminalWhenMoveDone
}

func (uc *unitConfig) setIsRadarGuided() {
	var isRadarGuided bool
	switch uc.name {
	case missileRadarActive:
		isRadarGuided = true
	case infantry, fighter, bomber, bombGravity, chaff, missileHeatSeeking, missileCruise,
		missileBallistic, missileAntiBallistic, flare:
	default:
		panic(fmt.Sprintf("setIsRadarGuided: cannot set isRadarGuided for unknown unit %s",
			uc.name.string()))
	}
	uc.isRadarGuided = isRadarGuided
}

func (uc *unitConfig) setIsHeatSeeking() {
	var isHeatSeeking bool
	switch uc.name {
	case missileHeatSeeking:
		isHeatSeeking = true
	case infantry, fighter, bomber, bombGravity, chaff, missileRadarActive, missileCruise,
		missileBallistic, missileAntiBallistic, flare:
	default:
		panic(fmt.Sprintf("setIsHeatSeeking: cannot set isHeatSeeking for unknown unit %s",
			uc.name.string()))
	}
	uc.isHeatSeeking = isHeatSeeking
}

func (uc *unitConfig) setIsSkipSpawnRadiusCheck() {
	var isSkipSpawnRadiusCheck bool
	switch uc.name {
	case chaff:
		isSkipSpawnRadiusCheck = true
	case missileHeatSeeking, infantry, fighter, bomber, bombGravity, missileRadarActive,
		missileCruise, missileBallistic, missileAntiBallistic, flare:
	default:
		panic(fmt.Sprintf(
			"setIsSkipSpawnRadiusCheck: cannot set isSkipSpawnRadiusCheck for unknown unit %s",
			uc.name.string()))
	}
	uc.isSkipSpawnRadiusCheck = isSkipSpawnRadiusCheck
}

func (uc *unitConfig) setCollisionLayers() {
	var collisionLayers layers
	switch uc.name {
	case infantry, fighter, bomber, bombGravity, missileCruise, missileBallistic,
		missileAntiBallistic, flare:
		collisionLayers = []layer{uc.produceLayer}
	case missileRadarActive:
		collisionLayers = []layer{lOrdnance, lAircraft, lChaff}
	case missileHeatSeeking:
		collisionLayers = []layer{lOrdnance, lAircraft}
	case chaff:
		collisionLayers = []layer{lOrdnance}
	default:
		panic(fmt.Sprintf("setCollisionLayers: cannot set collision layers for unknown unit %s",
			uc.name.string()))
	}
	uc.collisionLayers = collisionLayers
}

func (uc *unitConfig) setTracks() {
	tracks := func(uc unitConfig) bool { return false }
	switch uc.name {
	case infantry:
		tracks = func(uc unitConfig) bool { return uc.name == infantry }
	case missileHeatSeeking:
		tracks = func(uc unitConfig) bool { return uc.isEmitsHeat }
	case missileRadarActive:
		tracks = func(uc unitConfig) bool { return uc.isEmitsRadar }
	case missileAntiBallistic:
		tracks = func(uc unitConfig) bool { return uc.name == missileBallistic }
	case fighter, bomber, bombGravity, missileCruise, missileBallistic, flare, chaff:
	default:
		panic(fmt.Sprintf("setTracks : cannot set tracks for unknown unit %s", uc.name.string()))
	}
	uc.tracks = tracks
}

func (uc *unitConfig) setEquipments() {
	switch uc.name {

	// TODO: config these hard coded values. It will be easier to figure how best way to configure
	// once more unit equipments and equipment types are supported.
	case fighter:
		x := 2.0
		y := 20.0
		rotaryCannon := newOriginPolygon([]line{
			{Pos{X: 0, Y: 0}, Pos{X: 0, Y: y}},
			{Pos{X: 0, Y: y}, Pos{X: x, Y: y}},
			{Pos{X: x, Y: y}, Pos{X: x, Y: 0}},
			{Pos{X: x, Y: 0}, Pos{X: 0, Y: 0}},
		})
		uc.equipments = []equipment{
			{effectiveArea: effectiveArea{
				static: static{
					polygon: &rotaryCannon,
					offset:  Pos{X: -8, Y: -30},
				},
			}},
			{effectiveArea: effectiveArea{
				static: static{
					polygon: &rotaryCannon,
					offset:  Pos{X: 3, Y: -30},
				},
			}},
			// arbitrary equipment to prove radius (not only polygon) equipment behaves correctly.
			{effectiveArea: effectiveArea{
				static: static{
					radius: intPtr(10),
					offset: Pos{X: 10, Y: 40},
				},
			}},
		}
	case infantry, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic,
		missileAntiBallistic, flare, bombGravity, bomber, chaff:
	default:
		panic(fmt.Sprintf("setEquipments: cannot set equipments for unknown unit %s", uc.name.string()))
	}
}

// This log method is frequently used for debugging. Also it is typically logged at the end of
// unit.log(). Log frequently used info at the end for quicker debugging. Less used info can be
// parsed via tooling e.g. grep
func (uc unitConfig) log() []slog.Attr {
	var a = []slog.Attr{}
	a = append(a, slog.Int("radius", uc.radius))
	a = append(a, slog.Int("maxHealth", uc.maxHealth))
	a = append(a, slog.Int("attack", uc.attack))
	a = append(a, slog.Int("vision", uc.vision))
	a = append(a, slog.Int("reach", uc.reach))
	a = append(a, slog.Int("reactionsCount", len(uc.reactions)))
	a = append(a, slog.Bool("staticTarget", uc.isStaticTarget))
	a = append(a, slog.Bool("heatSeeking", uc.isHeatSeeking))
	a = append(a, slog.Bool("explodeWhenTerminal", uc.isExplodeWhenTerminal))
	a = append(a, slog.Bool("despawnWhenTerminal", uc.isDespawnWhenTerminal))
	a = append(a, slog.Uint64("terminalStepDelay", uc.terminalStepDelay))
	a = append(a, slog.Uint64("terminalStep", uc.terminalStep))
	a = append(a, slog.String("collisionLayers", uc.collisionLayers.string()))
	a = append(a, slog.String("trackLayers", uc.trackLayers.string()))
	a = append(a, slog.String("produceLayer", uc.produceLayer.string()))
	a = append(a, uc.consumables.log()...)
	a = append(a, slog.String("name", uc.name.string()))
	return a
}
