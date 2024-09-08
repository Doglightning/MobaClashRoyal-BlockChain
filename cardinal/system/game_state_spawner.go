package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/iterators"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "MobaClashRoyal/component"
	"MobaClashRoyal/msg"
)

// Spawns Game state for a new match.
// called by create_match.go msg.
func GameStateSpawnerSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage(world,
		func(create cardinal.TxData[msg.CreateMatchMsg]) (msg.CreateMatchResult, error) {
			//create filter for matching ID's
			matchFilter := cardinal.ComponentFilter(func(m comp.MatchId) bool {
				return m.MatchId == create.Msg.MatchID
			})

			// Search for existing matches.
			existingMatchSearch := cardinal.NewSearch().Entity(filter.Contains(filter.Component[comp.MatchId]())).Where(matchFilter)
			count, err := existingMatchSearch.Count(world)
			if err != nil {
				return msg.CreateMatchResult{Success: false}, fmt.Errorf("error during search (game_state_spawner.go): %w", err)
			}

			// No match found.
			if count == 0 {
				//Create new gamestate
				teamStateID, err := cardinal.Create(world,
					comp.MatchId{MatchId: create.Msg.MatchID},
					comp.UID{UID: 0},
					comp.Player1{
						Nickname:    create.Tx.PersonaTag,
						Hand:        []string{"Vampire", "ArcherLady", "Mage"},
						Deck:        []string{"ArcherLady", "Vampire"},
						RemovalList: make(map[int]bool),
						Gold:        5,
					},
					comp.SpatialHash{Cells: make(map[string]comp.SpatialCell),
						CellSize: SpatialGridCellSize,
						StartX:   float32(MapDataRegistry[create.Msg.MapName].StartX),
						StartY:   float32(MapDataRegistry[create.Msg.MapName].StartY)},
				)

				if err != nil {
					return msg.CreateMatchResult{Success: false}, fmt.Errorf("error creating match (game_state_spawner.go): %w", err)
				}

				// get spatial hash for collision map
				hash, err := cardinal.GetComponent[comp.SpatialHash](world, teamStateID)
				if err != nil {
					return msg.CreateMatchResult{Success: false}, fmt.Errorf("error getting hash component (game_state_spawner.go): %w", err)
				}

				//spawn bases
				err = spawnBasesGSS(world, create.Msg.MatchID, teamStateID, create.Msg.MapName, hash)
				if err != nil {
					return msg.CreateMatchResult{Success: false}, err
				}

				return msg.CreateMatchResult{Success: true}, nil // end logic for player1
			}

			// If a match is found, add the Player2 component to the existing entity
			matchFound, err := existingMatchSearch.First(world)

			if err != nil {
				return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component (game_state_spawner.go): %w", err)
			}

			//add player2 component
			err = cardinal.AddComponentTo[comp.Player2](world, matchFound)
			if err != nil {
				return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component (game_state_spawner.go): %w", err)
			}

			//set player2 compoenent
			err = cardinal.SetComponent(world, matchFound,
				&comp.Player2{
					Nickname:    create.Tx.PersonaTag,
					Hand:        []string{"Vampire", "ArcherLady", "Mage"},
					Deck:        []string{"ArcherLady", "Vampire"},
					RemovalList: make(map[int]bool),
					Gold:        5,
				})

			if err != nil {
				// if error remove the empty player2 component
				err = cardinal.RemoveComponentFrom[comp.Player2](world, matchFound)
				if err != nil {
					return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component (game_state_spawner.go): %w", err)
				}
				return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component (game_state_spawner.go): %w", err)
			}

			return msg.CreateMatchResult{Success: true}, nil

		})
}

