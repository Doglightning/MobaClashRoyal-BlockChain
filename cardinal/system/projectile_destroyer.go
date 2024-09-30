package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// Destroys projectiles whos destroyed compoenent is marked true
func ProjectileDestroyerSystem(world cardinal.WorldContext) error {

	// Filter for destoryed projectile
	destroyedFilter := cardinal.ComponentFilter(func(m comp.Destroyed) bool {
		return m.Destroyed
	})

	//go over each destroyed projectile id
	err := cardinal.NewSearch().Entity(
		filter.Contains()).
		Where(destroyedFilter).Each(world, func(id types.EntityID) bool {

		//get matchid and uid of projectile
		MatchID, uid, err := GetComponents2[comp.MatchId, comp.UID](world, id)
		if err != nil {
			fmt.Printf("get projectile components (projectile_destroyer.go): %v", err)
			return false
		}

		//get game state
		gameState, err := getGameStateGSS(world, MatchID)
		if err != nil {
			fmt.Printf("(projectile_destroyer.go) - %v", err)
			return false
		}

		//add projectile id to player1 removal list
		cardinal.UpdateComponent(world, gameState, func(player1 *comp.Player1) *comp.Player1 {
			if player1 == nil {
				fmt.Printf("error retrieving player1 component (projectile_destroyer.go)")
				return nil
			}
			//player1.RemovalList = append(player1.RemovalList, uid.UID)
			player1.RemovalList[uid.UID] = true
			return player1
		})

		//add projectile id to player2 removal list
		cardinal.UpdateComponent(world, gameState, func(player2 *comp.Player2) *comp.Player2 {
			if player2 == nil {
				fmt.Printf("error retrieving player1 component (projectile_destroyer.go)")
				return nil
			}
			//player1.RemovalList = append(player1.RemovalList, uid.UID)
			player2.RemovalList[uid.UID] = true
			return player2
		})

		//remove projectile
		if err := cardinal.Remove(world, id); err != nil {
			fmt.Println("Error removing entity (projectile_destroyer.go):", err)
			return false
		}

		return true
	})

	if err != nil {
		return fmt.Errorf("error retrieving unit entities (projectile_destroyer.go): %w", err)
	}

	return nil
}
