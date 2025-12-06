// A unit config provider simplifies the complex process of creating unit configs. It has several benefits
// but the core benefit is handling unit config fields dependent on other unit config fields.
//
// Benefits:
// - Unit config fields can have respective verifications to ensure valid values (e.g. radius cannot be negative).
// - Unit config fields which are dependent on other fields within the same unit are easy to set.
//   - The order for every set field is controlled. The dependent field is easily set after its required field.
//   - The dependent field can have respective verifications to ensure valid values (e.g. vision >= radius)
// - Unit config fields which are dependent on other unit configs are easy to set.
//   - All independent fields, for all units can be set before the fields dependent on other units are set.
// - It simplfies file exports.
//   - Consumers only interact with the single provider and the relatively simple provides method.
//   - Does not export multiple global unit configs
//   - In the future, I may make this a separate package, but for now using "set" methods clarifies exports.
// - Default values are more flexible.
//   - Implicit values can be set per field via a switch statement.
//   - Required explicit values can be set per field via a switch statement.
// - Better documentation and easier to tune configurations
//   - it is easy to understand a field when it is set for every unit within a single method.
// - It is easier to extend in the future since most likely upgrades will require dynamic unit configs.
//   - Requests to change configs will all be processed by the config provider.
//   - Multiple approachs are viable for consuming dynamic configs:
//     - Consumers may use the unit config provider as the single source of truth
//     - The provider may be the single source to provide updates to all consumers.
//
// Dependent fields have two anticipated sources:
// - The first source is sometimes it is easier to understand a config when it is defined via field X. But the config
//   must be processed via field Y. e.g. It is easier to understand unit A attacks unit B via an attack field. Then this
//   attack field is used to determine what layers the unit should search for a unit to attack.
// - The second source is when a field has high CPU usage to calculate at runtime. It is better for provider to
//   calculate and cache the value at startup, and retreive it during steady state runtime.

package main

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
)

// Provides unit configs to various consumers
type unitConfigProvider struct {
	unitConfigs
}

type unitConfigs map[name]unitConfig

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

// A unit action which should occur when a condition is met.
type reaction struct {
	when  func(*unit) bool
	do    func(*unit)
	next  uint64
	delay uint64
}

type consumable struct {
	count         int
	countMax      int
	fireDelay     step
	firesAt       func(uc unitConfig) bool
	firesAtLayers layers
}

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
	consumables            map[name]consumable
	isReactive             bool
	isEmitsHeat            bool
	isEmitsRadar           bool
	isHeatSeeking          bool
	isRadarGuided          bool
	isRequester            bool
	isSurvivesCollision    bool
	isSkipSpawnRadiusCheck bool
	isExplodeWhenTerminal  bool
	isDespawnWhenTerminal  bool
	isTerminalWhenMoveDone bool
	isStaticTarget         bool
	ignore                 func(u, t *unit) bool
	ignoreTarget           func(t *unit) bool
	ignoreCollision        func(other *unit) bool
	// This flag should rarely be used. If it is used often then write a function to avoid bugs.
	// This is safe to enable iff the unit only collides with units that have radiuses less
	// than state.spawnUnitRadiusLimit. For example, chaff only collides with radarGuided units
	// which are all less than the 20 pixels.
}

func newUnitConfigProvider() unitConfigProvider {
	configs := make(unitConfigs)
	for _, name := range names() {
		configs[name] = unitConfig{name: name}
	}
	for k, uc := range configs {
		uc.setUCFields()
		configs[k] = uc
	}
	ucProvider := unitConfigProvider{configs}
	ucProvider.setUCFieldsDependentOnOtherUCs()
	return ucProvider
}

func (ucProvider unitConfigProvider) setUCFieldsDependentOnOtherUCs() {
	// TODO: revist setConsumables once implementation is finished.
	// It may end up not being a dependent field but placing it here to be conservative.
	ucProvider.setConsumables()
	ucProvider.setReactions()
	ucProvider.setTracksLayers()
}

func (ucProvider unitConfigProvider) setTracksLayers() {
	for name, uc := range ucProvider.unitConfigs {
		layers := []layer{}
		if uc.tracks != nil {
			for _, other := range ucProvider.unitConfigs {
				if uc.tracks(other) {
					layers = append(layers, other.produceLayer)
				}
			}
		}
		uc.trackLayers = layers
		ucProvider.unitConfigs[name] = uc
	}
}

// TODO: Update consumables. It needs to handle independent and dependent fields similar to unitConfigProvider
func (ucProvider unitConfigProvider) setConsumables() {
	for _, uc := range ucProvider.unitConfigs {
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
		case infantry, missileCruise, missileBallistic, missileAntiBallistic, missileHeatSeeking, missileRadarActive, flare, bombGravity, chaff:
		default:
			panic(fmt.Sprintf("setConsumables: cannot set consumables for unknown unit %s", uc.name.string()))
		}
	}
}