// spawns bases and towers for both teams
func spawnBasesGSS(world cardinal.WorldContext, matchID string, teamStateID types.EntityID, mapName string, spatialHash *comp.SpatialHash) error {
	// check that map name exists in map registry
	if _, exsist := MapDataRegistry[mapName]; !exsist {
		return fmt.Errorf("map does not exist in MapDataRegistry (game_state_spawner.go/spawnBasesGSS)")
	}

	//get UID component
	uid, err := cardinal.GetComponent[comp.UID](world, teamStateID)
	if err != nil {
		return fmt.Errorf("error getting UID component (game_state_spawner.go/spawnBasesGSS): %w", err)
	}

	//spawn Blue Base
	baseID, err := cardinal.Create(world,
		comp.MatchId{MatchId: matchID},
		comp.UID{UID: uid.UID},
		comp.MapName{MapName: mapName},
		comp.UnitName{UnitName: "Base"},
		comp.Team{Team: "Blue"},
		comp.Health{CurrentHP: StructureDataRegistry["Base"].Health, MaxHP: StructureDataRegistry["Base"].Health},
		comp.Position{PositionVectorX: float32(MapDataRegistry[mapName].Bases[0][0]), PositionVectorY: float32(MapDataRegistry[mapName].Bases[0][1]), PositionVectorZ: float32(MapDataRegistry[mapName].Bases[0][2])},
		comp.UnitRadius{UnitRadius: StructureDataRegistry["Base"].Radius},
		comp.State{State: "Default"},
		comp.Attack{Combat: false, Damage: StructureDataRegistry["Base"].Damage, Rate: StructureDataRegistry["Base"].AttackRate, Frame: 0, DamageFrame: StructureDataRegistry["Base"].DamageFrame, Class: StructureDataRegistry["Base"].Class, AttackRadius: StructureDataRegistry["Base"].AttackRadius, AggroRadius: StructureDataRegistry["Base"].AggroRadius},
		comp.CenterOffset{CenterOffset: StructureDataRegistry["Base"].CenterOffset},
	)

	if err != nil {
		return fmt.Errorf("error creating blue base ((game_state_spawner.go/spawnBasesGSS)): %w", err)
	}

	//add structure to spatial hash collision map
	AddObjectSpatialHash(spatialHash, baseID, float32(MapDataRegistry[mapName].Bases[0][0]), float32(MapDataRegistry[mapName].Bases[0][1]), StructureDataRegistry["Base"].Radius, "Blue", "structure")

	//incriment UID
	uid.UID++

	//spawn Red Base
	baseID, err = cardinal.Create(world,
		comp.MatchId{MatchId: matchID},
		comp.UID{UID: uid.UID},
		comp.MapName{MapName: mapName},
		comp.UnitName{UnitName: "Base"},
		comp.Team{Team: "Red"},
		comp.Health{CurrentHP: StructureDataRegistry["Base"].Health, MaxHP: StructureDataRegistry["Base"].Health},
		comp.Position{PositionVectorX: float32(MapDataRegistry[mapName].Bases[1][0]), PositionVectorY: float32(MapDataRegistry[mapName].Bases[1][1]), PositionVectorZ: float32(MapDataRegistry[mapName].Bases[1][2])},
		comp.UnitRadius{UnitRadius: StructureDataRegistry["Base"].Radius},
		comp.State{State: "Default"},
		comp.Attack{Combat: false, Damage: StructureDataRegistry["Base"].Damage, Rate: StructureDataRegistry["Base"].AttackRate, Frame: 0, DamageFrame: StructureDataRegistry["Base"].DamageFrame, Class: StructureDataRegistry["Base"].Class, AttackRadius: StructureDataRegistry["Base"].AttackRadius, AggroRadius: StructureDataRegistry["Base"].AggroRadius},
		comp.CenterOffset{CenterOffset: StructureDataRegistry["Base"].CenterOffset},
	)

	if err != nil {
		return fmt.Errorf("error creating red base (team state spawner (game_state_spawner.go/spawnBasesGSS): %w", err)
	}

	//add structure to spatial hash collision map
	AddObjectSpatialHash(spatialHash, baseID, float32(MapDataRegistry[mapName].Bases[1][0]), float32(MapDataRegistry[mapName].Bases[1][1]), StructureDataRegistry["Base"].Radius, "Red", "structure")

	//incriment UID
	uid.UID++

	//spawn all towers
	for i := 0; i < MapDataRegistry[mapName].numTowers; i++ {
		//spawn Blue towers
		towerID, err := cardinal.Create(world,
			comp.MatchId{MatchId: matchID},
			comp.UID{UID: uid.UID},
			comp.MapName{MapName: mapName},
			comp.UnitName{UnitName: "Tower"},
			comp.Team{Team: "Blue"},
			comp.Health{CurrentHP: StructureDataRegistry["Tower"].Health, MaxHP: StructureDataRegistry["Tower"].Health},
			comp.Position{PositionVectorX: float32(MapDataRegistry[mapName].TowersBlue[i][0]), PositionVectorY: float32(MapDataRegistry[mapName].TowersBlue[i][1]), PositionVectorZ: float32(MapDataRegistry[mapName].TowersBlue[i][2])},
			comp.UnitRadius{UnitRadius: StructureDataRegistry["Tower"].Radius},
			comp.State{State: "Default"},
			comp.Attack{Combat: false, Damage: StructureDataRegistry["Tower"].Damage, Rate: StructureDataRegistry["Tower"].AttackRate, Frame: 0, DamageFrame: StructureDataRegistry["Tower"].DamageFrame, Class: StructureDataRegistry["Tower"].Class, AttackRadius: StructureDataRegistry["Tower"].AttackRadius, AggroRadius: StructureDataRegistry["Tower"].AggroRadius},
			comp.CenterOffset{CenterOffset: StructureDataRegistry["Tower"].CenterOffset},
		)

		if err != nil {
			return fmt.Errorf("error creating blue tower ((game_state_spawner.go/spawnBasesGSS)): %w", err)
		}

		//add structure to spatial hash collision map
		AddObjectSpatialHash(spatialHash, towerID, float32(MapDataRegistry[mapName].TowersBlue[i][0]), float32(MapDataRegistry[mapName].TowersBlue[i][1]), StructureDataRegistry["Tower"].Radius, "Blue", "structure")

		//incriment UID
		uid.UID++

		//spawn Blue towers
		towerID, err = cardinal.Create(world,
			comp.MatchId{MatchId: matchID},
			comp.UID{UID: uid.UID},
			comp.MapName{MapName: mapName},
			comp.UnitName{UnitName: "Tower"},
			comp.Team{Team: "Red"},
			comp.Health{CurrentHP: StructureDataRegistry["Tower"].Health, MaxHP: StructureDataRegistry["Tower"].Health},
			comp.Position{PositionVectorX: float32(MapDataRegistry[mapName].TowersRed[i][0]), PositionVectorY: float32(MapDataRegistry[mapName].TowersRed[i][1]), PositionVectorZ: float32(MapDataRegistry[mapName].TowersRed[i][2])},
			comp.UnitRadius{UnitRadius: StructureDataRegistry["Tower"].Radius},
			comp.State{State: "Default"},
			comp.Attack{Combat: false, Damage: StructureDataRegistry["Tower"].Damage, Rate: StructureDataRegistry["Tower"].AttackRate, Frame: 0, DamageFrame: StructureDataRegistry["Tower"].DamageFrame, Class: StructureDataRegistry["Tower"].Class, AttackRadius: StructureDataRegistry["Tower"].AttackRadius, AggroRadius: StructureDataRegistry["Tower"].AggroRadius},
			comp.CenterOffset{CenterOffset: StructureDataRegistry["Tower"].CenterOffset},
		)

		if err != nil {
			return fmt.Errorf("error creating red tower ((game_state_spawner.go/spawnBasesGSS)): %w", err)
		}

		//add structure to spatial hash collision map
		AddObjectSpatialHash(spatialHash, towerID, float32(MapDataRegistry[mapName].TowersRed[i][0]), float32(MapDataRegistry[mapName].TowersRed[i][1]), StructureDataRegistry["Tower"].Radius, "Red", "structure")

		//incriment UID
		uid.UID++

	}

	//set UID in game state
	if err := cardinal.SetComponent[comp.UID](world, teamStateID, uid); err != nil {
		return fmt.Errorf("error updating UID (team state spawner (game_state_spawner.go/spawnBasesGSS): %w", err)
	}

	return nil
}

