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