func (ucProvider unitConfigProvider) provide(name name) unitConfig {
	return ucProvider.unitConfigs[name]
}

func (ucProvider unitConfigProvider) setReactions() {
	for name := range ucProvider.unitConfigs {
		var reactions []reaction
		switch name {
		case infantry, bombGravity, chaff, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare:
		case bomber:
			bomberUC := ucProvider.unitConfigs[bomber]
			reactions = []reaction{
				{
					when: func(bomber *unit) bool {
						// TODO: bresenham is not guaranteed to get to final position. I have not seen
						// aircraft not reach their final destination but missiles have. Need to ensure
						// this is not error prone or make it more conservative.
						return bomber.Pixel.equals(bomber.sourcePosition) && bomber.isPayloadNotFull()
					},
					do: func(bomber *unit) {
						bomber.refillPayload()
					},
					delay: bomberUC.reloadPayloadDelay,
				},
				ucProvider.reactFlare(bomberUC.consumables[flare].fireDelay),
				ucProvider.reactChaff(bomberUC.consumables[chaff].fireDelay),
			}
		case fighter:
			fighterUC := ucProvider.unitConfigs[bomber]
			missileRadarActiveUC := ucProvider.unitConfigs[missileRadarActive]
			missileHeatSeekingUC := ucProvider.unitConfigs[missileHeatSeeking]
			reactions = []reaction{
				{
					when: func(fighter *unit) bool {
						// TODO: bresenham is not guaranteed to get to final position. I have not seen
						// aircraft not reach their final destination but missiles have. Need to ensure
						// this is not error prone or make it more conservative.
						return fighter.Pixel.equals(fighter.sourcePosition) && fighter.isPayloadNotFull()
					},
					do: func(fighter *unit) {
						fighter.refillPayload()
						// TODO: when a missile terminates, should it remove the target from skip targets
						// rather than do this during restock. Not sure if this is a realistic upgrade.
						// Also, multiple ordnance might make this complex. Just because a single ordnance skips
						// a target does not mean others should or can this be avoided?
						fighter.skipTargets = []id{}
					},
					delay: fighterUC.reloadPayloadDelay,
				},
				{
					when: func(fighter *unit) bool { return fighter.nearest != nil },
					do: func(fighter *unit) {
						var uc unitConfig
						// TODO: determine fighter missile preference? Configurable by player?
						// Arbitrarily set radar to have higher priority than heat seeking.
						switch {
						case fighter.consumables[missileRadarActive].count > 0 && !missileRadarActiveUC.ignoreTarget(fighter.nearest):
							uc = missileRadarActiveUC
						case fighter.consumables[missileHeatSeeking].count > 0 && !missileHeatSeekingUC.ignoreTarget(fighter.nearest):
							uc = missileHeatSeekingUC
						default:
							return
						}
						missiles := fighter.consumables[uc.name]
						missiles.count -= 1
						fighter.consumables[uc.name] = missiles
						spawn := newUnit(
							uc,
							position{layer: uc.produceLayer, Pixel: fighter.Pixel},
							fighter.player,
							[]position{fighter.nearest.position},
							[]action{},
							GridPosition{},
							fighter.selected)
						spawn.sourceUnit = fighter.id
						fighter.requests = append(fighter.requests, request{spawn: &spawn})
					},
					// TODO: fix this. delay needs to be conditional on some unit config.
					delay: 50,
				},
				ucProvider.reactChaff(fighterUC.consumables[chaff].fireDelay),
				ucProvider.reactFlare(fighterUC.consumables[flare].fireDelay),
			}
		default:
			panic(fmt.Sprintf("setReactions: cannot set reactions for unknown unit %s", name.string()))
		}
		uc := ucProvider.unitConfigs[name]
		uc.reactions = reactions
		ucProvider.unitConfigs[name] = uc
	}
}

func (ucProvider unitConfigProvider) reactChaff(fireOrdnanceDelay uint64) reaction {
	chaffUC := ucProvider.unitConfigs[chaff]
	return reaction{
		when: func(u *unit) bool {
			return u.nearest != nil && u.nearest.name == missileRadarActive && u.consumables[chaff].count > 0
		},
		do: func(u *unit) {
			newChafff := newUnit(
				chaffUC,
				position{layer: chaffUC.produceLayer, Pixel: u.Pixel},
				u.player,
				[]position{u.position},
				[]action{},
				GridPosition{},
				u.selected)
			u.requests = append(u.requests, request{spawn: &newChafff})
			consumeable := u.consumables[chaff]
			consumeable.count -= 1
			u.consumables[chaff] = consumeable
		},
		delay: fireOrdnanceDelay,
	}
}