func getNextUID(world cardinal.WorldContext, matchID string) (int, error) {

	//create filter for matchID to get game state
	matchFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
		return m.MatchId == matchID
	})

	gameStateSearch := cardinal.NewSearch().Entity(
		filter.Exact(GameStateFilters())).
		Where(matchFilter)
	//game state
	gameState, err := gameStateSearch.First(world)

	if err != nil {
		return 0, fmt.Errorf("error searching for match (getNextUID): %w", err)
	}

	if gameState == iterators.BadID {
		return 0, fmt.Errorf("no match found with ID or missing components (getNextUID): %s", matchID)
	}

	//get UID compoenent
	UID, err := cardinal.GetComponent[comp.UID](world, gameState)
	if err != nil {
		return 0, fmt.Errorf("error getting UID component (getNextUID): %w", err)
	}
	returnUID := UID.UID
	//increment UID
	UID.UID++
	//Set updated UID component
	err = cardinal.SetComponent(world, gameState, UID)
	if err != nil {
		return 0, fmt.Errorf("error setting UID component (getNextUID): %w", err)
	}

	return returnUID, nil
}

// fetches unit components needed for spatial hash removal
func getPlayerComponentsGSS(world cardinal.WorldContext, id types.EntityID) (p1 *comp.Player1, p2 *comp.Player2, err error) {
	//get player1 component
	p1, err = cardinal.GetComponent[comp.Player1](world, id)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving player1 component (getPlayerComponentsGSS): %s", err)
	}

	// get player2 components
	p2, err = cardinal.GetComponent[comp.Player2](world, id)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving player2 component (getPlayerComponentsGSS): %s", err)
	}

	return p1, p2, nil
}

