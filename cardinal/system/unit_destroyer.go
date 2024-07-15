package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func UnitDestroyerSystem(world cardinal.WorldContext) error {

	// Filter for unit with no HP
	unitFilter := cardinal.ComponentFilter[comp.Health](func(m comp.Health) bool {
		return m.CurrentHP <= 0
	})

	err := cardinal.NewSearch().Entity(
		filter.Exact(UnitFilters())).
		Where(unitFilter).Each(world, func(id types.EntityID) bool {

		MatchID, uid, UnitPosition, UnitRadius, err := getUnitComponentsUD(world, id)
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
			fmt.Printf("error retrieving player1 component (unit destroyer): %s", err)
			return false
		}

		//get player2 team state
		player2, err := cardinal.GetComponent[comp.Player2](world, gameState)
		if err != nil {
			fmt.Printf("error retrieving player2 component (unit destroyer): %s", err)
			return false
		}

		////////////////////////////Filter for unit with no HP/////////////////////////////////////////
		targetFilter := cardinal.ComponentFilter[comp.Attack](func(m comp.Attack) bool {
			return m.Target == id
		})

		err = cardinal.NewSearch().Entity(
			filter.Exact(UnitFilters())).
			Where(targetFilter).Each(world, func(enemyID types.EntityID) bool {

			cardinal.UpdateComponent(world, enemyID, func(attack *comp.Attack) *comp.Attack {
				if attack == nil {
					fmt.Printf("error retrieving enemy attack component (unit destroyer): ")
					return nil
				}
				attack.Combat = false
				attack.Frame = 0
				return attack
			})

			return true
		})

		if err != nil {
			fmt.Printf("error retrieving unit entities (unit destroyer): %s", err)
			return false
		}
		////////////////////////inner search reminder area////////////////////////////////////////////

		////////////////////////////Projectile Destroyer/////////////////////////////////////////

		err = cardinal.NewSearch().Entity(
			filter.Exact(ProjectileFilters())).
			Where(targetFilter).Each(world, func(projectileID types.EntityID) bool {

			//get player1 team state
			projectileUID, err := cardinal.GetComponent[comp.UID](world, projectileID)
			if err != nil {
				fmt.Printf("error retrieving projectile UID component (unit destroyer): %s", err)
				return false
			}

			player1.RemovalList[projectileUID.UID] = true
			player2.RemovalList[projectileUID.UID] = true
			//remove entity
			if err := cardinal.Remove(world, projectileID); err != nil {
				fmt.Println("Error removing entity projectile (unit destroyer):", err)
				return false
			}
			return true
		})

		if err != nil {
			fmt.Printf("error retrieving projectile entities (unit destroyer): %s", err)
			return false
		}
		////////////////////////inner search reminder area////////////////////////////////////////////

		//get Spatial Hash
		CollisionSpartialHash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
		if err != nil {
			fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (unit destroyer): %s", err)
			return false
		}
		RemoveObjectFromSpatialHash(CollisionSpartialHash, id, UnitPosition.PositionVectorX, UnitPosition.PositionVectorY, UnitRadius.UnitRadius)

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
			fmt.Printf("error updating player1 component (unit destroyer): %s", err)
			return false
		}

		//add removed unit to player2 removal list component
		if err := cardinal.SetComponent[comp.Player2](world, gameState, player2); err != nil {
			fmt.Printf("error updating player2 component (unit destroyer): %s", err)
			return false
		}

		return true
	})

	if err != nil {
		return fmt.Errorf("error retrieving unit entities (unit destroyer): %w", err)
	}

	return nil
}

// fetches unit components needed for spatial hash removal
func getUnitComponentsUD(world cardinal.WorldContext, id types.EntityID) (matchID *comp.MatchId, uid *comp.UID, unitPosition *comp.Position, unitRadius *comp.UnitRadius, err error) {

	unitPosition, err = cardinal.GetComponent[comp.Position](world, id)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving enemy Position component (unit destroyer): %v", err)
	}
	unitRadius, err = cardinal.GetComponent[comp.UnitRadius](world, id)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving enemy Radius component (unit destroyer): %v", err)
	}
	matchID, err = cardinal.GetComponent[comp.MatchId](world, id)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving MatchID component (unit destroyer): %v", err)
	}
	uid, err = cardinal.GetComponent[comp.UID](world, id)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving UID component (unit destroyer): %v", err)
	}
	return matchID, uid, unitPosition, unitRadius, nil
}
