 # Gridlock

A prototype for a simple war/4x game. Simple like mini metro or mini motorways. 
War game like civ, hearts of iron, 9 kings etc. (also other popular war games on steamdb). 
Think about how this game relates to civilization revolution, catan, RISK, 
starcraft desert strike, totally accurate battle simulator, etc.
Graphics may end up similar to the gnorp analogue.
Is there a mobile equivalent for war/4x games?

---

## Top 3 Ranked Concerns : Risks

1. Design and iterate on Local path finding (maybe flux) -> Most unknown but fundamental game engine decision.
2. Faster iteration from better inspection and debugging tooling. -> Don't waste engineering effort due to poor prioritization.
3. Playtesting MVP. -> Don't invest major engineering effort into wrong tasks.

## Checkpoints : Purposes

### Next Checkpoints

- Fix units colliding into hexagon obstacles:
  - Can add collision detection for hexagon obstacles. Probably need spat partitioning
  - Can make hexagon obstacles visually smaller

- Fix HACKS added during multi obstacle implementation.
  - Determine if high level unit tests should be added now or punted till later.
  - Handle setting up visiblity grids edge cases e.g. obstacle on border, etc

- Optimize multi obstacle pathing:
  - Long term:
    - May need to use local collision avoidance with other units since obstacles should funnel units together
    - May want to not use center of hexagon as target. 
  - Investigate life long planning more. I think this is technically what I am doing since
    targets are stored and I have not read much about this.
  - If I do hex edge detection for LOS instead of Bresenham. multiple optimizations are possible:
    - remove internal obstacle edges
    - combine external edges that form straight lines
    - I can reduce the number of edges checked by doing a high pass then a low pass. e.g performance 
      4 edge checks per obstacle boxed by 90 degree angles. If I hit conservtive box then do a deep
      check for any boxs hit.
  - I am very unsure about this optimization, but for single obstacle, I can maybe reduce the number 
    of visiblity grid points checked  by doing a greedy apporach and jumping as far as possible. Then
    expanding as needed but it seems very complex. This also might be worst in some cases e.g. if the 
    target is within the obstacle and the jump missed.
  - Can store what edge intersects are between what points cw and ccw as well but right now
    this is done at runtime. I think it is a bit too complex to optimize ahead of time.
  - Need to handle edge case where line of sight reaches target, but unit travel
    collides with impassable terran due to unit radius. To solve this for circular units,
    At game startup, store the distance of each pixelPos to nearest impassable terrain and
    also the pixelPos of the impassable terrain. Then colliison can be easily checked.
    Non circular collision can be handled at a later time.
  - I think caching at startup is going to be the best speed up for global targeting. If I can 
    verify what visbiltiy grids a unit should take when targetting a pixel pos with a given grid 
    position, I can then only calculate the final path. If obstacles and nodes are not dynamic this 
    will probably be easy. Essentially, avoid cache invalidation.
  - Current theta* implementation is not finding shortest path for multi obstacle. 
  - For pixel pathing may optimize by adding two fields:
    - nearestUnwalkabledDistance int
    - nearestUnwalkable pixelPosition

- Determine if LOS algorithm should be updated from bresenham to check obstacle line intersects
  - line intersect check will be probably be much more performant
  - line intersect check will probably be much better than switching from map to array for cache locality.
  - current CPU profiling shows map access during lineOfSight checks are taking up significant CPU
  - can reduce number of line checks at game start e.g. 2 hexagons connecetd only have 10 sides, not 12.

- Determine if I need to implement separating axis theorem for collision other than (circle <-> circle)
  - Right now node collision is hacky. It is just a circle that is smaller than the hexagon.
  - Unpassable terrain e.g. water, high altitude mountains will probably require this for a* pathing.
  - Nodes are hexagon. Determine if SAT is the simplest and most performant collision check for hexagons.
  - Units like tank may need represented as a square. 

- Determine if spatial partitiong node colliison optimization is needed.
  - right now node collision naively checks all nodes.

- Improve guided air missile
  - If analyze target is migrated to linked list, need to update number of hexagons checked based on missile vision.
  - Fix bug where ground missile (not air missile) is not moving for computer player
  - auto despawns for non computer players. This is probably wrong behavior.
  - Determine target locking behavior for human player and computer player.

