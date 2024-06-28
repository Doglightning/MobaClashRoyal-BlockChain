package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"

	comp "MobaClashRoyal/component"
	"MobaClashRoyal/msg"
)

// RemoveAllEntitiesSystem removes all entities associated with a given MatchId when a game ends.
func RemoveAllEntitiesSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.RemoveAllEntitiesMsg, msg.RemoveAllEntitiesResult](
		world,
		func(create cardinal.TxData[msg.RemoveAllEntitiesMsg]) (msg.RemoveAllEntitiesResult, error) {
			// Create a filter to match entities with the specified MatchId.
			matchFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
				return m.MatchId == create.Msg.MatchID
			})

			entitySearch := cardinal.NewSearch().Entity(
				filter.Contains(filter.Component[comp.MatchId]())).
				Where(matchFilter)

			// Attempt to remove each entity found that matches the filter.
			err := entitySearch.Each(world, func(id types.EntityID) bool {

				if err := cardinal.Remove(world, id); err != nil {
					fmt.Println("Error removing entity:", err) // Log error if any
					return false                               // Stop iteration on error
				}
				return true // Continue if successful
			})

			if err != nil {
				return msg.RemoveAllEntitiesResult{Success: false}, fmt.Errorf("error during entity removal: %w", err)
			}

			return msg.RemoveAllEntitiesResult{Success: true}, nil
		})
}