func (ucProvider unitConfigProvider) reactFlare(fireOrdnanceDelay uint64) reaction {
	flareUC := ucProvider.unitConfigs[flare]
	return reaction{
		when: func(u *unit) bool {
			return u.nearest != nil && u.nearest.name == missileHeatSeeking && u.consumables[flare].count > 0
		},
		do: func(u *unit) {
			newFlare := newUnit(
				flareUC,
				position{layer: flareUC.produceLayer, Pixel: u.Pixel},
				u.player,
				[]position{u.nearest.position},
				[]action{},
				GridPosition{},
				u.selected)
			u.requests = append(u.requests, request{spawn: &newFlare})
			consumeable := u.consumables[flare]
			consumeable.count -= 1
			u.consumables[flare] = consumeable
		},
		delay: fireOrdnanceDelay,
	}
}

func names() []name {
	return []name{infantry, fighter, bomber, bombGravity, missileHeatSeeking, missileRadarActive, missileBallistic,
		missileCruise, missileAntiBallistic, flare, chaff}
}

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
		firesAt = func(uc unitConfig) bool { return uc.isEmitsHeat && uc.name != flare }
	case missileRadarActive:
		firesAt = func(uc unitConfig) bool { return uc.isEmitsRadar && uc.name != chaff }
	case bombGravity, flare, chaff:
	default:
		panic(fmt.Sprintf("setFiresAt: cannot set fires at for unknown consumable %s", name.string()))
	}
	c.firesAt = firesAt
}

