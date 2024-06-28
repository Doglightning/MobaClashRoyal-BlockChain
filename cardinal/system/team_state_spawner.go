package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/iterators"
	"pkg.world.dev/world-engine/cardinal/search/filter"

	comp "MobaClashRoyal/component"
	"MobaClashRoyal/msg"
)

// Spawns a new match entity.
func MatchSpawnerSystem(world cardinal.WorldContext) error {
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
				founds, err := existingMatchSearch.First(world)

				if err != nil {
					return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component: %w", err)
				}

				err = cardinal.AddComponentTo[comp.Player2](world, founds)
				if err != nil {
					return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component: %w", err)
				}

				err = cardinal.SetComponent(world, founds, &comp.Player2{Nickname2: create.Tx.PersonaTag})
				if err != nil {
					err = cardinal.RemoveComponentFrom[comp.Player2](world, founds)
					if err != nil {
						return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component: %w", err)
					}
					return msg.CreateMatchResult{Success: false}, fmt.Errorf("error adding Player2 component: %w", err)
				}

				return msg.CreateMatchResult{Success: true}, nil
			}

			mapFilter := cardinal.ComponentFilter[comp.MapName](func(m comp.MapName) bool {
				return m.MapName == create.Msg.MapName
			})

			mapSearch := cardinal.NewSearch().Entity(filter.Exact(filter.Component[comp.DirectionMap](), filter.Component[comp.GridUtils](), filter.Component[comp.MapName]())).Where(mapFilter)

			foundMap, err := mapSearch.First(world)

			if err != nil {
				return msg.CreateMatchResult{Success: false}, fmt.Errorf("error searching for map (team state spawner): %w", err)
			}

			if foundMap == iterators.BadID { // Assuming cardinal.NoEntity represents no result found
				return msg.CreateMatchResult{Success: false}, fmt.Errorf("no map found in mapSearch (team state spawner):  %w", err)
			}

			// Get Map gridUtils component
			gridUtils, err := cardinal.GetComponent[comp.GridUtils](world, foundMap)
			if err != nil {
				return msg.CreateMatchResult{Success: false}, fmt.Errorf("error retrieving GridUtils component (team state spawner): %w", err)
			}

			//create a new match
			_, err = cardinal.Create(world, comp.MatchId{MatchId: create.Msg.MatchID}, comp.UID{UID: 0}, comp.Player1{Nickname: create.Tx.PersonaTag}, comp.SpatialHash{Cells: make(map[string]comp.SpatialCell), CellSize: 150, StartX: float32(gridUtils.StartX), StartY: float32(gridUtils.StartY)})
			if err != nil {
				return msg.CreateMatchResult{Success: false}, fmt.Errorf("error creating match: %w", err)
			}

			return msg.CreateMatchResult{Success: true}, nil
		})
}
