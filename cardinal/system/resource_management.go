package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

var goldGen float32 = .1

// function to regenerate gold at the rate of goldGen per tick
func GoldGeneration(world cardinal.WorldContext) error {

	err := cardinal.NewSearch().Entity(
		filter.Contains(GameStateFilters())).
		Each(world, func(id types.EntityID) bool {

			//increment player1 gold
			err := cardinal.UpdateComponent(world, id, func(player1 *comp.Player1) *comp.Player1 {
				if player1 == nil {
					fmt.Printf("error getting player1 gold (resource_management.go):\n")
					return nil
				}
				player1.Gold += goldGen
				//cap gold to 10
				if player1.Gold > 10 {
					player1.Gold = 10
				}

				return player1
			})

			if err != nil {
				return false
			}

			//increment player2 gold
			err = cardinal.UpdateComponent(world, id, func(player2 *comp.Player2) *comp.Player2 {
				if player2 == nil {
					fmt.Printf("error getting player2 gold (resource_management.go):\n")
					return nil
				}
				player2.Gold += goldGen
				//cap gold to 10
				if player2.Gold > 10 {
					player2.Gold = 10
				}
				return player2
			})

			return err == nil
		})

	return err
}
