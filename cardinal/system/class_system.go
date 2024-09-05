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
func spSpawner(world cardinal.WorldContext, id types.EntityID, name string, sp *comp.Sp) error {
	var err error
	if name == "ArcherLady" {
		err = archerLadySpawn(world, id)
	}

	if name == "Vampire" {
		err = vampireSpawnSP(world, id, sp)
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
