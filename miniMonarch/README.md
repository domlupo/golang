# Mini Monarch

A prototype for a simple war/4x game. Simple like mini metro or mini motorways. 
War game like civ, hearts of iron, 9 kings etc. (also other popular war games on steamdb). 
Think about how this game relates to civilization revolution, catan, RISK, 
starcraft desert strike, totally accurate battle simulator, etc.
Graphics may end up similar to the gnorp analogue.
Is there a mobile equivalent for war/4x games?

---

## Top 3 Ranked Concerns : Risks

1. Design and iterate on Local path finding and flux. -> Most unknown but fundamental game engine decision.
2. Faster iteration from better inspection and debugging tooling. -> Don't waste engineering effort due to poor prioritization.
3. Playtesting MVP. -> Don't invest major engineering effort into wrong tasks.

## Checkpoints : Purposes

### Next Checkpoints

- Determine if I should implement grid chunking or something else for game engine to support higher unit count.
  - Investigate why collides with Node and collides with Unit are taking up so much CPU. 
  - It is very intersting that collides with node is taking up similar CPU as collides with Units.
  - I am not convinced grid chunking will improve collides with performance.
  - The fix for collides with may be lower hanging fruit than grid chunking.
  - Collides fix may be a big enough benefit that I can punt grid chunking for awhile.
  - Analyze target is too CPU inefficent
  - Collides with unit is too CPU inefficient

- Improve guided air missile
  - Fix bug where ground missile (not air missile) is not moving for computer player
  - auto despawns for non computer players. This is probably wrong behavior.
  - Determine target locking behavior for human player and computer player.

- Determine how to handle flux for units that never had targets or units just produced. Lots of issues.
  - Do you make flux never move the unit? This creates bimodal behavior where some units respect flux and others dont?
  - Do you make the unit that was never moved by flux add more or less to flux magnitude?
  - Not actually spawning the unit until a target is set might be a very simple solution but the
    player experiencee might be terrible.
  - How do you handle missiles blowing up from flux moving them.
  - How do you handle missiles moving without target? Can players abuse this?

- Invest in a visual debugger:
  - First invest in: It can auto select/highlight a specific unit when it spawns
  - Other much lower priority improvements that can be added when needed:
    - It can can ping the pixel pos where the bug is occurring.
    - It can ping differently based on layer.
    - It can also do something like log unit and pause game so I can read logs?

- Fix Bresenham's algorithm edge cases:
  - Bresenham error needs to reset itself for several reasons. I still need to handle:
    - Determine if the target moving is handled. There may be an elegant solution here.
    - If it is not handled determine how to fix. Structure may need updated.
  - A potentially good contract:
    - Multiple sources (state, computer, etc) update units various targets (flux, local, global, etc)
    - Unit must keep itself in sync
  - Determine if a different structure is better for Bresenham

- Iterate on group and local pathfinding
  - Flux can be iterated and tuned a lot. I only implemented a basic solution.
  - Flux resolution can be improved by changing flux grid size at compile time. 
  - Flux grid can be changed at runtime as well but probably unneded complexity. 
  - How units react to flux can depend on unit type or be the same for all units.
  - How units contribute to flux can depend on unit type or be the same for all units.
  - Determine better flux structure and contracts.

- Determine path finding goals for individual units to decide on path finding algorithm:
  - Unit -> player experience -> path finding algorithm
  - Infantry -> Form tight territory lines of a specific thickness -> 
  - Sniper -> Clump with other snipers behind infantry line to allow other units pass ->
  - unguidedGroundMissile ->
  - guidedAirMissile ->
  - mobile SAM ->
  - tank -> 
  - bomer -> 
  - fighter ->

- Add more ground units for local path finding iteration
  - tank (counters infantry, secondary mode to push unit?)
  - mobile SAM (counters missiles)
    - Ground missile exploding in air does no damage?
    - Should air missile exploding in air do damage?
    - Should air missile only get one retarget on missle? An upgrade could allow multiple retargets besides the first?
    - Units may need different pathing algorithm or reactions to flux? 
      - How to avoid this?
      - How to simplify this?
      - How to make this CPU performant?
    - Targets different layer than it moves
    - How will AI handle placement? 
      - An easy implementation might be to be near node.
      - An easy implementation might be to stay away from infantry
      - Counter missile will either need to travel faster than tracked missile or anticipate movement?
    - Added guided missile variant will probably help but might be low priority?

- Improve flux debugging
  - Add flux square selection
  - Add selected flux logging for detailed tracing

