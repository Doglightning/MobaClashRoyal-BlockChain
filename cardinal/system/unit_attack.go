package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func UnitAttackSystem(world cardinal.WorldContext) error {

	// Filter for current map
	unitFilter := cardinal.ComponentFilter[comp.Attack](func(m comp.Attack) bool {
		return m.Combat
	})

	err := cardinal.NewSearch().Entity(
		filter.Exact(UnitFilters())).
		Where(unitFilter).Each(world, func(id types.EntityID) bool {

		unitAttack, err := cardinal.GetComponent[comp.Attack](world, id)
		if err != nil {
			fmt.Printf("error retrieving unitAttack component (Unit Attack): %v", err)
			return false
		}

		if unitAttack.Frame == unitAttack.DamageFrame {
			if unitAttack.Class == "melee" {
				enemyHealth, err := cardinal.GetComponent[comp.UnitHealth](world, unitAttack.Target)
				if err != nil {
					fmt.Printf("error getting enemy Health component (Unit Attack): %v", err)
					return false
				}
				enemyHealth.CurrentHP -= float32(unitAttack.Damage)
				if enemyHealth.CurrentHP < 0 {
					enemyHealth.CurrentHP = 0
				}
				err = cardinal.SetComponent(world, unitAttack.Target, enemyHealth)
				if err != nil {
					fmt.Printf("error setting Health component (Unit Attack): %v", err)
					return false
				}
			}

			if unitAttack.Class == "range" {

				unitPosition, matchID, mapName, unitName, err := GetUnitComponentsUA(world, id)
				if err != nil {
					fmt.Printf("%v", err)
					return false
				}

				UID, err := getNextUID(world, matchID.MatchId)
				if err != nil {
					fmt.Printf("(Unit Attack): %v", err)
					return false
				}

				cardinal.Create(world,
					comp.MatchId{MatchId: matchID.MatchId},
					comp.UID{UID: UID},
					comp.UnitName{UnitName: ProjectileRegistry[unitName.UnitName].Name},
					comp.Movespeed{CurrentMS: ProjectileRegistry[unitName.UnitName].Speed},
					comp.Position{PositionVectorX: unitPosition.PositionVectorX, PositionVectorY: unitPosition.PositionVectorY, PositionVectorZ: unitPosition.PositionVectorZ, RotationVectorX: unitPosition.RotationVectorX, RotationVectorY: unitPosition.RotationVectorY, RotationVectorZ: unitPosition.RotationVectorZ},
					comp.MapName{MapName: mapName.MapName},
					comp.Attack{Target: unitAttack.Target, Class: "projectile", Damage: unitAttack.Damage},
					comp.Destroyed{Destroyed: false},
				)
			}
		}

		if unitAttack.Frame >= unitAttack.Rate {
			unitAttack.Frame = -1
		}

		unitAttack.Frame++

		if err := cardinal.SetComponent[comp.Attack](world, id, unitAttack); err != nil {
			fmt.Printf("error updating attack component (unit attack): %s", err)
			return false
		}

		return true
	})

	if err != nil {
		return fmt.Errorf("error retrieving unit entities (unit attack): %w", err)
	}

	return nil
}

// GetUnitComponents fetches all necessary components related to a unit entity.
func GetUnitComponentsUA(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.MatchId, *comp.MapName, *comp.UnitName, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving Position component (unit attack): %v", err)
	}

	matchId, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving MatchId component (unit attack): %v", err)
	}

	mapName, err := cardinal.GetComponent[comp.MapName](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving Distance component (unit attack): %v", err)
	}

	unitName, err := cardinal.GetComponent[comp.UnitName](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving UnitName component (unit attack): %v", err)
	}
	return position, matchId, mapName, unitName, nil
}
