package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// system to deal with objects attacking each other
func AttackPhaseSystem(world cardinal.WorldContext) error {
	// Filter for in combat
	combatFilter := cardinal.ComponentFilter(func(m comp.Attack) bool {
		return m.Combat
	})
	//for every object in combats id
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.Attack]())).
		Where(combatFilter).Each(world, func(id types.EntityID) bool {
		// get attack component
		atk, err := cardinal.GetComponent[comp.Attack](world, id)
		if err != nil {
			fmt.Printf("error retrieving unit Attack component (Unit_Attack.go): %v", err)
			return false
		}

		if atk.Class == "projectile" {
			err = ProjectileAttack(world, id, atk)
			if err != nil {
				fmt.Printf("%v \n", err)
				return false
			}
		}

		//if unit is in its damage frame
		if atk.Frame == atk.DamageFrame {
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
				err = spSpawner(world, id, unitName.UnitName, unitSp)
				if err != nil {
					fmt.Printf("error spawning special attack (Unit_Attack.go): %v - ", err)
					return false
				}

			} else { // normal attack

				if atk.Class == "melee" { //if melee
					//reduce health by units attack damage
					cardinal.UpdateComponent(world, atk.Target, func(health *comp.Health) *comp.Health {
						if health == nil {
							fmt.Printf("error retrieving Health component (Unit_Attack.go)")
							return nil
						}
						health.CurrentHP -= float32(atk.Damage)
						if health.CurrentHP < 0 {
							health.CurrentHP = 0 //never have negative health
						}
						return health
					})
				}

				if atk.Class == "range" { //if range
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
						comp.Attack{Target: atk.Target, Class: "projectile", Damage: UnitRegistry[unitName.UnitName].Damage},
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
		if atk.Frame >= atk.Rate {
			atk.Frame = -1
		}
		atk.Frame++
		// set updated attack component
		if err := cardinal.SetComponent(world, id, atk); err != nil {
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

// handles projectiles in combat (they are in range to deal dmg to enemy)
func ProjectileAttack(world cardinal.WorldContext, id types.EntityID, projectileAttack *comp.Attack) error {
	//get targets health compoenent from the projectiles attack target
	enemyHealth, err := cardinal.GetComponent[comp.Health](world, projectileAttack.Target)
	if err != nil {
		return fmt.Errorf("error getting enemy Health component (projectile_Attack.go): %v", err)
	}

	//reduce enemy HP
	enemyHealth.CurrentHP -= float32(projectileAttack.Damage)
	if enemyHealth.CurrentHP < 0 {
		enemyHealth.CurrentHP = 0
	}
	//set enemy HP compoenent
	err = cardinal.SetComponent(world, projectileAttack.Target, enemyHealth)
	if err != nil {
		return fmt.Errorf("error setting Health component (projectile_Attack.go): %v", err)
	}
	//set projectime combat to false
	projectileAttack.Combat = false
	//set attack component
	if err := cardinal.SetComponent(world, id, projectileAttack); err != nil {
		return fmt.Errorf("error updating attack component (projectile_Attack.go): %v", err)
	}

	//update projectiles destroyed component to True
	cardinal.UpdateComponent(world, id, func(destroyed *comp.Destroyed) *comp.Destroyed {
		if destroyed == nil {
			fmt.Printf("error retrieving enemy destroyed component (projectile_Attack.go): ")
			return nil
		}
		destroyed.Destroyed = true
		return destroyed
	})
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
