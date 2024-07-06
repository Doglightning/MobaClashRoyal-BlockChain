package system

import (
	comp "MobaClashRoyal/component"
	"MobaClashRoyal/msg"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
)

func RemovalListSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage[msg.RemoveUnitMsg, msg.RemoveUnitResult](
		world,
		func(create cardinal.TxData[msg.RemoveUnitMsg]) (msg.RemoveUnitResult, error) {

			teamFilter := cardinal.ComponentFilter[comp.MatchId](func(m comp.MatchId) bool {
				return m.MatchId == create.Msg.MatchId
			})

			teamSearch, err := cardinal.NewSearch().Entity(
				filter.Exact(TeamFilters())).
				Where(teamFilter).First(world)

			if err != nil {
				return msg.RemoveUnitResult{Succsess: false}, fmt.Errorf("error searching for team (Removal State Query): %w", err)
			}

			if create.Msg.Team == "Blue" {
				// Get Player1 component
				player1, err := cardinal.GetComponent[comp.Player1](world, teamSearch)
				if err != nil {
					return msg.RemoveUnitResult{Succsess: false}, fmt.Errorf("error retrieving Player1 component (Removal State Query): %w", err)
				}

				for _, key := range create.Msg.RemovalList {
					delete(player1.RemovalList, key)
				}

				//add removed unit to player2 removal list component
				if err := cardinal.SetComponent(world, teamSearch, player1); err != nil {

					return msg.RemoveUnitResult{Succsess: false}, fmt.Errorf("error updating player1 component (Removal State Query): %w", err)
				}
			} else {
				// Get Player1 component
				player2, err := cardinal.GetComponent[comp.Player2](world, teamSearch)
				if err != nil {
					return msg.RemoveUnitResult{Succsess: false}, fmt.Errorf("error retrieving Player2 component (Removal State Query): %w", err)
				}

				for _, key := range create.Msg.RemovalList {
					delete(player2.RemovalList, key)
				}
				//add removed unit to player2 removal list component
				if err := cardinal.SetComponent(world, teamSearch, player2); err != nil {
					return msg.RemoveUnitResult{Succsess: false}, fmt.Errorf("error updating player2 component (Removal State Query): %w", err)
				}
			}

			return msg.RemoveUnitResult{Succsess: true}, nil
		})
}
