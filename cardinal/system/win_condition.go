package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// conditions on game to win
func WinCondition(world cardinal.WorldContext) error {
	// Filter for no health
	healthFilter := cardinal.ComponentFilter(func(m comp.UnitName) bool {
		return m.UnitName == "Base"
	})
	//check all structures with no health
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.StructureTag]())).
		Where(healthFilter).Each(world, func(id types.EntityID) bool {

		//get structure Health
		health, err := cardinal.GetComponent[comp.Health](world, id)
		if err != nil {
			fmt.Printf("error getting health component (win condition): %s\n", err)
			return false
		}

		//check if base has died
		if health.CurrentHP <= 0 {
			//get structure matchID
			matchID, err := cardinal.GetComponent[comp.MatchId](world, id)
			if err != nil {
				fmt.Printf("error getting matchID component (win condition): %s\n", err)
				return false
			}
			//remove all entites
			RemoveAllEntitiesSystem(world, matchID.MatchId)
		}

		return true
	})

	return err
}
