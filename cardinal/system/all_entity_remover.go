package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "MobaClashRoyal/component"
	"MobaClashRoyal/msg"
)

// RemoveAllEntitiesSystem removes all entities associated with a given MatchId when recieve remove_all_entities.go msg
func RemoveAllEntitiesMsgSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage(world,
		func(create cardinal.TxData[msg.RemoveAllEntitiesMsg]) (msg.RemoveAllEntitiesResult, error) {
			// Create a filter to match entities with the specified MatchId.
			matchFilter := cardinal.ComponentFilter(func(m comp.MatchId) bool {
				return m.MatchId == create.Msg.MatchID
			})
			entitySearch := cardinal.NewSearch().Entity(
				filter.Contains(filter.Component[comp.MatchId]())).
				Where(matchFilter)

			// for each ID remove entity based on MatchID search
			err := entitySearch.Each(world, func(id types.EntityID) bool {
				//remove entity
				if err := cardinal.Remove(world, id); err != nil {
					fmt.Println("Error removing entity (all entity remover/RemoveAllEntitiesMsgSystem):", err)
					return false
				}
				return true
			})

			if err != nil {
				return msg.RemoveAllEntitiesResult{Success: false}, fmt.Errorf("error during entity removal (all entity remover/RemoveAllEntitiesMsgSystem): %w", err)
			}

			return msg.RemoveAllEntitiesResult{Success: true}, nil
		})
}

// RemoveAllEntitiesSystem removes all entities associated with a given MatchId.
func RemoveAllEntitiesSystem(world cardinal.WorldContext, matchID string) error {

	// Create a filter to match entities with the specified MatchId.
	matchFilter := cardinal.ComponentFilter(func(m comp.MatchId) bool {
		return m.MatchId == matchID
	})
	entitySearch := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.MatchId]())).
		Where(matchFilter)

	// for each ID remove entity based on MatchID search
	err := entitySearch.Each(world, func(id types.EntityID) bool {
		//remove entity
		if err := cardinal.Remove(world, id); err != nil {
			fmt.Println("Error removing entity (all entity remover/RemoveAllEntitiesSystem):", err)
			return false
		}
		return true
	})

	if err != nil {
		return fmt.Errorf("error during entity removal (all entity remover/RemoveAllEntitiesSystem): %w", err)
	}

	return nil

}
