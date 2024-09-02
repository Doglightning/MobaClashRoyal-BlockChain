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
			fmt.Printf("error retrieving unit Attack component (phase_Attack.go): %v \n", err)
			return false
		}

		// projectile attack logic
		if atk.Class == "projectile" {
			err = ProjectileAttack(world, id, atk)
			if err != nil {
				fmt.Printf("%v \n", err)
				return false
			}

			// basic melee/range attack logic
		} else if atk.Class == "melee" || atk.Class == "range" {
			err = MeleeRangeAttack(world, id, atk)
			if err != nil {
				fmt.Printf("%v \n", err)
				return false
			}
		} else if atk.Class == "structure" {
			err = StructureAttack(world, id, atk)
			if err != nil {
				fmt.Printf("%v \n", err)
				return false
			}
		}

		return true
	})

	if err != nil {
		return fmt.Errorf("error retrieving unit entities (phase_Attack.go): %w ", err)
	}
	return nil
}

// handles basic range / melee units in combat
func MeleeRangeAttack(world cardinal.WorldContext, id types.EntityID, atk *comp.Attack) error {
	//if unit is in its damage frame
	if atk.Frame == atk.DamageFrame {
		//get special power component
		unitSp, err := cardinal.GetComponent[comp.Sp](world, id)
		if err != nil {
			return fmt.Errorf("error retrieving special power component (phase_Attack.go): %v", err)
		}

		//get units name
		unitName, err := cardinal.GetComponent[comp.UnitName](world, id)
		if err != nil {
			return fmt.Errorf("error retrieving unit name component (phase_Attack.go): %v", err)
		}

		//if unit is ready to use Special power attack
		if unitSp.CurrentSp == unitSp.MaxSp {
			//spawn special power
			err = spSpawner(world, id, unitName.UnitName, unitSp)
			if err != nil {
				return fmt.Errorf("error spawning special attack (phase_Attack.go): %v - ", err)
			}

		} else { // normal attack

			err := ClassAttack(world, id, unitName.UnitName, atk)
			if err != nil {
				return err
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
			return fmt.Errorf("error updating special power component (phase_Attack.go): %s ", err)
		}
	}
	//if our attack frame is at the attack rate reset
	if atk.Frame >= atk.Rate {
		atk.Frame = -1
	}
	atk.Frame++
	// set updated attack component
	if err := cardinal.SetComponent(world, id, atk); err != nil {
		return fmt.Errorf("error updating attack component (phase_Attack.go): %s ", err)
	}

	return nil
}

// handles structure units in combat
func StructureAttack(world cardinal.WorldContext, id types.EntityID, atk *comp.Attack) error {
	//if unit is in its damage frame
	if atk.Frame == atk.DamageFrame {

		//get units name
		unitName, err := cardinal.GetComponent[comp.UnitName](world, id)
		if err != nil {
			return fmt.Errorf("error retrieving unit name component (phase_Attack.go): %v", err)
		}

		err = ClassAttack(world, id, unitName.UnitName, atk)
		if err != nil {
			return err
		}

	}
	//if our attack frame is at the attack rate reset
	if atk.Frame >= atk.Rate {
		atk.Frame = -1
	}
	atk.Frame++
	// set updated attack component
	if err := cardinal.SetComponent(world, id, atk); err != nil {
		return fmt.Errorf("error updating attack component (phase_Attack.go): %s ", err)
	}

	return nil
}

// handles projectiles in combat (they are in range to deal dmg to enemy)
func ProjectileAttack(world cardinal.WorldContext, id types.EntityID, projectileAttack *comp.Attack) error {
	//get targets health compoenent from the projectiles attack target
	enemyHealth, err := cardinal.GetComponent[comp.Health](world, projectileAttack.Target)
	if err != nil {
		return fmt.Errorf("error getting enemy Health component (projectile_Attack - phase_Attack.go): %v ", err)
	}

	//reduce enemy HP
	enemyHealth.CurrentHP -= float32(projectileAttack.Damage)
	if enemyHealth.CurrentHP < 0 {
		enemyHealth.CurrentHP = 0
	}
	//set enemy HP compoenent
	err = cardinal.SetComponent(world, projectileAttack.Target, enemyHealth)
	if err != nil {
		return fmt.Errorf("error setting Health component (projectile_Attack - phase_Attack.go): %v ", err)
	}
	//set projectime combat to false
	projectileAttack.Combat = false
	//set attack component
	if err := cardinal.SetComponent(world, id, projectileAttack); err != nil {
		return fmt.Errorf("error updating attack component (projectile_Attack - phase_Attack.go): %v ", err)
	}

	//update projectiles destroyed component to True
	cardinal.UpdateComponent(world, id, func(destroyed *comp.Destroyed) *comp.Destroyed {
		if destroyed == nil {
			fmt.Printf("error retrieving enemy destroyed component (projectile_Attack - phase_Attack.go): \n")
			return nil
		}
		destroyed.Destroyed = true
		return destroyed
	})
	return nil
}