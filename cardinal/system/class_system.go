package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// updates all SP's spawned
func SpUpdater(world cardinal.WorldContext) error {

	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.SpEntity]())).Each(world, func(id types.EntityID) bool {
		//get sp name
		spEntity, err := cardinal.GetComponent[comp.SpEntity](world, id)
		if err != nil {
			fmt.Printf("error getting sp name component (SpUpdater): %v", err)
			return false
		}

		if spEntity.SpName == "ArcherLadySP" {
			err = archerLadyUpdate(world, id)
			if err != nil {
				fmt.Printf("%v", err)
				return false
			}

		}

		if spEntity.SpName == "MageSP" {
			err = MageUpdate(world, id)
			if err != nil {
				fmt.Printf("%v", err)
				return false
			}
		}

		if spEntity.SpName == "VampireSP" {
			err = vampireUpdateSP(world, id)
			if err != nil {
				fmt.Printf("%v", err)
				return false
			}
		}

		return true
	})

	return err
}

// spawns the special attack
func spSpawner(world cardinal.WorldContext, id types.EntityID, name string) error {
	var err error
	if name == "ArcherLady" {
		err = archerLadySpawn(world, id)
	}
	if name == "FireSpirit" {
		err = fireSpiritSpawn(world, id)
	}

	if name == "Mage" {
		err = MageSpawnSP(world, id)
	}

	if name == "Vampire" {
		err = vampireSpawnSP(world, id)
	}
	return err
}

// triggers unit attack
func ClassAttack(world cardinal.WorldContext, id types.EntityID, name string, atk *comp.Attack) error {
	var err error

	if name == "ArcherLady" {
		err = archerLadyAttack(world, id, atk)
	}

	if name == "Mage" {
		err = mageAttack(world, id, atk)
	}

	if name == "Tower" || name == "Base" {
		err = towerAttack(world, id, atk)
	}

	if name == "Vampire" {
		err = vampireAttack(world, atk)
	}
	return err
}

// sets attack system for how units engage in combat
func ClassAttackSystem(world cardinal.WorldContext, id types.EntityID, atk *comp.Attack) error {

	name, err := cardinal.GetComponent[comp.UnitName](world, id)
	if err != nil {
		return fmt.Errorf("error getting name component (class attack system): %v", err)
	}

	if name.UnitName == "FireSpirit" {
		err = FireSpiritAttack(world, id, atk)

	} else {
		err = MeleeRangeAttack(world, id, atk)
	}

	return err
}

// on desetry resets combat for units targeting
func ClassResetCombat(world cardinal.WorldContext, id types.EntityID, name string) error {

	var err error

	if name == "FireSpirit" {
		err = fireSpiritResetCombat(world, id)
	} else {
		//reset attack component
		err := cardinal.UpdateComponent(world, id, func(attack *comp.Attack) *comp.Attack {
			if attack == nil {
				fmt.Printf("error retrieving enemy attack component (Phase attack.go): ")
				return nil
			}
			attack.Combat = false
			attack.Frame = 0
			return attack
		})
		if err != nil {
			return fmt.Errorf("error updating attack comp (Phase attack.go): %v", err)
		}

	}

	return err
}