- Fix Bresenham's algorithm edge cases:
  - Bresenham error needs to reset itself for several reasons. I still need to handle:
    - Determine if the target moving is handled. There may be an elegant solution here.
    - If it is not handled determine how to fix. Structure may need updated.
  - A potentially good contract:
    - Multiple sources (state, computer, etc) update units various targets (flux, local, global, etc)
    - Unit must keep itself in sync
  - Determine if a different structure is better for Bresenham

- Determine path finding goals for individual units to decide on path finding algorithm:
  - Determine what units I want to remove. i think most ground units are unneeded.
  - Unit -> player experience -> path finding algorithm
  - Infantry -> Form tight territory lines of a specific thickness -> 
  - etc 

- Add more ground units for local path finding iteration
  - Determine what units I want to remove. i think most ground units are unneeded.
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

- Set a deadline and plan for a polished MVP
  - This should be multiplayer playable by friends via browser.
  - A vertical slice of MVP features. 
  - If feedback is good, then I can work on adding more features.
  - If feedback is good, then I can work on polishing existing features.

- Analyze and iterate on game design and game loop:
  - Is game fun?
  - Are you designing features to improve player experience? 
  - Are you coding features to improve player experience?
  - Are you priotizing and cutting features as needed?

- clean up bugs and hacks:
  - There are enough to cause bugs while testing game design and game loop

- Refactor for more idiomatic go code as I read learning Go.

---

### Future Checkpoints

- Iterate on local path finding solutions:
  - If I use noticeable CPU I need to guarantee it makes the global and local situation better. I also need
    to make sure it is deterministic and helps the player experience. I don't want unknown experience in the game.
    I would prefer less total unit count with deterministic pathfinding than more unknown unit pathing.
  - Desired player experiences:
    - Don't allow full blockages
    - When the gap is fully blocked, units pack in. 
    - Units apporaching other units attacking, path around to form infantry lines. 
  - Start with simple methods that may have more future uses:
    - flow fields
    - flocking behaviors
    - steering
  - More complex methods can be looked into later if needed:
    - CBS (a* or theta* probably possible)
    - Enhanced Conflict Based Search (Not optimal, but faster computation) https://www.ijcai.org/Proceedings/15/Papers/110.pdf
    - Meta agent conflict based search (optimization for units which collide frequently with eachother)
    - Geometric MAPF or MAPF for large agents (LAMAPF) https://ojs.aaai.org/index.php/AAAI/article/view/4756 

- Update game setups where nodes and obstacles are combined

- Determine if for CPU reasons, I should not mix layers in the same map and always manage them by state.go

- Refactor how state, setups, and pathing is initialized. The current structure is poor.

- How can I use more unit tests? Or how can i write automated high level tests.

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

- Do I need to think about layers more?
  - Right now I think this decision can only be known far in the future and it should be 
    much easier to refactor three layer back to two layer.
  - two (groun vs air), three (ground, low altitude, high altitude) have major game differences
  - Should anything be on multiple layers e.g. high altitude mountains technically are on all 3.

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

- Release to web:
  - learn how to deploy
  - easier to share with friends
  - test web based issues rather than only mac

- Release to phone: 
  - check phone controls
  - check resolution

- local multiplayer:
  - Initial research shows webRTC is probably better but more complex than websockets

- Determine if flux iteration is useful:
  - Need to implement specific local path finding before flux iteration.
  - Even if group behavior is needed, methods other than flux might be better to invest in.
  - Flux resolution can be improved by changing flux grid size at compile time. 
  - Flux grid can be changed at runtime as well but probably unneded complexity. 
  - How units react to flux can depend on unit type or be the same for all units.
  - How units contribute to flux can depend on unit type or be the same for all units.
  - Determine better flux structure and contracts.
  - Does a flat top hexagon grid perform better for flux? If yes, is it worth it for potentially
    worse asethetics? I think flux can most likely be designed around this but the monitor display
    limitation for wide resolution suggests left vs right pathing is more likely where as phone 
    vertical display suggests top vs bottom pathing is more likely.

- Determine how to handle flux for units that never had targets or units just produced. Lots of issues:
  - Do you make flux never move the unit? This creates bimodal behavior where some units respect flux and others dont?
  - Do you make the unit that was never moved by flux add more or less to flux magnitude?
  - Not actually spawning the unit until a target is set might be a very simple solution but the
    player experiencee might be terrible.
  - How do you handle missiles blowing up from flux moving them.
  - How do you handle missiles moving without target? Can players abuse this?

