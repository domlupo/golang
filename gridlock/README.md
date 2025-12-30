Gridlock

A prototype for a simple war/4x game. Simple like mini metro or mini motorways. War game like Civ,
hearts of iron, 9 kings etc. (also other popular war games on Steamdb). Think about how this game
relates to civilization revolution, Catan, RISK, Starcraft desert strike, totally accurate battle
simulator, etc. Graphics may end up similar to the Gnorp analogue. Is there a mobile equivalent for
war/4x games?

## Top 3 Ranked Concerns : Risks

1. Design and iterate on Local path finding -> Most unknown but fundamental game engine decision.
2. Faster iteration from better inspection and debugging tooling. -> Don't waste engineering effort
   due to poor prioritization.
3. Playtesting MVP. -> Don't invest major engineering effort into wrong tasks.
  - Playtesting MVP
    - units:
      - fighter (F-22)
      - bomber (B-1 Lancer)
      - infantry
      - node ballistic missile
      - node cruise missile
      - node anti ballistic missile
      - node anti air missile (can reuse other missiles e.g. heat seeking, etc)
    - risk vs reward mechanics
      - produce infantry
      - produce advanced unit
      - research tech/upgrades
    - randomized tech tree
      - select 1 of 3 upgrades
      - comeback mechanics (skip, salvage, etc)
    - human vs computer
    - deploy to web
    - UI
       - two main options
      - bar at bottom of screen (star craft)
        - helps with CPU
      - small overlays (mini metro)
        - more units allowed?
    - user controls
    - art
    - Maybe
      - unit attack plans or attack patterns
      - randomized map generation
      - comeback mechanics
    - Don't include
      - menu or screen transitions
      - sound
      - save file

## Closed Decisions

Stale: 15

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
4. Same cost Orthogonal and Cardinal movement (Chebyshev distance):
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
  - The aesthetics of a vertical layout hexagonal grid is much nice than a square grid
    - The pattern is more pleasing.
    - The border is more pleasing.
    - More variety to connect grid nodes.
    - Multiple connected hexagon nodes is more pleasing.
    - It makes the game seem more organic rather than synthetic which match war theme.
    - Specifically, the same infantry lines are formed but are more pleasing with a hexagon
      background.
    - Hexagon is closer to a circle shape for easier art, hacky collision detection, etc
  - Pixel Position to Hexagon index calculation is not CPU intensive (floating point nor sqrt math
    needed)
  - Grid chunking only requires 6 adjacent rather than 8.
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
  - I assume it may be platform dependent e.g. phone GPU is less powerful than computer integrated
    GPU is less powerful than computer dedicated GPU.
11. Use cardinal directions for hexagon points and sides (north, northeast, etc) because the order
is more obvious e.g. northeast not eastnorth where as topRight and rightTop are both used.
12. Use {19,10} as initial hexagon offsets. These offsets greatly simplify getGridPos but makes
getPixelPos slightly more complex. If hexagon configuration changes in the future, probably should
prioritize simplifying getGridPos over getPixelPos.
13. Infantry is only simple ground unit. Less artwork, less animations, less design decisions.
Infantry are simple background visual but main visuals are air units above them. If needed, infantry
can have upgrades which can helps tech tree idea.
  - sniper -> range upgrade
  - marine -> cross shallow water
14. Use json for files (cache, saves, etc)
  - human readable
  - Go standard library support
  - Go standard library does not support toml or yaml
  - Avoids cross platform issues with sqlite
15. Use 3 layers (land, ordnance, aircraft) instead of only two layers (land, air)
  - 3 layers are easier to code and understand
  - aircraft is for long living air units (bombers, fighters, etc)
  - ordnance is for temporary air units (bombs, paratroopers, etc)
  - It is not the most realistic but I think it is fine:
    - brood war only uses a single air layer
    - Simplifies a lot of design questions
      - avoids "what units fly in what layer?"
      - avoids "Is fighter only unit able to switch layers?"
    - Avoids implementation work
      - avoids implement fighter switching layers
    - Makes art easier to implement so players can see units clearly
    - Makes layer concept easier for players to understand
16. No side to side camera movement
  - I still might add camera zoom out (think mini metro or mini motorways) during game start.
  - I don't think this is needed for the game experience I am trying to create.
  - Add a lot of complexity (programming, UI, player , etc) with no provable benefit.
