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
		filter.Contains(filter.Component[comp.SpName]())).Each(world, func(id types.EntityID) bool {
		//get sp name
		spName, err := cardinal.GetComponent[comp.SpName](world, id)
		if err != nil {
			fmt.Printf("error getting sp name component (SpUpdater): %v", err)
			return false
		}

		if spName.SpName == "ArcherLadySP" {
			err = archerLadyUpdate(world, id)
			if err != nil {
				fmt.Printf("%v", err)
				return false
			}

		}

		if spName.SpName == "VampireSP" {
			fmt.Println("hello")
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

	if name == "Vampire" {
		fmt.Println("hello")
	}
	return err
}