- Improve flux debugging
  - Add flux square selection
  - Add selected flux logging for detailed tracing

---

### Punted Checkpoints

- Check if hexagon look up table can be optimized with static aka compile time setup rather 
  than runtime setup at game start and accessed during game time via heap allocated pointer.

- Set FPS or add setting for low performance platforms e.g. mobile?

- Should I add a yaml, toml or json config file that is loaded at runtime?
  - Toggles are already allowing runtime updates, rather than manually updating a text file.

- Invest in a visual debugger:
  - First feature would most likely be it can auto select/highlight a specific unit when it spawns
  - Other much lower priority improvements that can be added when needed:
    - It can can ping the pixel pos where the bug is occurring.
    - It can ping differently based on layer.
    - It can also do something like log unit and pause game so I can read logs?

- Logging affects game performance. Especially hot loop logs e.g. unit refresh
  - Can I use build tags for a development vs release build?
  - Can I conditionally compile?

- Polish logger
  - Improve always log logic. Right now it is a bit confusing if always log surpasses logging level. I should
    make this more clear probably through simplifying rather than adding more complexity.
  - Add logs throughout (global path, local path, attack, move, collision, etc)
  - Add runtime dynamic log sampling rates. Is this even possible? Are the benefits worth it?
  - Paused debugging? e.g. Updating unit during state.go execution can log state at game level.

- Controls interface which work on computer, steam deck, phone, etc: 
  - multiple platform means more sales
  - customers can pick favorite device if multiple owned

- Polish missile:
  - Damage rather than instant kill?
  - Fall off damage?
  - Add unit that counters missile and expands game logic?

- initial map generation:
  - Omitting certain counters or land features from maps might give ticket to ride vibes where
    a different map significantly changes gameplay.
  - Omitting certain tech tree options is another way to randomize games.
  - Omitting certain tech tree options is a bit like randomized civ pick in AOE

- Tech Tree:
  - Can randomize what is availably per game for all players or everyone has the same random tree
  - Can randommize connections between the tech tree somewhat but keep is thematic.
  - Perhaps certain node connections should be consistent e.g. plane node always leads to (plane -> bomb -> heat seaking)

- Different terrains: 
  - Careful, this may be hard to implement due to pathfinding algorithms assuming same cost movement

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
4. Same cost Orgthogonal and Cardinal movement (Chebyshev distance):
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
9. Avoid quadtree optimization for as long as possible.
  - Seems more complex and error prone than other methods like spatial partitioning
  - Dynamic performance seems harder to tune and profile
  - Can't reuse existing hexagon logic.
10. Avoid GPU optimization for as long as possible.
  - I don't have any experience with this type of optimization.
  - I assume it may be platform dependent e.g. phone GPU is less powerful than computer integrated GPU
     is less powerful than computer dedicated GPU.
11. Use cardinal directions for hexagon points and sides (north, northeast, etc) because the order is more
    obvious e.g. northeast not eastnorth where as topRight and rightTop are both used.
12. Use {19,10} as initial hexagon offsets. These offsets greatly simplify getGridPos but makes getPixelPos
    slightly more complex. If hexagon configuration changes in the future, probably should prioritze
    simplifying getGridPos over getPixelPos.
13. Infantry is only simple ground unit. Less artwork, less animations, less design decisions.
    If needed, infantry can have upgrades which can helps tech tree idea.
  - sniper -> range upgrade
  - anti tank infantry -> anti tank ammunition
  - marine -> cross shallow water

---

## Game Design Questions

- Risk vs reward for how many units to produce and what type vs how production to invest elsewhere could be nice. 
  e.g. I invested too much, my enemy lines are breaking apart. Infantry units may need to be ranged in order for
  more units to matter. It might also be nice to allow bimodal behavior. E.g. Charge infantry vs infantry shoots
  from far away. This affects all infantry to keep game controls simple.
- Having different constraints on layers may help optimize. e.g. ground layer restricts unit movement to 1 pixel
  but low altitude and high altitude units allow more than 1 pixel movement. Can aslo sub tick only 1 layer
  and regular tick the others.