17. Punt all land units besides infantry for initial release.
  - There is enough complexity, polish and work for only infantry and air units.
  - I don't think non infantry land units are required to create desired game loop.
  - Would increase complexity of infantry behavior a lot.
  - If the game is popular enough, it would be a good theme for an expansion.
  - It avoids distracted or overwhelmed players trying to closely follow 2 different altitudes.
18. Punt multiplayer for initial release
  - There is enough complexity, polish and work for single player.
  - Single player games can be very popular and successful.
  - If I ever implement, initial research shows webRTC is probably better but more complex than
    WebSocket.
19. Ballistic missiles only collide with anti ballistic missiles. They do not collide with each
other.
  - They need to different from cruise missiles to justify their game design investment and player
    use.
  - They are designed to fly in their own layer and have few counters.
  - Realistically they will only be countered by anti ballistic missiles.

## Game Design Questions

- It would be cool optional polish setting if over time the battle field slightly changed colors
  - frequently traveled hex changed colors -> can reuse hex linked list logic
  - plane explosions left burn marks
  - bombs left marks
- Should payload affect stealth? Stuff like this seems minor and hard to implement so probably not.
- Spy upgrade should allow blocking tech, lowering satellite launch success, etc
- Wind could be an interesting mechanic but unneeded. It could mess with aircraft.
- It might be interesting if the tech tree has a direct conflict via a race. e.g. your opponent can
  block you from researching something rather than them researching something for themselves. Or
  they can block 1 of the 3 options.
- I think slay the spire ascensions is a great way to adjust computer difficulty. Especially if I
  invest in all of these upgrades and tech trees. I can easily reuse them.
- A few ideas for a good comeback mechanic or high risk high reward play
  - higher risk, higher award unit upgrades
  - A bet mechanic where you can get a random upgrade for a discounted cost since its random
  - A skip or reroll mechanic to try and counter you opponent
  - A remove mechanic where you randomly get 2 of the remaining
  - upgrade salvage mechanic where you sell off tech to acquire other tech
  - node salvage mechanic where you sell off a node to acquire stuff
  - I feel like the tech tree might end up looking like a shop you can essentially barter with.
- I wonder if you should be able to see your opponents upgrades or if it should be hidden. If it is
  the same tech tree there could be a way to track what the player picked. It add an element of
  tracking to the game. I also can add an upgrade or unit where you can see more of your opponents
  upgrades. It could have a nice UI. There can also be an upgrade to see further in the tech tree or
  see what upgrade count your opponent is at.
- I like the look of mosaics and mosaic the word sounds cool. The hexagon grid pattern reminds me of
  mosaics a lot. Should I incorporate this in this game or a future game? Should I use mosaics to
  inspire art direction over other games like Catan, Risk and Gnorp analogue? Maybe polish is mosaic
  filter.
- I like how Balatro has a skip mechanic. If i make a tech tree which is randomized, a select 1 of 3
  can be used. Then I can allow either cheap skips to progress further or a reroll option. Initial
  thought is skip would be more interesting than reroll and more skill based.
- Risk vs reward for how many units to produce and what type vs how production to invest elsewhere
  could be nice. e.g. I invested too much, my enemy lines are breaking apart. Infantry units may
  need to be ranged in order for more units to matter. It might also be nice to allow bimodal
  behavior. E.g. Charge infantry vs infantry shoots from far away. This affects all infantry to keep
  game controls simple.
- It looks visually nice when units in combat vs out of combat have different animations. It also
  clearly shows game state. In final game, perhaps there should be an attack animations for infantry
  and a non attack animation.
- Should air units have impassable terrain e.g. high mountain? Will make path finding harder for air
  units
- Stopping node unit production in order to advance tech tree, improve node is probably a good form
  of risk -> reward.
- Should nodes support unit cycling e.g. produce missile, then produce infantry, then missile etc
- Chemical gas can push enemies back unless they have gas masks
- Should nodes have a max unit capacity e.g. This node stops producing since 10 units it produced
  are still alive
- Should bombs dropped from planes have health and can die?
- Should paratroopers be able to be shot at first?
- Plane collisions would be awesome?
- air units can produce other units e.g. missiles, paratroops, etc?
- Should units besides explosives detonate?
  - e.g. A plane being hit could explode damaging nearby units. Can use unit explode radius variable
    for this. A plane can fall from sky and cause ground damage from debris.
  - A missile detonating should detonate other missiles in air
  - This could create interesting explosion chains