- Improve computer
  - Update computer to be beatable by humans
  - Should computer units prioritize nodes or units?
  - Investigate, there may be a bug where when a node is captured by a computer player, unit targets incorrectly become node
    when a local target is already acquired.

- Create more node setup test scenarios
  - Able to easily reproduce and test precise setups for edge case issues
  - Able to randomly generated setups to find hidden edge case issues

- Polish UI
  - Currently it is very hard to fight computer due to clunky UI
  - Easier way to select all nodes of same type with double click
  - Better automated targetting from units. e.g. units don't auto target nodes when they should.

- Analyze and iterate on game design and game loop:
  - Check design feasability 
  - Improve game

- clean up bugs and hacks:
  - There are enough to cause bugs while testing game design and game loop

- Set a deadline and plan for a polished MVP
  - This should be multiplayer playable by friends via browser.
  - A vertical slice of MVP features. 
  - If feedback is good, then I can work on adding more features.
  - If feedback is good, then I can work on polishing existing features.

- Refactor for more idiomatic go code as I read learning Go. Already did var declare and switch statements.

---

### Future Checkpoints

- Invest in a visual debugger:
  - First invest in: It can auto select/highlight a specific unit when it spawns
  - Other much lower priority improvements that can be added when needed:
    - It can can ping the pixel pos where the bug is occurring.
    - It can ping differently based on layer.
    - It can also do something like log unit and pause game so I can read logs?

- Write unit tests for missing automated tests for function getHexPixelPosition

- Verify if Bresenham error is handle correctly for Flux movements.

- Deterministic seeds and deterministic game engine:
  - I can easily run simulations to analyze game balance

- Determine if a node should be a type of unit rather than a separate struct
  - I am leaning towards overall it would be better to just make a node a unit
  - Other units will probably produce units e.g. plane create flare, plane create paratrooper, etc
  - if node count stays static there is probably some benefits to stay separate 
  - if nodes do not move there is probably some benefit to stay separate

- Determine if unit collision not occuring for static units is correct behavior.
  - This behavior probably redcues CPU usage
  - This behavior probably leads to more bugs where dynamic unit colliding into static unit has no effect occur on static unit.

- Determine behavior for overkill
  - what order should main game actions occur e.g. target acquisition, attack units, attack nodes, attack positions, move units?
  - How is overkill handled?
  - Should units retarget if they lose their initial target?

- Should attack units occur deterministicly?
  - maps in go by default are not deterministic. Currently, a unit can oscillate attacks. Probably wrong behavior

- Determine behavior for multi hit melee combat 
  - Expose issues with combat
  - Expose issues with timing
  - Expose issues with pathing
  - Will probably lead to bug fix for units not killing correctly
  - May help or hurt efforts towards local pathfinding. I suspect this will be coupled with local pathfinding.

- Tool / CLI which launches multiple windows at once to test multiple scenarios
  - Faster test loop

- save file: 
  - Will expose constraints
  - Will expose platform concerns
  - May expose issues with work so far 

- Determine if frames per second lock is needed:
  - fundmental issue that I have not considered

- Release to web:
  - learn how to deploy
  - easier to share with friends
  - test web based issues rather than only mac

- Release to phone: 
  - check phone controls
  - check resolution

- local multiplayer:
  - Initial research shows webRTC is probably better but more complex than websockets

---

### Punted Checkpoints

- Should I add a config file?

- Investigate if logging impacts game performance
  - Might be able to fix via build tags to compile different code for release build

- Polish logger
  - Improve always log logic. Right now it is a bit confusing if always log surpasses logging level. I should
    make this more clear probably through simplifying rather than adding more complexity.
  - Add logs throughout (global path, local path, attack, move, collision, etc)
  - Add runtime dynamic log sampling rates. Is this even possible? Are the benefits worth it?
  - Paused debugging? e.g. Updating unit during state.go execution can log state at game level.

- Controls which work on computer, steam deck, phone, etc: 
  - multiple platform means more sales
  - customers can pick favorite device if multiple owned

- Polish missile:
  - Damage rather than instant kill?
  - Fall off damage?
  - Add unit that counters missile and expands game logic?

