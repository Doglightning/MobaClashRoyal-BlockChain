package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

var towerHealing float32 = 2

// Heal towers while converting
func TowerConverterSystem(world cardinal.WorldContext) error {
	// Filter for no HP
	stateFilter := cardinal.ComponentFilter(func(m comp.State) bool {
		return m.State == "Converting"
	})
	//for each tower still converting teams
	err := cardinal.NewSearch().Entity(
		filter.Contains(StructureFilters())).
		Where(stateFilter).Each(world, func(id types.EntityID) bool {

		// increase tower hp until full
		err := cardinal.UpdateComponent(world, id, func(health *comp.Health) *comp.Health {
			if health == nil {
				fmt.Printf("error retrieving health component (tower conversion.go)")
				return nil
			}
			health.CurrentHP += towerHealing
			if health.CurrentHP >= health.MaxHP {
				health.CurrentHP = health.MaxHP

				// set tower state to Default
				err := cardinal.UpdateComponent(world, id, func(state *comp.State) *comp.State {
					if state == nil {
						fmt.Printf("error retrieving state component (tower conversion.go)")
						return nil
					}
					state.State = "Default"
					return state
				})

				if err != nil {
					fmt.Printf("error updating state component (tower conversion.go): %s", err)
					return health
				}

			}
			return health
		})

		if err != nil {
			fmt.Printf("error updating health component (tower conversion.go): %s", err)
			return false
		}

		return true
	})
	return err
}
