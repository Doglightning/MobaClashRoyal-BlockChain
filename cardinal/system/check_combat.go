package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// make sure any units who are in range to attack get set for attack b4 attack phase
func CombatCheckSystem(world cardinal.WorldContext) error {
	// filter not in combat units
	combatFilter := cardinal.ComponentFilter(func(m comp.Attack) bool {
		return !m.Combat
	})
	//for each unit with no hp's ids
	err := cardinal.NewSearch().Entity(
		filter.Exact(UnitFilters())).
		Where(combatFilter).Each(world, func(id types.EntityID) bool {
		//get attack component
		uAtk, err := cardinal.GetComponent[comp.Attack](world, id)
		if err != nil {
			fmt.Printf("failed to get attack comp (check_combat.go): %v", err)
			return false
		}

		if !uAtk.Combat { //not in combat
			//get Unit Components
			uPos, uRadius, uAtk, uTeam, MatchID, err := getUnitComponentsCC(world, id)
			if err != nil {
				fmt.Printf("(check_combat.go) -%v", err)
				return false
			}

			// get game state
			gameState, err := getGameStateGSS(world, MatchID)
			if err != nil {
				fmt.Printf("(check_combat.go): - %v", err)
				return false
			}

			// get collision Hash
			collisionHash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
			if err != nil {
				fmt.Printf("error retrieving SpartialHash component on tempSpartialHash (check_combat.go): %s", err)
				return false
			}
			//find closest enemy
			eID, eX, eY, eRadius, found := findClosestEnemy(collisionHash, id, uPos.PositionVectorX, uPos.PositionVectorY, uAtk.AggroRadius, uTeam.Team)
			if found { //found enemy
				// Calculate squared distance between the unit and the enemy, minus their radii
				adjustedDistance := distanceBetweenTwoPoints(uPos.PositionVectorX, uPos.PositionVectorY, eX, eY) - float32(eRadius) - float32(uRadius.UnitRadius)
				//if within attack range
				if adjustedDistance <= float32(uAtk.AttackRadius) {
					uAtk.Combat = true
					uAtk.Target = eID
					//set attack component
					if err = cardinal.SetComponent(world, id, uAtk); err != nil {
						fmt.Printf("error setting attack component (check_combat.go): %v", err)
						return false
					}
				}

			}
		}
		return true
	})

	return err
}

// GetUnitComponents fetches all necessary components related to a unit entity.
func getUnitComponentsCC(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.UnitRadius, *comp.Attack, *comp.Team, *comp.MatchId, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Position component (getUnitComponentsCC): %v", err)
	}
	unitRadius, err := cardinal.GetComponent[comp.UnitRadius](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit Radius component (getUnitComponentsCC): %v", err)
	}
	unitAttack, err := cardinal.GetComponent[comp.Attack](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Unit attack component (getUnitComponentsCC): %v", err)
	}
	team, err := cardinal.GetComponent[comp.Team](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving Team component (getUnitComponentsCC): %v", err)
	}
	matchId, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving MatchId component (getUnitComponentsCC): %v", err)
	}
	return position, unitRadius, unitAttack, team, matchId, nil
}
