package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// spawns the special attack
func spSpawner(world cardinal.WorldContext, id types.EntityID, name string) error {
	var err error
	if name == "ArcherLady" {
		err = archerLadyAttack(world, id)

	}

	if name == "Vampire" {
		fmt.Println("hello")
	}
	return err
}