- Update unit pathing algorithm to Theta* with same cost orthgonal and diagonal movement: (https://en.wikipedia.org/wiki/Theta*)
  - Units will be able to go around obstacles

- initial map generation:
  - Omitting certain counters or land features from maps might give ticket to ride vibes where
    a different map significantly changes gameplay.
  - Omitting certain tech tree options is another way to randomize games.
  - Omitting certain tech tree options is a bit like randomized civ pick in AOE

- Tech Tree:
  - ?

- Different terrains: 
  - Careful, this may be hard to implement due to pathfinding algorithms assuming same cost movement

---

## Commands

```find . -type f -name "*.go" -exec wc -l {} +```

```go tool pprof -http=:6060 cpu.prof```

```go tool pprof -http=:6060 mem.prof```

---

## Closed Decisions

1. No worker unit or resource collection
    - Adds complexity
    - Nodes can produce cargo units which more aligned with multi use nodes
    - No guaranteed benefit
2. No civilization era progressions:
    - Directly adds design and balancing complexity
    - Indirectly adds complexity and scope throughout game e.g. more units
    - No guaranteed benefit
    - Tech tree is a simpler feature I can implement that has a lot of overlap with era progression.
3. Save center rather than top left of unit, nodes, etc into state:
   - Easier to design game/engine
   - Easier to test game/engine
   - Easier to debug game/engine
4. Same cost Orgthogonal and Cardinal movement:
   - Simplifies path finding
   - Lowers CPU cost
   - Lowers overall game time
5. Don't do 3d globe visual for game:
   - Very complex
   - No guaranteed benefit
   - Performance concerns
6. Theme will be modern land and air units
   - Naval units are unneeded for main mechanic and game loop. Also huge scope creep.
   - There is a large enough variety with only land and air units.
   - Medieval theme would not allow air units
7. Hexagonal grid rather than square grid
  - The aesthetics of a verticl layout hexagonal grid is much nice than a square grid
    - The pattern is more pleasing.
    - The border is more pleasing.
    - More variety to connect grid nodes.
    - Multiple connected hexagon nodes is more pleasing.
    - It makes the game seem more organic rather than synthetic which match war theme.
    - Specifically, the same infantry lines are formed but are more pleasing with a hexagon background.
    - Hexagon is closer to a circle shape for easier art, hacky collision detection, etc
  - Pixel Position to Hexagon index calcuation is not CPU intensive (floating point nor sqrt math needed)
  - Grid chunking and flux only requires 6 adjacent rather than 8.
  - Matches existing games better like Catan, Civ, etc
  - Optional hexagonal and square grid is more code complexity, work, bimodal behavior, etc
8. Hexagons are left side side inclusive, right side exclusive. 
  - This was chosen mostly arbitrarily. Probably simplifies some code.
  - The remaining 4 sides are sloped which determines tie breaks nicely.

---

## Game Design Questions

- Should air units have impassable terrain e.g. high mountain? Will make path finding harder for air units
- Stopping node unit production in order to advance tech tree, improve node is probably a good form of risk -> reward.
- Should nodes support unit cycling e.g. produce missile, then produce infantry, then produce tank, etc
- Are both low altitude and high altitude needed or would a single air layer be sufficient?
- Chemical gas can push enemies back unless they have gas masks
- Should nodes have a max unit capacity e.g. This node stops producing since 10 units it produced are still alive
- Should there be low altitude and high altitude flux? I suspect it is better to have separate maps rather
  than a combined map as for units because there will be less interactions between flux layers but units
  do interact somewhat between layers.
- Should bombs dropped from planes have health and can die?
- Should paratroopers be able to be shot at first?
- Plane collisions would be awesome?
- air units can produce other units e.g. missiles, paratroops, etc?
- Should units besides explosives detonate?
  - e.g. A tank being hit should explode damaging nearby units. Can use unit explode radius variable for this.
  - A missile detonating should detonate other missiles in air
  - This could create interesting explosion chains
- Should there be a game camera:
  - No camera can easily be the best option
  - One option is starcraft / civ where it panes around a larger area
  - One option is mini metro and mini motorways style that slowly zooms out to the full area?
- Should the focus of the game be positioning, logistics and decisions rather than micro?
  - It seems interesting whens units get into long term fights
  - Reminds me of the board game go RISK and Catan
  - Reminds me of the military diagrams that show battles over time
  - May be able to utilize the idea of target change units
  - Early fight for position, long term fight for ?
- top down graphics or iso or 3d?
- Players get to pick starting positions?
- Hidden or fully visible grid?
- APM lock on nodes or units so you cannot spam actions? Think about sokpop game clickyland where you are limited by clicks.
- Campaign mode. Specific battles, war or countries? Similar to cities in mini metro and mini motorway.
- Tech tree?
- Special nodes?
- Should nodes only have a single grid offset (offset = 0 is acceptable) or should offset differ? 
- Node expansion overtime?
- Grids expand over time to cities, cities can merge together, merged cities are more powerful?
- Buildings?
- Do I want to stand out from civ, or be a simpler civ e.g. factorio -> shapez. How does this affect marketing?
- Should players be able to target thier own nodes?
- Should a node be able to target itself?
- Cities have different units, different unit ratios, or different counts produced? e.g. Catan
- Traders or cargo can go to a node which you own to transfer production of city production
- Ammo limit for units?
- Nodes with defense system
- How is node unit production handled?
- Theme is the gnorp analogue?
- Game name?
  - gridlock
  - gridlocked
  - miniMilitary
  - miniMilitia
- Replay feature where you can see a battle recap?
- Should unit behavior be toggable on a global scale? E.g. tanks have fire mode vs push mode

---

## Misc Notes

- Research civ concepts and features e.g. culture victory, trade, etc
- Research polytopia game
- Research other mobile platform "civ" games like unciv

---

### Multi-agent and Local Pathfinding Brainstorm

- Units with different count and layers can employ different path finding algorithms.
  - total number of bombers will be low relatively fewer collisions
  - total number of infantry will be high with many collisions

- Pathfinding algorithm buckets:
  - Static obstacles: 
    - Theta*
    - A*
  - Dynamic obstacles: 
    - D*
    - D* focused
    - D* lite
    - LPA* aka Lifelong Planning A*
    - Conflict Based Search
    - Enhanced Conflict Based Search
    - Increasing  cost tree seach
    - Geometric MAPF or MAPF for large agents (LAMAPF)
    - Lifelong planning conflict-based search (LPCBS)
  - Group:
    - flow fields
    - flocking behaviors
    - dynamic obstacle avoidance
    - continuum crowds

- Other consideration:
  - Anonymous Multi-agent pathfinding: Set of targets to reach but specific agent does not matter
  - Can use nav meshes to optimize or simplify searches
  - Can use window searches aka use smaller searches to lower CPU use 
  - Rapidly-exploring random tree can create randomized pathing 

### Unit List Brainstorm

- Questions:
  - Which units will counter other units?
  - Should units have upgrades?
  - Should certain "civs" have better units?
  - Should planes need to fly back to airport to refuel?
  - Should planes need to fly back to airport to rearm?
  - Should planes need to fly back to airport to repair?
  - Should plane maintenance be much quicker than building a new plane?
  - How can other units be unique and fit with node mechanic?
  - Should I do purpose exercise with units like I do with checkpoints
  - Is there a good way to counter artillery?
  - Heat seeking missile (it can be an upgrade but does not have to be)
  - Jet heat flare (it can be an upgrade but does not have to be)
  - Should ground units like tank, artillery, etc need to refuel and resupply?
  - Missile upgrade makes it smaller?
  - SAM upgrade or configurable number of rockets sent?
  - Prior to the game and during game allow players to tune unit configs? 
    e.g. 1 SAM in army send 1 rocket, 2 SAM in army send 1 rocket from all in range or 
    send 2 rockets from closest? Army could have seensible defaults but then players
    can customize?

- Units:
  - Land:
    - Done:
      - infantry
      - sniper
    - High Priority:
      - tank
      - mobile SAM (shoots missile in low altitude or high altitude)
    - Low Priority:
      - mobile artillery / rocket artillery (thought is low priority because SAM pathing more complex. Could be wrong)
      - marine (cross water)
      - anti tank infantry
      - troop transport
      - cargo truck (can transfer production between nodes?)
      - unit to change targets
      - a unit which pushes other units out of way (tank?)
      - train track builder?
  - Air:
    - Done:
      -
    - High Priority:
      - bomber
      - stealth bomber
      - fighter 
      - stealth fighter
    - Low Priority:
      - helicopter / gun ship
      - paratrooper
      - cargo plane (can transfer production between nodes?)
      - drone
      - unit to change targets?
  - City aka Node:
    - Done
      - missle (low altitude)
    - High Priority
      - guided missle (low altitude)
      - anti missile system
      - anti artillery system
      - anti air gun
    - Low Priority
      - nuke
      - satelite
      - a unit to change targets
  - Other
    - air missile (low altitude, high altitude)
    - ground missile (low altitude, high altitude)
    - bomb (low altitude)
    - artillery shell (low altitude)
    - chemical gas
    - jet flare

---