- Should there be a zoom out camera:
  - Think mini metro and mini motorways style that slowly zooms out to the full area.
- Should the focus of the game be positioning, logistics and decisions rather than micro?
  - It seems interesting whens units get into long term fights
  - Reminds me of the board game go RISK and Catan
  - Reminds me of the military diagrams that show battles over time
  - May be able to utilize the idea of target change units
  - Early fight for position, long term fight for ?
- top down graphics or isometric or 3d?
- Players get to pick starting positions?
- Hidden or fully visible grid?
- APM lock on nodes or units so you cannot spam actions? Think about Sokpop game clickyland where
  you are limited by clicks.
- Campaign mode. Specific battles, war or countries? Similar to cities in mini metro and mini
  motorway.
- Tech tree?
- Special nodes?
- Should nodes only have a single grid offset (offset = 0 is acceptable) or should offset differ?
- Node expansion overtime?
- Grids expand over time to cities, cities can merge together, merged cities are more powerful?
- Buildings?
- Do I want to stand out from Civ, or be a simpler Civ e.g. Factorio -> Shapez. How does this affect
  marketing?
- Should players be able to target their own nodes?
- Should a node be able to target itself?
- Cities have different units, different unit ratios, or different counts produced? e.g. Catan
- Traders or cargo can go to a node which you own to transfer production of city production
- Nodes with defense system
- How is node unit production handled?
- Theme is the Gnorp analogue?
- Game name?
  - choke point
  - choke points
  - gridlock
  - gridlocked
  - hawk nest
  - eagle eye
  - iron wing
  - sky claw
  - talon strike
  - hawk point
- Replay feature where you can see a battle recap?
- Should unit behavior be toggleable on a global scale? E.g. bombers have different bomb plan options
- Should nodes produce multiple units e.g. Nodes always create infantry, then get upgraded over time.
  - Initially, this seems awesome but makes unit production harder
  - Capturing nodes has more severe swings? which is good?
  - Risk vs Reward?
  - Can add a sabotage mechanic.
- Missile upgrade can guarantee different areas are targeted or not. Right now missile sends out
  multiple missile which hit the same area since it doesn't factor in it will kill units. Another
  upgrade could be missile coordination e.g. multiple noes coordinate.
- Should obstacles be higher resolution than hexagonal? e.g. allow non node hexagons to have partial
  obstacles.
- It would be good polish to have end game stats/summary
  - infantry killed
  - infantry died
  - node control
  - units with the most kills
    - a single bomber had X kills
- It would be good polish to allow players to configure certain aspects e.g. highlight widths
- Should nodes support multiple bombers? If yes should there also be options
  - fly in formation and bomb same targets
  - fly separate but bomb same targets
  - fly separate and bomb differently (complex due to path finding)
- A "heat" seeking mechanic can be added to several different flying units to prioritize targets
  - counter missile
  - bomber / fighter has heat seeking missiles
  - flares
- There should probably be some type of comeback mechanic.
  - As players lose nodes how do you make nodes more valuable?
    - easier to control fewer nodes
    - If nodes are able to be upgraded, most likely better node are in the back. This means it is likely a
      bomber would have more frequent trips for the losing player than the winning player if the
      battle field lines are close.
    - What other comeback mechanics can be added?
- If I add a unit like paratrooper, Hex control probably needs to remain. I can make hex control be
  less micro able by making a "war zone area" between opposing nodes which are not controllable.
- Fog of war? It may be needed or it might be very beneficial while also being simple to implement.
  A compromise could be to keep only certain info hidden then something like UAV or satellite can
  expose the info.
- Tech tree could be a source of randomness. The choices and order can vary per game. Players can
  have the same or different pick 1 of three that slowly goes through the tech tree. I like 1 of 3
  because its like rock paper scissors.
- fuel mechanic for planes could be risk vs reward for building closer airports but this already
  exists because planes need to rearm.
- Should infantry have upgrades like gas masks, cross shallow water? Cross shallow water is hard
  since it is dynamic pathing.

## Misc Notes

- Research Civ concepts and features e.g. culture victory, trade, etc
- Research Polytopia game
- Research other mobile platform "Civ" games like Unciv

## Brainstorms

### Brainstorm Spatial Partitioning

In general:

- Start simple, refactor and optimize later. I will most likely need to make significant changes in
  the future no matter what I build. Get a frame work in place that unblocks other tasks like local
  path finding and adding more units types.