- Using 3 layers instead of 2 may help with spatial partitioning
- Flux looks visually nice because units in combat vs out of combat have different animations. In final game,
  perhaps there should be an attack animations for infantry and a non attack animation.
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

###  Brainstorm Spatial Partitioning

In general:

- Start simple, refactor and optimize later. I will most likely need to make significant changes 
  in the future no matter what I build. Get a frame work in place that unblocks other tasks like 
  local path finding and adding more units types.
- This may unblock performing more CPU intensive local path finding (e.g. theta*) for
  single hexagon of unit or local 7 hexagons. This might work much better than flux.
  Also, the unit can calculate for additional hexagons or add more more adjacent 
  hexagons as needed.

High level Implementation:

- Units store next and prior pointer. Hexagon should also store next and prior or maybe call it something else.
- Use double linked list that store pointers to relate units and next units.
- Not sure best way to break up player and layer. Probably best to start with a map
  with a key of (player + layer) and see if performance optimization is needed.
- Use existing hexagons as grid
  - coupling seems fine since hexagons are so fundemental to game.
  - most likely there is a way to reuse hexagon collision logic for units changing hexagons 
  - writing save file is easier if I reuse hexagon
  - hexagon has a soft limit of <20 or <=20 pixel for units. Seems fine for now and can be worked around.

Low level implementation:

- Go through issues to make sure they are all addressed.
- start per unit heap allocation during new unit production. 
- start with using getGridPos on every turn to figure out hexagon. There many ways to optimize this in future.
- Avoid two way checks where possible with tagging. Some two ways checks are easier to avoid than others.
- knowing how many units of a certain player are in a given hexagon 
  might be good for guided missile targets. Can easily increment and 
  decrement a counter on node rather than iterate through heap 
  allocations. Can calculate unit delta e.g. opponenet has 10 units, I have 2, so missile is +8 overall.

Careful of current issues:

- Will probably need to support checking a lot of hexagons for certain scenarios:
  - A large explosion will have collision > one hexagon. This is something I am concerned about
    because this will probablhy be the most performance sensitive computation since players
    will want to most likely target areas with alot of units.
  - A unit with large vision e.g. guided air missile will probably need to check multiple hexagons.
  - If I implement in the future units which capture a city atttempt to capture the next closest city
    to themself rather than to the city. It might be more on theme that they all attack the same city
    which is an easier calculation? 
- For nearest unit I can either check optimially by iterating through more distance hexagons or brute
  force check all hexagons.
- On unit kill need to handle its double linked list pointers
- Need to sync or verify or handle units switching layers if layer optimization is needed
- Units able to move more than one pos per turn cannot transfer hexagons via single line check.
  Need to use something more robust like current getGridPos.
- nodes switching ownership should not use newNode anymore unless pointers are preserved.

Potential future issues:

- Units largers than one hexagon have work arounds e.g. add unit to more than one hexagon.
- Same work around is possible for units moving fast.
- Any work arounds needed will probably be a pain but this issue is far in future.

Post Work Reflection:

- The implementation worked because CPU performance has improved multiple times:
  - collidesWithUnit became negligible and allowed ~1300 infantry units at 60 FPS.
  - analyzeTargets became negligible and allowed ~2500 infantry units with radius of 3 pixels at 60 FPS.
    Performance is understanted since it was tested with developer logging taking 50% of update game CPU usage.
- The implementation is simpler than what was designed
  - I didn't split layers. This was a good decision since I needed to start simple and this is an optimization.
  - Hexagon only stores a single pointer, not two which is simpler.
  - units per player in hex would require multiple counters. Only added a single total unit counter per hex.
  - Hexagons are separate from nodes, so node ownership update logic did not need an update.
- The implementation feels much more elegant than prior code. 
  - collidesWithUnit signature simplified from (newPos, unit, s.units) -> (newPos, unit)
  - Same updatePointers function is used for unit production and unit movement between hexagons.

---

### Brainstorm Implementation for Local Pathfinding with Spatial Partioning

Overall Thoughts:

- Start with an optimal solution. This can be used for few units. If this is too CPU intense, 
  can use large agent count algorithms or methods when more than a few units are present. 
- Probably want a default path finding algorithm for most units and specialized for infantry.
- Think about and define common player experiences and the desired behavior.
- First implementation goal is to allow better gridlock/clumping that forms lines of engagement efficiently.
- Units moving the same speed mean that only static units will collide. So don't perform
  expensive path find unless static units are nearby.
