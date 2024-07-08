package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "MobaClashRoyal/component"
	"MobaClashRoyal/msg"
)

//var SpatialGridCellSize = 300

// Spawns a new match entity.
func GameStateSpawnerSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.CreateMatchMsg, msg.CreateMatchResult](
		world,
		//iterate over every create match msgs
		func(create cardinal.TxData[msg.CreateMatchMsg]) (msg.CreateMatchResult, error) {
			//create filter for matching ID's
			matchFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
				return m.MatchId == create.Msg.MatchID
			})

			// Search for existing matches
			existingMatchSearch := cardinal.NewSearch().Entity(filter.Contains(filter.Component[comp.MatchId]())).Where(matchFilter)
			count, err := existingMatchSearch.Count(world)
			if err != nil {
				return msg.CreateMatchResult{Success: false}, fmt.Errorf("error during search: %w", err)
			}

			if count > 0 {
				// If a match is found, add the Player2 component to the existing entity
				matchFound, err := existingMatchSearch.First(world)

				if err != nil {
					return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component: %w", err)
				}

				err = cardinal.AddComponentTo[comp.Player2](world, matchFound)
				if err != nil {
					return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component: %w", err)
				}

				err = cardinal.SetComponent(world, matchFound, &comp.Player2{Nickname2: create.Tx.PersonaTag, RemovalList: make(map[int]bool)})
				if err != nil {
					err = cardinal.RemoveComponentFrom[comp.Player2](world, matchFound)
					if err != nil {
						return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component: %w", err)
					}
					return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component: %w", err)
				}

				return msg.CreateMatchResult{Success: true}, nil
			}

			//create a new match
			teamStateID, err := cardinal.Create(world, comp.MatchId{MatchId: create.Msg.MatchID}, comp.UID{UID: 0}, comp.Player1{Nickname: create.Tx.PersonaTag, RemovalList: make(map[int]bool)}, comp.SpatialHash{Cells: make(map[string]comp.SpatialCell), CellSize: SpatialGridCellSize, StartX: float32(MapDataRegistry[create.Msg.MapName].StartX), StartY: float32(MapDataRegistry[create.Msg.MapName].StartY)})
			if err != nil {
				return msg.CreateMatchResult{Success: false}, fmt.Errorf("error creating match: %w", err)
			}

			hash, err := cardinal.GetComponent[comp.SpatialHash](world, teamStateID)
			if err != nil {
				return msg.CreateMatchResult{Success: false}, fmt.Errorf("error getting hash component (team state spawner): %w", err)
			}

			//spawn bases
			err = spawnBasesTSS(world, create.Msg.MatchID, teamStateID, create.Msg.MapName, hash)
			if err != nil {
				return msg.CreateMatchResult{Success: false}, err
			}

			return msg.CreateMatchResult{Success: true}, nil
		})
}

// spawns bases and towers for both teams
func spawnBasesTSS(world cardinal.WorldContext, matchID string, teamStateID types.EntityID, mapName string, spatialHash *comp.SpatialHash) error {
	//get UID
	uid, err := cardinal.GetComponent[comp.UID](world, teamStateID)
	if err != nil {
		return fmt.Errorf("error getting UID component (team state spawner): %w", err)
	}

	if _, exsist := MapDataRegistry[mapName]; !exsist {
		return fmt.Errorf("map does not exsist in MapDataRegistry (team state spawner): %w", err)
	}

	//spawn Blue Base
	baseID, err := cardinal.Create(world,
		comp.MatchId{MatchId: matchID},
		comp.UID{UID: uid.UID},
		comp.MapName{MapName: mapName},
		comp.UnitName{UnitName: "Base"},
		comp.Team{Team: "Blue"},
		comp.UnitHealth{CurrentHP: StructureDataRegistry["Base"].Health, MaxHP: StructureDataRegistry["Base"].Health},
		comp.Position{PositionVectorX: float32(MapDataRegistry[mapName].Bases[0][0]), PositionVectorY: float32(MapDataRegistry[mapName].Bases[0][1])},
		comp.UnitRadius{UnitRadius: StructureDataRegistry["Base"].Radius},
	)

	if err != nil {
		return fmt.Errorf("error creating blue base (team state spawner): %w", err)
	}

	//add unit to spatial hash collision map
	AddObjectSpatialHash(spatialHash, baseID, float32(MapDataRegistry[mapName].Bases[0][0]), float32(MapDataRegistry[mapName].Bases[0][1]), StructureDataRegistry["Base"].Radius, "Blue")

	//incriment UID
	uid.UID++

	//spawn Red Base
	baseID, err = cardinal.Create(world,
		comp.MatchId{MatchId: matchID},
		comp.UID{UID: uid.UID},
		comp.MapName{MapName: mapName},
		comp.UnitName{UnitName: "Base"},
		comp.Team{Team: "Red"},
		comp.UnitHealth{CurrentHP: StructureDataRegistry["Base"].Health, MaxHP: StructureDataRegistry["Base"].Health},
		comp.Position{PositionVectorX: float32(MapDataRegistry[mapName].Bases[1][0]), PositionVectorY: float32(MapDataRegistry[mapName].Bases[1][1])},
		comp.UnitRadius{UnitRadius: StructureDataRegistry["Base"].Radius},
	)

	if err != nil {
		return fmt.Errorf("error creating red base (team state spawner): %w", err)
	}

	//add unit to spatial hash collision map
	AddObjectSpatialHash(spatialHash, baseID, float32(MapDataRegistry[mapName].Bases[1][0]), float32(MapDataRegistry[mapName].Bases[1][1]), StructureDataRegistry["Base"].Radius, "Red")

	//incriment UID
	uid.UID++

	if err := cardinal.SetComponent[comp.UID](world, teamStateID, uid); err != nil {
		return fmt.Errorf("error updating UID (team state spawner): %w", err)
	}

	return nil
}
