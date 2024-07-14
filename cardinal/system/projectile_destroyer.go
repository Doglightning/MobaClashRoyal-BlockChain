package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func ProjectileDestroyerSystem(world cardinal.WorldContext) error {

	// Filter for destoryed projectile
	destroyedFilter := cardinal.ComponentFilter[comp.Destroyed](func(m comp.Destroyed) bool {
		return m.Destroyed
	})

	err := cardinal.NewSearch().Entity(
		filter.Exact(ProjectileFilters())).
		Where(destroyedFilter).Each(world, func(id types.EntityID) bool {

		MatchID, uid, err := getProjectileComponentsPD(world, id)
		if err != nil {
			fmt.Printf("%v", err)
			return false
		}

		//get team state
		gameState, err := getGameStateUM(world, MatchID)
		if err != nil {
			fmt.Printf("%v", err)
			return false
		}

		//get player1 team state
		player1, err := cardinal.GetComponent[comp.Player1](world, gameState)
		if err != nil {
			fmt.Printf("error retrieving player1 component (projectile destroyer): %s", err)
			return false
		}

		//get player2 team state
		player2, err := cardinal.GetComponent[comp.Player2](world, gameState)
		if err != nil {
			fmt.Printf("error retrieving player2 component (projectile destroyer): %s", err)
			return false
		}

		//remove entity
		if err := cardinal.Remove(world, id); err != nil {
			fmt.Println("Error removing entity:", err) // Log error if any
			return false                               // Stop iteration on error
		}

		//player1.RemovalList = append(player1.RemovalList, uid.UID)
		player1.RemovalList[uid.UID] = true
		//player2.RemovalList = append(player2.RemovalList, uid.UID)
		player2.RemovalList[uid.UID] = true

		//add removed unit to player1 removal list component
		if err := cardinal.SetComponent[comp.Player1](world, gameState, player1); err != nil {
			fmt.Printf("error updating player1 component (projectile destroyer): %s", err)
			return false
		}

		//add removed unit to player2 removal list component
		if err := cardinal.SetComponent[comp.Player2](world, gameState, player2); err != nil {
			fmt.Printf("error updating player2 component (projectile destroyer): %s", err)
			return false
		}

		return true
	})

	if err != nil {
		return fmt.Errorf("error retrieving unit entities (projectile destroyer): %w", err)
	}

	return nil
}

// fetches projectile components needed
func getProjectileComponentsPD(world cardinal.WorldContext, id types.EntityID) (matchID *comp.MatchId, uid *comp.UID, err error) {

	matchID, err = cardinal.GetComponent[comp.MatchId](world, id)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving MatchID component (projectile destroyer): %v", err)
	}
	uid, err = cardinal.GetComponent[comp.UID](world, id)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving UID component (projectile destroyer): %v", err)
	}
	return matchID, uid, nil
}