- Unfiorm cost grid greatly simplifies algorithms and helps reduce CPU usage. This is probably
  needed for optimal game experience since pathfinding is fundamental and dynamic movement cost is 
  a minor idea.
- Can use different path finding techniques based on how many local units are nearby. Probably
  don't want to change techniques based on total game units since it would add a lot of complexity
  and modality.
- Can use different path finding techniques based on unit type e.g. infantry will have high count
  and high number of collisions. Infantry needs optimized, low CPU, path finding which helps clumps
  them together. Bombers will be low count, low number of collisions.
- Might need to perform specific path finding for a few units, then large unit updates e.g. missile
  kills off several units, perform large agent count path finding.

Potential optimizations

- p100 runtime calcuation should be prioritzed. What can i compute at startup and during other
  game actions or turns to speed up computation?
- When a change occurs, only recalculate nodes affected 
- What heuristics can I use for better performacne e.g euclidean distance, octile distance, other distances, etc
- What optimizations does uniform grid cost allow?
- What optimizatoin does 1 pixel movemnt cost allow?
- JPS pruning gains (avoid duplicate paths via uniform movement, allow jumps)
- JPS* precompute gains (save runtime calculations to determine JPS nodes)

Implementation order:
- Implement globa theta*. Global is probably needed long term, seems easier to implement, small stepping stone.
- Implement local theta*. Use visibility grids for hexagons and units. Use hexagon spatial partitioning.
- Determine how "extra" units are handled? Do they pack in tightly? or sit back? 
- If conflicts are causing issues, determine if normal CBS vs enhanced CBS is better for handling conflicts.

---

### Brainstorm Multi-agent and Local Pathfinding

Pathfinding algorithm types:

- Static obstacles: 
  - Theta* (any angle, more natural paths) https://www.gameaipro.com/GameAIPro2/GameAIPro2_Chapter16_Theta_Star_for_Any-Angle_Pathfinding.pdf
  - A* (grid based movement) without post processing
  - A* with post processing
  - Visibility graphs (find shortest path, find collisions, find vertexs of collision, use these for paths)
- Dynamic obstacles: 
  - D* https://www.ri.cmu.edu/pub_files/pub3/stentz_anthony__tony__1994_2/stentz_anthony__tony__1994_2.pdf
  - D* focused
  - D* lite
  - Conflict Based Search
  - Enhanced Conflict Based Search (Not optimal, but faster computation) https://www.ijcai.org/Proceedings/15/Papers/110.pdf
  - Meta agent conflict based search (optimization for units which collide frequently with eachother)
  - Geometric MAPF or MAPF for large agents (LAMAPF) https://ojs.aaai.org/index.php/AAAI/article/view/4756 
    (Large agents are just not single point)
- Group:
  - flow fields
  - flocking behaviors
  - steering
  - dynamic obstacle avoidance
  - continuum crowds

Other considerations:

- Anonymous Multi-agent pathfinding: 
  - Set of targets to reach but specific agent does not matter
  - Might be useful for infantry collide together on front lines
- Nav mesh is being used via free movement hexagons vs obstacle hexagon.
- Use window searches aka use smaller searches to lower CPU usage
- If randomness is desired, Rapidly-exploring random tree can be used.

Probably don't want:

- life long planning since I am not using multiple objective plans. Multiple targets or types of 
  target (local, global and flux) might make this relevant.
  - LPA* Life long planning A*
  - LPCBS Lifelong planning conflict-based search