- This may unblock performing more CPU intensive local path finding (e.g. theta*) for single hexagon
  of unit or local 7 hexagons. Also, the unit can calculate for additional hexagons or add more more
  adjacent hexagons as needed.

High level Implementation:

- Units store next and prior pointer. Hexagon should also store next and prior or maybe call it
  something else.
- Use double linked list that store pointers to relate units and next units.
- Not sure best way to break up player and layer. Probably best to start with a map with a key of
  (player + layer) and see if performance optimization is needed.
- Use existing hexagons as grid
  - coupling seems fine since hexagons are so fundamental to game.
  - most likely there is a way to reuse hexagon collision logic for units changing hexagons
  - writing save file is easier if I reuse hexagon
  - hexagon has a soft limit of <20 or <=20 pixel for units. Seems fine for now and can be worked
    around.

Low level implementation:

- Go through issues to make sure they are all addressed.
- start per unit heap allocation during new unit production.
- start with using getGridPos on every turn to figure out hexagon. There many ways to optimize this
  in future.
- Avoid two way checks where possible with tagging. Some two ways checks are easier to avoid than
  others.
- knowing how many units of a certain player are in a given hexagon might be good for guided missile
  targets. Can easily increment and decrement a counter on node rather than iterate through heap
  allocations. Can calculate unit delta e.g. opponent has 10 units, I have 2, so missile is +8
  overall.

Careful of current issues:

- Will probably need to support checking a lot of hexagons for certain scenarios:
  - A large explosion will have collision > one hexagon. This is something I am concerned about
    because this will probably be the most performance sensitive computation since players will want
    to most likely target areas with a lot of units.
  - A unit with large vision e.g. anti air missile will probably need to check multiple hexagons.
  - If I implement in the future units which capture a city attempt to capture the next closest city
    to them self rather than to the city. It might be more on theme that they all attack the same
    city which is an easier calculation?
- For nearest unit I can either check optimally iterating through more distance hexagons or
  brute force check all hexagons.
- On unit kill need to handle its double linked list pointers
- Need to sync or verify or handle units switching layers if layer optimization is needed
- Units able to move more than one pos per turn cannot transfer hexagons via single line check. Need
  to use something more robust like current getGridPos.
- nodes switching ownership should not use newNode anymore unless pointers are preserved.

Potential future issues:

- Units larger than one hexagon have work around e.g. add unit to more than one hexagon.
- Same work around is possible for units moving fast.
- Any work around needed will probably be a pain but this issue is far in future.

Post Work Reflection:

- The implementation worked because CPU performance has improved multiple times:
  - collidesWithUnit became negligible and allowed ~1300 infantry units at 60 FPS.
  - analyzeTargets became negligible and allowed ~2500 infantry units with radius of 3 pixels at 60
    FPS. Performance is understated since it was tested with developer logging taking 50% of update
    game CPU usage.
- The implementation is simpler than what was designed
  - I didn't split layers. This was a good decision since I needed to start simple and this is an
    optimization.
  - Hexagon only stores a single pointer, not two which is simpler.
  - units per player in hex would require multiple counters. Only added a single total unit counter
    per hex.
  - Hexagons are separate from nodes, so node ownership update logic did not need an update.
- The implementation feels much more elegant than prior code.
  - collidesWithUnit signature simplified from (newPos, unit, s.units) -> (newPos, unit)
  - Same updatePointers function is used for unit production and unit movement between hexagons.

### Brainstorm Implementation for Local Pathfinding with Spatial Partitioning

Overall Thoughts:

- Start with an optimal solution. This can be used for few units. If this is too CPU intense, can
  use large agent count algorithms or methods when more than a few units are present.
- Probably want a default path finding algorithm for most units and specialized for infantry.
- Think about and define common player experiences and the desired behavior.
- First implementation goal is to allow better gridlock/clumping that forms lines of engagement
  efficiently.
- Units moving the same speed mean that only static units will collide. So don't perform expensive
  path find unless static units are nearby.
- Uniform cost grid greatly simplifies algorithms and helps reduce CPU usage. This is probably
  needed for optimal game experience since pathfinding is fundamental and dynamic movement cost is a
  minor idea.
- Can use different path finding techniques based on how many local units are nearby. Probably don't
  want to change techniques based on total game units since it would add a lot of complexity and
  modality.