func (c *consumable) setFiresAtLayers() {
	layers := []layer{}
	// TODO: Fix in future commit. This commit is already large.
	/*
		for _, other := range unitConfigs {
			if c.firesAt(other) {
				layers = append(layers, other.produceLayer)
			}
		}
	*/
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

func (uc *unitConfig) setIsSurvivesCollision() {
	var survivesCollision bool
	switch uc.name {
	case infantry, chaff:
		survivesCollision = true
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare,
		bombGravity, fighter, bomber:
	default:
		panic(fmt.Sprintf("setIsSurvivesCollision: cannot set isSurvivesCollision for unknown unit %s", uc.name.string()))
	}
	uc.isSurvivesCollision = survivesCollision
}

func (uc *unitConfig) setIsEmitsHeat() {
	var isEmitsHeat bool
	switch uc.name {
	case fighter, bomber, flare:
		isEmitsHeat = true
	case infantry, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, bombGravity, chaff:
	default:
		panic(fmt.Sprintf("setIsEmitsHeat: cannot set isEmitsHeats for unknown unit %s", uc.name.string()))
	}
	uc.isEmitsHeat = isEmitsHeat
}

func (uc *unitConfig) setIsEmitsRadar() {
	var isEmitsRadar bool
	switch uc.name {
	case fighter, bomber:
		isEmitsRadar = true
	case infantry, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare, bombGravity, chaff:
	default:
		panic(fmt.Sprintf("setIsEmitsRadar: cannot set isEmitsRadar for unknown unit %s", uc.name.string()))
	}
	uc.isEmitsRadar = isEmitsRadar
}

func (uc *unitConfig) setIsRequester() {
	var isRequester bool
	switch uc.name {
	case fighter, bomber:
		isRequester = true
	case infantry, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare, bombGravity, chaff:
	default:
		panic(fmt.Sprintf("setIsRequester: cannot set isRequester for unknown unit %s", uc.name.string()))
	}
	uc.isRequester = isRequester
}

func (uc *unitConfig) setIsReactive() {
	var isReactive bool
	switch uc.name {
	case fighter, bomber:
		isReactive = true
	case infantry, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare, bombGravity, chaff:
	default:
		panic(fmt.Sprintf("setIsReactive: cannot set isReactive for unknown unit %s", uc.name.string()))
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
		panic(fmt.Sprintf("setProduceLayer: cannot set produceLayer for unknown unit %s", uc.name.string()))
	}
	uc.produceLayer = produceLayer
}

func (uc *unitConfig) setTerminalStepDelay() {
	var terminalStepDelay step
	switch uc.name {
	case bombGravity, chaff:
		terminalStepDelay = 300
	case infantry, fighter, bomber, missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare:
	default:
		panic(fmt.Sprintf("setTerminalStepDelay: cannot set terminal step delay for unknown unit %s", uc.name.string()))
	}
	uc.terminalStepDelay = terminalStepDelay
}

func (uc *unitConfig) setIsExplodeWhenTerminal() {
	var isExplodeWhenTerminal bool
	switch uc.name {
	case bombGravity, missileCruise, missileBallistic:
		isExplodeWhenTerminal = true
	case infantry, fighter, bomber, missileRadarActive, missileHeatSeeking, missileAntiBallistic, flare, chaff:
	default:
		panic(fmt.Sprintf("setIsExplodeWhenTerminal: cannot set is explode when terminal for unknown unit %s", uc.name.string()))
	}
	uc.isExplodeWhenTerminal = isExplodeWhenTerminal
}

func (uc *unitConfig) setIsStaticTarget() {
	var isStatictarget bool
	switch uc.name {
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare, bombGravity, chaff:
		isStatictarget = true
	case infantry, fighter, bomber:
	default:
		panic(fmt.Sprintf("setIsStaticTarget: cannot set isStaticTarget for unknown unit %s", uc.name.string()))
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
		panic(fmt.Sprintf("setIsDespawnWhenTermianl: cannot set isDespawnWhenTerminal for unknown unit %s", uc.name.string()))
	}
	uc.isDespawnWhenTerminal = isDespawnWhenTerminal
}

func (uc *unitConfig) setIsTerminalWhenMoveDone() {
	var isTerminalWhenMoveDone bool
	switch uc.name {
	case missileRadarActive, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare:
		isTerminalWhenMoveDone = true
	case infantry, fighter, bomber, bombGravity, chaff:
	default:
		panic(fmt.Sprintf("setIsTerminalWhenMoveDone: cannot set isTerminalWhenMoveDone for unknown unit %s", uc.name.string()))
	}
	uc.isTerminalWhenMoveDone = isTerminalWhenMoveDone
}

func (uc *unitConfig) setIsRadarGuided() {
	var isRadarGuided bool
	switch uc.name {
	case missileRadarActive:
		isRadarGuided = true
	case infantry, fighter, bomber, bombGravity, chaff, missileHeatSeeking, missileCruise, missileBallistic, missileAntiBallistic, flare:
	default:
		panic(fmt.Sprintf("setIsRadarGuided: cannot set isRadarGuided for unknown unit %s", uc.name.string()))
	}
	uc.isRadarGuided = isRadarGuided
}

func (uc *unitConfig) setIsHeatSeeking() {
	var isHeatSeeking bool
	switch uc.name {
	case missileHeatSeeking:
		isHeatSeeking = true
	case infantry, fighter, bomber, bombGravity, chaff, missileRadarActive, missileCruise, missileBallistic, missileAntiBallistic, flare:
	default:
		panic(fmt.Sprintf("setIsHeatSeeking: cannot set isHeatSeeking for unknown unit %s", uc.name.string()))
	}
	uc.isHeatSeeking = isHeatSeeking
}

func (uc *unitConfig) setIsSkipSpawnRadiusCheck() {
	var isSkipSpawnRadiusCheck bool
	switch uc.name {
	case chaff:
		isSkipSpawnRadiusCheck = true
	case missileHeatSeeking, infantry, fighter, bomber, bombGravity, missileRadarActive, missileCruise,
		missileBallistic, missileAntiBallistic, flare:
	default:
		panic(fmt.Sprintf("setIsSkipSpawnRadiusCheck: cannot set isSkipSpawnRadiusCheck for unknown unit %s", uc.name.string()))
	}
	uc.isSkipSpawnRadiusCheck = isSkipSpawnRadiusCheck
}

func (uc *unitConfig) setCollisionLayers() {
	var collisionLayers layers
	switch uc.name {
	case infantry, fighter, bomber, bombGravity, missileCruise, missileBallistic, missileAntiBallistic, flare:
		collisionLayers = []layer{uc.produceLayer}
	case missileRadarActive:
		collisionLayers = []layer{lOrdnance, lAircraft, lChaff}
	case missileHeatSeeking:
		collisionLayers = []layer{lOrdnance, lAircraft}
	case chaff:
		collisionLayers = []layer{lOrdnance}
	default:
		panic(fmt.Sprintf("setCollisionLayers: cannot set collision layers for unknown unit %s", uc.name.string()))
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

// This log method is frequently used for debugging. Also it is typically logged at the end of unit.log().
// Log frequently used info at the end for quicker debugging. Less used info can be parsed via tooling e.g. grep
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
	// TODO: refactor this lower e.g. into consumable log() method?
	if len(uc.consumables) == 0 {
		a = append(a, slog.String("consumables", "none"))
	}
	for k, v := range uc.consumables {
		a = append(a, nest(k.string(), v.log()))
	}
	a = append(a, slog.String("name", uc.name.string()))
	return a
}