// Returns the Game state
func getGameStateGSS(world cardinal.WorldContext, mID *comp.MatchId) (types.EntityID, error) {
	//get teamstate to get spatialhash tree
	teamFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
		return m.MatchId == mID.MatchId
	})
	foundTeam, err := cardinal.NewSearch().Entity(
		filter.Exact(GameStateFilters())).
		Where(teamFilter).First(world)

	if err != nil {

		fmt.Printf("error searching for match (game state spawner): %s", err)
		return foundTeam, err
	}

	if foundTeam == iterators.BadID { // Assuming cardinal.NoEntity represents no result found
		return foundTeam, fmt.Errorf("no match found with ID or missing components (game state spawner): %s", mID.MatchId)
	}
	return foundTeam, nil
}

// Returns the collision hash
func getCollisionHashGSS(world cardinal.WorldContext, mID *comp.MatchId) (*comp.SpatialHash, error) {
	//get teamstate to get spatialhash tree
	teamFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
		return m.MatchId == mID.MatchId
	})
	foundTeam, err := cardinal.NewSearch().Entity(
		filter.Exact(GameStateFilters())).
		Where(teamFilter).First(world)

	if err != nil {
		return nil, fmt.Errorf("error searching for match (game state spawner): %w", err)
	}

	if foundTeam == iterators.BadID { // Assuming cardinal.NoEntity represents no result found
		return nil, fmt.Errorf("no match found with ID or missing components (game state spawner): %s", mID.MatchId)
	}

	collisionHash, err := cardinal.GetComponent[comp.SpatialHash](world, foundTeam)
	if err != nil {
		return nil, fmt.Errorf("collision hash not found (game state spawner): %w", err)
	}
	return collisionHash, nil
}