- Can use different path finding techniques based on unit type e.g. infantry will have high count
  and high number of collisions. Infantry needs optimized, low CPU, path finding which helps clumps
  them together. Bombers will be low count, low number of collisions.
- Might need to perform specific path finding for a few units, then large unit updates e.g. missile
  kills off several units, perform large agent count path finding.

Potential optimizations

- p100 runtime calculation should be prioritized. What can i compute at startup and during other
  game actions or turns to speed up computation?
- When a change occurs, only recalculate nodes affected
- What heuristics can I use for better performance e.g euclidean distance, octile distance, other
  distances, etc
- What optimizations does uniform grid cost allow?
- What optimization  1 pixel movement cost allow?
- JPS pruning gains (avoid duplicate paths via uniform movement, allow jumps)
- JPS* precompute gains (save runtime calculations to determine JPS nodes)

Implementation order:
- Implement global theta*. Global is probably needed long term, seems easier to implement, small
  stepping stone.
- Implement local theta*. Use visibility grids for hexagons and units. Use hexagon spatial
  partitioning.
- Determine how "extra" units are handled? Do they pack in tightly? or sit back?
- If conflicts are causing issues, determine if normal CBS vs enhanced CBS is better for handling
  conflicts.

### Brainstorm Multi-agent and Local Pathfinding

Pathfinding algorithm types:

- Static obstacles:
  - Theta* (any angle, more natural paths)
    https://www.gameaipro.com/GameAIPro2/GameAIPro2_Chapter16_Theta_Star_for_Any-Angle_Pathfinding.pdf
  - A* (grid based movement) without post processing
  - A* with post processing
  - Visibility graphs (find shortest path, find collisions, find vertices of collision, use these
    for paths)
- Dynamic obstacles:
  - D*
    https://www.ri.cmu.edu/pub_files/pub3/stentz_anthony__tony__1994_2/stentz_anthony__tony__1994_2.pdf
  - D* focused
  - D* lite
  - Conflict Based Search
  - Enhanced Conflict Based Search (Not optimal, but faster computation)
    https://www.ijcai.org/Proceedings/15/Papers/110.pdf
  - Meta agent conflict based search (optimization for units which collide frequently with
    eachother)
  - Geometric MAPF or MAPF for large agents (LAMAPF)
    https://ojs.aaai.org/index.php/AAAI/article/view/4756 (Large agents are just not single point)
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
  target (local, global, etc) might make this relevant.
  - LPA* Life long planning A*
  - LPCBS Lifelong planning conflict-based search
