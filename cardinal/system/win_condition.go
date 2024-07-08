package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func WinCondition(world cardinal.WorldContext) error {

	// Filter for unit with no HP
	unitFilter := cardinal.ComponentFilter[comp.UnitHealth](func(m comp.UnitHealth) bool {
		return m.CurrentHP == 0
	})

	err := cardinal.NewSearch().Entity(
		filter.Exact(StructureFilters())).
		Where(unitFilter).Each(world, func(id types.EntityID) bool {

		//get structure health
		matchID, err := cardinal.GetComponent[comp.MatchId](world, id)
		if err != nil {
			fmt.Printf("error getting matchID component (win condition): %s\n", err)
			return false
		}
		RemoveAllEntitiesSystem(world, matchID.MatchId)

		_, err = cardinal.Create(world,
			comp.MatchId{MatchId: matchID.MatchId},
		)
		if err != nil {
			fmt.Printf("error getting matchID component (win condition): %s\n", err)
			return false
		}
		test
		return true
	})

	return err
}
