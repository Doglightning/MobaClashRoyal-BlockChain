package system

import (
	comp "MobaClashRoyal/component"
	"MobaClashRoyal/msg"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
)

// Given a list of units to remove from player removal list on game state.
// marks that client for that player has removed from their end.
// called with removal_list_msg.go
func RemovalListSystem(world cardinal.WorldContext) error {
	return cardinal.EachMessage(world,
		func(create cardinal.TxData[msg.RemoveUnitMsg]) (msg.RemoveUnitResult, error) {
			//filter for matchID
			matchFilter := cardinal.ComponentFilter(func(m comp.MatchId) bool {
				return m.MatchId == create.Msg.MatchId
			})
			//find game state
			gameState, err := cardinal.NewSearch().Entity(
				filter.Contains(filter.Component[comp.GameStateTag]())).
				Where(matchFilter).First(world)
			if err != nil {
				return msg.RemoveUnitResult{Succsess: false}, fmt.Errorf("error searching for team (Removal_list_system.go): %w", err)
			}

			//if blue team
			if create.Msg.Team == "Blue" {
				//remove all ids from msg in removal list for player1
				cardinal.UpdateComponent(world, gameState, func(player1 *comp.Player1) *comp.Player1 {
					if player1 == nil {
						fmt.Printf("error retrieving player1 component (Removal_list_system.go)")
						return nil
					}
					// for each id in msg
					for _, key := range create.Msg.RemovalList {
						//remove id key from list
						delete(player1.RemovalList, key)
					}
					return player1
				})
			} else {
				//remove all ids from msg in removal list for player2
				cardinal.UpdateComponent(world, gameState, func(player2 *comp.Player2) *comp.Player2 {
					if player2 == nil {
						fmt.Printf("error retrieving player1 component (Removal_list_system.go)")
						return nil
					}
					// for each id in msg
					for _, key := range create.Msg.RemovalList {
						//remove id key from list
						delete(player2.RemovalList, key)
					}
					return player2
				})
			}

			return msg.RemoveUnitResult{Succsess: true}, nil
		})
}