- ICT seems way too complex for minimal gains over a* (https://www.sciencedirect.com/science/article/pii/S0004370212001543)
- Continuum crowds. Large crowds moving alot seems rare e.g. missile blowing up a lot of units.
 
CBS in depth notes: https://www.sciencedirect.com/science/article/pii/S0004370214001386

- CBS usually attempts to get all units to destination while also minimizing the total cost.
- CBS can be one of two classes. Optimal or suboptimal. Large number of agents usually uses suboptimal. 
- Suboptimal can't guarantee the solution is optimal. 
- CBS can not group units or group units together based on how often they are colliding.
- Some algorithms allow units to leave the destiantion and return. I probably don't want this since units lock
  together during attack.
- CBS can priotize some units over other units. Or create a priority for units.
- Cost functions can be used to limit the total cost or the p100 amongst all agents
- CBS calculation can either be distributed or centralized on hosts. This game is centralized since its single computer.
- Units can move or idle per tick.
- Centralized solvers have three classes:
  - reduce MAPF to other well studied computer science problems
  - MAPF-specific sub-optimal solvers
  - MAPF-specific optimal solvers
- This paper focuses on optimal solvers
- Reduce MAPF to other well studied vs problems:
  - Examples are boolean satisfiability (SAT), integer linear programming (ILP), answer set programming (ASP)
  - Usually for low nuber of agents, solvers optimally, high CPU cost
- MAPF-specific sub-optimal solvers (Probably what I want to implement)
  - Examples are:
    - search based solvers (close to optimal) (Hierachical cooperate A* aka HCA*)
    - Rule based handle specific scenarios with specific movements. They run in polynomial time. (TASS, Push and SWAP).
    - Hybrid approaches combine search and rule e.g. large grid with specific pathways to move units thru
- MAPF optimal solvers 
  - usually search a global space that is a combination of all agents individual states
- Avoiding surplus nodes helps with speed up. Usually called node pruning.
- CBS uses two levels
  - The high level is called a constraint tree. It is a binary tree.
  - root has no constraints, child inherit parent constraints, add one contraint per agent.
  - Goal node when all set of paths for all agents have no conflicts.
  - high level does BFS
  - When a conflict is found between two agents, two children nodes are created with single a 
    constraint for each respective agent.
  - Nodes can extract constaints by traveling up the tree
  - The root node can have more than one agent but all other nodes have at most one agent and path.
  - conflicts of more than 2 agenst can be handled by generating more than two nodes or by generating
    deeper trees with only 2 children.
- CBS behavors poorly for strongly coupled units aka high rate of internal conflicts.
- MA-CBS tries to overcome this by merging agents into a meta agent.
- A* can sometimes outperform CBS. CBS sometimes outperforms A*. Generally CBS is worst when a lot of conflicts occur.

- Idea to limit CPU usage. What about token bucket and units have a priority queue.
  Only allow so much complex pathing to occur over a span of time.

- Flow fields
  - Using pixel pathing struct, I can check the 9 adjacent squares to find the direction of the nearest obstacle.
  - Can add flow fields based on nearby obstacles and sum thier fluxes to get a single final directon.
  - longer obstacles can have stronger fluxs so units don't start don't long path.
- Flocking
  - Three main types:
    - separation
    - alignment
    - cohesion
  - There are several other types available.
  - Can sum points for units within a hexagon or units within the nearest 7 hexagons
    - can sum units in different groups e.g. not moving since attacking, not moivng since blocked, etc
    - units can move towards or away from this point
- Streering
  - Can use spheres around non spherical objects for object avoidance e.g. hexagon has a sphere, if a unit
    collides with the sphere push it away. I could also create a sphere that grows bigger when X units
    are withn a hexagon. Object avoidance is different than collision detection.
  - Can use a sphere in front of a units direction to test for object avoidance
  - This includes path following and unaligned collision avoidance which may be needed for airplanes.
  - Can use steering behaviors to make infantry avoids tanks
  - Can prioritze steering behaviors e.g. do X unless 0. do Y if X is 0 and Y is not 0. Do Z if X and Y are 0.
  - Can merge streering behavior together to get combined direction and velocity.
- dynamic obstacle avoidance
  - can be used for non infantry units. It is probably CPU intensive.

Have read through several college publications and game ai pro articles not linked above:

- https://www.gameaipro.com/GameAIPro2/GameAIPro2_Chapter14_JPS_Plus_An_Extreme_A_Star_Speed_Optimization_for_Static_Uniform_Cost_Grids.pdf
- https://www.gameaipro.com/GameAIPro2/GameAIPro2_Chapter15_Subgoal_Graphs_for_Fast_Optimal_Pathfinding.pdf
- https://www.gameaipro.com/GameAIPro3/GameAIPro3_Chapter23_Faster_Dijkstra_Search_on_Uniform_Cost_Grids.pdf

---

### Brainstorm Unit List

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
    - High Priority:
      - tank
      - mobile SAM (shoots missile in low altitude or high altitude)
    - Low Priority:
      - mobile artillery / rocket artillery (thought is low priority because SAM pathing more complex. Could be wrong)
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