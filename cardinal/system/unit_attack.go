package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// system to deal with units attacking each other
func UnitAttackSystem(world cardinal.WorldContext) error {
	// Filter for in combat
	combatFilter := cardinal.ComponentFilter(func(m comp.Attack) bool {
		return m.Combat
	})
	//for every unit in combats id
	err := cardinal.NewSearch().Entity(
		filter.Exact(UnitFilters())).
		Where(combatFilter).Each(world, func(id types.EntityID) bool {
		// get unit attack component
		unitAtk, err := cardinal.GetComponent[comp.Attack](world, id)
		if err != nil {
			fmt.Printf("error retrieving unit Attack component (Unit_Attack.go): %v", err)
			return false
		}

		//if unit is in its damage frame
		if unitAtk.Frame == unitAtk.DamageFrame {
			//get special power component
			unitSp, err := cardinal.GetComponent[comp.Sp](world, id)
			if err != nil {
				fmt.Printf("error retrieving special power component (Unit_Attack.go): %v", err)
				return false
			}
			//if unit is ready to use Special power attack
			if unitSp.CurrentSp == unitSp.MaxSp {
				//get units name
				unitName, err := cardinal.GetComponent[comp.UnitName](world, id)
				if err != nil {
					fmt.Printf("error retrieving unit name component (Unit_Attack.go): %v", err)
					return false
				}
				//spawn special power
				err = spSpawner(world, id, unitName.UnitName)
				if err != nil {
					fmt.Printf("error spawning special attack (Unit_Attack.go): %v - ", err)
					return false
				}

			} else { // normal attack

				if unitAtk.Class == "melee" { //if melee
					//reduce health by units attack damage
					cardinal.UpdateComponent(world, unitAtk.Target, func(health *comp.Health) *comp.Health {
						if health == nil {
							fmt.Printf("error retrieving Health component (Unit_Attack.go)")
							return nil
						}
						health.CurrentHP -= float32(unitAtk.Damage)
						if health.CurrentHP < 0 {
							health.CurrentHP = 0 //never have negative health
						}
						return health
					})
				}

				if unitAtk.Class == "range" { //if range
					//get units component
					unitPosition, matchID, mapName, unitName, err := GetUnitComponentsUA(world, id)
					if err != nil {
						fmt.Printf("%v", err)
						return false
					}
					//get next uid
					UID, err := getNextUID(world, matchID.MatchId)
					if err != nil {
						fmt.Printf("(Unit_Attack.go): %v", err)
						return false
					}
					//create projectile entity
					cardinal.Create(world,
						comp.MatchId{MatchId: matchID.MatchId},
						comp.UID{UID: UID},
						comp.UnitName{UnitName: ProjectileRegistry[unitName.UnitName].Name},
						comp.Movespeed{CurrentMS: ProjectileRegistry[unitName.UnitName].Speed},
						comp.Position{PositionVectorX: unitPosition.PositionVectorX, PositionVectorY: unitPosition.PositionVectorY, PositionVectorZ: unitPosition.PositionVectorZ, RotationVectorX: unitPosition.RotationVectorX, RotationVectorY: unitPosition.RotationVectorY, RotationVectorZ: unitPosition.RotationVectorZ},
						comp.MapName{MapName: mapName.MapName},
						comp.Attack{Target: unitAtk.Target, Class: "projectile", Damage: ProjectileRegistry[unitName.UnitName].Damage},
						comp.Destroyed{Destroyed: false},
					)
				}
			}
			//if our CurrentSp is at the MaxSp threshhold
			if unitSp.CurrentSp >= unitSp.MaxSp {
				unitSp.CurrentSp = 0
			} else {
				unitSp.CurrentSp += unitSp.SpRate //increase sp after attack
				// make sure we are not over MaxSp
				if unitSp.CurrentSp >= unitSp.MaxSp {
					unitSp.CurrentSp = unitSp.MaxSp
				}
			}
			// set updated attack component
			if err := cardinal.SetComponent(world, id, unitSp); err != nil {
				fmt.Printf("error updating special power component (Unit_Attack.go): %s", err)
				return false
			}
		}
		//if our attack frame is at the attack rate reset
		if unitAtk.Frame >= unitAtk.Rate {
			unitAtk.Frame = -1
		}
		unitAtk.Frame++
		// set updated attack component
		if err := cardinal.SetComponent(world, id, unitAtk); err != nil {
			fmt.Printf("error updating attack component (Unit_Attack.go): %s", err)
			return false
		}

		return true
	})

	if err != nil {
		return fmt.Errorf("error retrieving unit entities (Unit_Attack.go): %w", err)
	}
	return nil
}

// GetUnitComponents fetches all necessary components related to a unit entity.
func GetUnitComponentsUA(world cardinal.WorldContext, unitID types.EntityID) (*comp.Position, *comp.MatchId, *comp.MapName, *comp.UnitName, error) {
	position, err := cardinal.GetComponent[comp.Position](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving Position component (Unit_Attack.go): %v", err)
	}

	matchId, err := cardinal.GetComponent[comp.MatchId](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving MatchId component (Unit_Attack.go): %v", err)
	}

	mapName, err := cardinal.GetComponent[comp.MapName](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving mapname component (Unit_Attack.go): %v", err)
	}

	unitName, err := cardinal.GetComponent[comp.UnitName](world, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving UnitName component (Unit_Attack.go): %v", err)
	}
	return position, matchId, mapName, unitName, nil
}