- ICT seems way too complex for minimal gains over a*
  (https://www.sciencedirect.com/science/article/pii/S0004370212001543)
- Continuum crowds. Large crowds moving alot seems rare e.g. missile blowing up a lot of units.

CBS in depth notes: https://www.sciencedirect.com/science/article/pii/S0004370214001386

- CBS usually attempts to get all units to destination while also minimizing the total cost.
- CBS can be one of two classes. Optimal or suboptimal. Large number of agents usually uses
  suboptimal.
- Suboptimal can't guarantee the solution is optimal.
- CBS can not group units or group units together based on how often they are colliding.
- Some algorithms allow units to leave the destination and return. I probably don't want this since
  units lock together during attack.
- CBS can prioritize some units over other units. Or create a priority for units.
- Cost functions can be used to limit the total cost or the p100 amongst all agents
- CBS calculation can either be distributed or centralized on hosts. This game is centralized since
  its single computer.
- Units can move or idle per tick.
- Centralized solvers have three classes:
  - reduce MAPF to other well studied computer science problems
  - MAPF-specific sub-optimal solvers
  - MAPF-specific optimal solvers
- This paper focuses on optimal solvers
- Reduce MAPF to other well studied vs problems:
  - Examples are boolean satisfiability (SAT), integer linear programming (ILP), answer set
    programming (ASP)
  - Usually for low number of agents, solvers optimally, high CPU cost
- MAPF-specific sub-optimal solvers (Probably what I want to implement)
  - Examples are:
    - search based solvers (close to optimal) (Hierarchical cooperate A* aka HCA*)
    - Rule based handle specific scenarios with specific movements. They run in polynomial time.
      (TASS, Push and SWAP).
    - Hybrid approaches combine search and rule e.g. large grid with specific pathways to move units
      thru
- MAPF optimal solvers
  - usually search a global space that is a combination of all agents individual states
- Avoiding surplus nodes helps with speed up. Usually called node pruning.
- CBS uses two levels
  - The high level is called a constraint tree. It is a binary tree.
  - root has no constraints, child inherit parent constraints, add one constraint per agent.
  - Goal node when all set of paths for all agents have no conflicts.
  - high level does BFS
  - When a conflict is found between two agents, two children nodes are created with single a
    constraint for each respective agent.
  - Nodes can extract constraints by traveling up the tree
  - The root node can have more than one agent but all other nodes have at most one agent and path.
  - conflicts of more than 2 agenst can be handled by generating more than two nodes or by
    generating deeper trees with only 2 children.
- CBS behaviors poorly for strongly coupled units aka high rate of internal conflicts.
- MA-CBS tries to overcome this by merging agents into a meta agent.
- A* can sometimes outperform CBS. CBS sometimes outperforms A*. Generally CBS is worst when a lot
  of conflicts occur.

- Idea to limit CPU usage. What about token bucket and units have a priority queue. Only allow so
  much complex pathing to occur over a span of time.

- Flow fields
  - Using obstacle info struct, I can check the 9 adjacent squares to find the direction of the
    nearest obstacle.
  - Can add flow fields based on nearby obstacles and sum their forces to get a single final
    direction.
  - longer obstacles can have stronger forces so units don't start don't long path.
- Flocking
  - Three main types:
    - separation
    - alignment
    - cohesion
  - There are several other types available.
  - Can sum points for units within a hexagon or units within the nearest 7 hexagons
    - can sum units in different groups e.g. not moving since attacking, not moivng since blocked,
      etc
    - units can move towards or away from this point
- Streering
  - Can use spheres around non spherical objects for object avoidance e.g. hexagon has a sphere, if
    a unit collides with the sphere push it away. I could also create a sphere that grows bigger
    when X units are withn a hexagon. Object avoidance is different than collision detection.
  - Can use a sphere in front of a units direction to test for object avoidance
  - This includes path following and unaligned collision avoidance which may be needed for
    airplanes.
  - Can prioritize steering behaviors e.g. do X unless 0. do Y if X is 0 and Y is not 0. Do Z if X
    and Y are 0.
  - Can merge streering behavior together to get combined direction and velocity.
- dynamic obstacle avoidance
  - can be used for non infantry units. It is probably CPU intensive.

Have read through several college publications and game ai pro articles not linked above:

- https://www.gameaipro.com/GameAIPro2/GameAIPro2_Chapter14_JPS_Plus_An_Extreme_A_Star_Speed_Optimization_for_Static_Uniform_Cost_Grids.pdf
- https://www.gameaipro.com/GameAIPro2/GameAIPro2_Chapter15_Subgoal_Graphs_for_Fast_Optimal_Pathfinding.pdf
- https://www.gameaipro.com/GameAIPro3/GameAIPro3_Chapter23_Faster_Dijkstra_Search_on_Uniform_Cost_Grids.pdf

### Brainstorm Unit List

- Questions:
  - Which units will counter other units?
  - Should units have upgrades?
  - Should different players have different tech trees? Different unit costs?
  - Should planes need to fly back to airport to refuel?
  - Should planes need to fly back to airport to rearm?
  - Should planes need to fly back to airport to repair?
  - Should plane maintenance be much quicker than building a new plane?
  - How can other units be unique and fit with node mechanic?
  - Should I do purpose exercise with units like I do with checkpoints
  - Is there a good way to counter artillery?
  - Heat seeking missile (it can be an upgrade but does not have to be)
  - Jet heat flare (it can be an upgrade but does not have to be)
  - Missile upgrade makes it smaller?
  - Prior to the game and during game allow players to tune unit configs? e.g. 1 SAM in army send 1
    rocket, 2 SAM in army send 1 rocket from all in range or send 2 rockets from closest? Army could
    have seensible defaults but then players can customize?

- Units:
  - Land:
    - infantry
  - Air
    - bomber
    - fighter
    - paratroopers
    - ac130 gunship
    - apache
    - stealth bomber
  - Node
    - missile
    - counter missile
  - Other

- Units to consider:
  - satellite
  - rocket
  - drone
  - uav
  - air refueler
  - Growler (electronic warfare support)
  - black hawk (stealth helicopter)
  - nuke (long time to build, but devastating?)
  - spy or espionage

### Brainstorm Bomber Unit

- size
  - what basic size should it be?
  - relative scale to other art (nodes, infantry, etc)
  - collision shape (I don't think circle will work here long term)
  - turn radius size
- behavior
  - is it going to be harder to target corners?
  - is it going to be harder to target edges?
  - When it crashes does it do ground damage? I think it should.
  - Does crashing with different payload amount matters? I think it should.
  - Should be less maneuverable than a fighter
- control
  - multi click for path to bombing run
    - highlight cone shows potential next clicks
  - click + drag + release for bomb run.
    - the bombs auto space out based on the number
  - Then multi click return path same as prior control
  - This means only one bomb run can happen. This seems fine. It is complex and hard to allow
    multiple bomb runs to occur since a future run would change the auto targets of the prior run.
  - Should bomber path and change bomb target after it has taken off be viable?
  - Need to display current and future path irrespective off how often update occurs.
- targeting
  - X bombs, dragging shows path
- programming
  - bomb as a separate unit?
  - bomb created at take off or during destruction / bomb drop?
  - multiple types of bomb? (should this change flight path?)
    - gravity (drops right below aircraft)
    - glide (drop at angle? I am not sure the benefit here unless multiple drop a)
    - guided (JDAM)
    - missile
  - targeting is much different than other units will help understand long term goal here.
- upgrades?
  - bomb types
  - speed
  - bomb count
  - turn radius
  - faster bomb drop
  - (should bomb accuracy not be exact, but upgradable to be better?)
  - fuel (longer distance). This is hard because targeting UI becomes more complex and can fail.
- bomb type:
  - gravity bombs
  - guided bombs (drop directly over target and steer towards units or something else)
  - bunker buster (hits node?, seems like a good high risk high reward option)
  - chemical (infantry gas masks counter)
  - incendiary (fire, can reuse plane crash fire)
  - glide bombs (fired from range and fly towards target)
  - nuclear bomb


### Brainstorm upgrades

There are lots of open questions for upgrades:
  - do upgrades make sense to add to game?
  - Upgrades would most likely be the randomized tech tree.
    - each player has the same tree (rock vs paper vs scissors)
    - More complex than rock vs papyer vs scissors because future tech can be shown?
    - More complex than rock vs papyer vs scissors because skips?
    - select 1 of 3
  - What is an upgrade vs what is a config? Certain things are not always better
  - Upgrades can make units perform better e.g. more range, more damage
  - Upgrades can make play easier e.g. manual flares vs auto flares
  - Should all or only some upgrades be toggleable?
  - what does the upgrade UI look like.
  - upgrades will probably need prereq
    - need to research heat seeking missile before adding count
    - need to research collision avoidance before fly in formation

### Brainstorm upgrades for fighter unit

Listing out potential upgrades for a single unit might help:

- Better vision / radar
- Faster fighter
- Faster projectile speed
  - missiles
  - rotary cannon?
  - bomb step delay
- Faster fire rate
  - basic missile
  - heat seeking missile
  - all other missile types
  - rotary cannon?
- Research ammo type
  - basic missile
  - semi active radar missile (fighter provides radar guidance)
  - active radar missile (fire and forget)
  - heat seeking missile (fire and forget)
  - bomb
  - flares
  - chaff
  - radar jam
  - rotary cannon
- Add ammo
  - Add missile count
  - Add heat seeking missile count
  - Add bomb count
  - Add flare
- Add way point
- Faster maintenance / reloads
- Faster build time
- Improve missile accuracy
- Improve missile damage Smart targeting?
  - Only attempt to intercept ballistic when it will reach in time
  - Only attempt to intercept cruise when it will reach in time
  - etc
- Steer away from:
  - colliding with other aircraft
  - opposing ammo (missiles)
- Add fighter count (This should work for any advanced unit combination e.g. fighter + bomber,
  fighter + ac130)
  - fly together in formation
  - fly separately in rotation

### Pixel Art Notes

- b1 lancer bomber
  - https://en.wikipedia.org/wiki/Rockwell_B-1_Lancer#Specifications_(B-1B)
  - b1 lancer is 146' long and 40 pixels -> 3.65
  - b1 lancer is 137' wide and 38 pixels -> 3.6
- f22 fighter
  - https://en.wikipedia.org/wiki/Lockheed_Martin_F-22_Raptor#Specifications_(F-22A)
  - f22 is 62' long -> 17
  - f22 is 44' wide -> 12
