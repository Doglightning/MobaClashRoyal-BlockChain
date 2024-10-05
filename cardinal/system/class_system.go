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

		switch spEntity.SpName {
		case "ArcherLadySP":
			err = archerLadyUpdate(world, id)
		case "MageSP":
			err = MageUpdate(world, id)
		case "VampireSP":
			err = vampireUpdateSP(world, id)
		}

		if err != nil {
			fmt.Printf("%v \n", err)
			return false
		}

		return true
	})

	return err
}

// spawns the special attack
func spSpawner(world cardinal.WorldContext, id types.EntityID, name string) error {
	var err error

	switch name {
	case "ArcherLady":
		err = archerLadySpawn(world, id)
	case "FireSpirit":
		err = fireSpiritSpawn(world, id)
	case "Mage":
		err = MageSpawnSP(world, id)
	case "Vampire":
		err = vampireSpawnSP(world, id)
	}

	return err
}

// triggers unit attack
func ClassAttack(world cardinal.WorldContext, id types.EntityID, name string, atk *comp.Attack) error {
	var err error

	switch name {
	case "ArcherLady":
		err = archerLadyAttack(world, id, atk)
	case "Base":
		err = towerAttack(world, id, atk)
	case "Mage":
		err = mageAttack(world, id, atk)
	case "Tower":
		err = towerAttack(world, id, atk)
	case "Vampire":
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

	switch name.UnitName {
	case "FireSpirit":
		err = FireSpiritAttack(world, id, atk)
	case "LeafBird":
		err = leafBirdAttackSystem(world, id, atk)
	default:
		err = MeleeRangeAttack(world, id, atk)
	}

	return err
}

// on desetry resets combat for units targeting
func ClassResetCombat(world cardinal.WorldContext, id types.EntityID, atk *comp.Attack) error {

	name, err := cardinal.GetComponent[comp.UnitName](world, id)
	if err != nil {
		return fmt.Errorf("error getting name component (class attack system): %v", err)
	}

	switch name.UnitName {
	case "FireSpirit":
		err = channelingResetCombat(world, id, atk)
	case "LeafBird":
		err = channelingResetCombat(world, id, atk)
	default:
		resetCombat(atk)
	}

	return err
}

// logic of how a unit destroys itself
func ClassDestroySystem(world cardinal.WorldContext, id types.EntityID) error {
	// name, err := cardinal.GetComponent[comp.UnitName](world, id)
	// if err != nil {
	// 	return fmt.Errorf("error getting name component (class attack system): %v", err)
	// }

	//default
	err := unitDestroyerDefault(world, id)

	return err
}
