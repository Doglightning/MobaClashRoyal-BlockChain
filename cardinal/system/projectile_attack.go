package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

func ProjectileAttackSystem(world cardinal.WorldContext) error {

	// Filter for projectile in combat
	combatFilter := cardinal.ComponentFilter[comp.Attack](func(m comp.Attack) bool {
		return m.Combat
	})

	err := cardinal.NewSearch().Entity(
		filter.Exact(ProjectileFilters())).
		Where(combatFilter).Each(world, func(id types.EntityID) bool {

		projectileAttack, err := cardinal.GetComponent[comp.Attack](world, id)
		if err != nil {
			fmt.Printf("error retrieving projectile Attack component (projectile Attack): %v", err)
			return false
		}

		enemyHealth, err := cardinal.GetComponent[comp.UnitHealth](world, projectileAttack.Target)
		if err != nil {
			fmt.Printf("error getting enemy Health component (projectile Attack): %v", err)
			return false
		}
		enemyHealth.CurrentHP -= float32(projectileAttack.Damage)
		if enemyHealth.CurrentHP < 0 {
			enemyHealth.CurrentHP = 0
		}
		err = cardinal.SetComponent(world, projectileAttack.Target, enemyHealth)
		if err != nil {
			fmt.Printf("error setting Health component (projectile Attack): %v", err)
			return false
		}

		projectileAttack.Combat = false

		if err := cardinal.SetComponent[comp.Attack](world, id, projectileAttack); err != nil {
			fmt.Printf("error updating attack component (projectile attack): %s", err)
			return false
		}

		cardinal.UpdateComponent(world, id, func(destroyed *comp.Destroyed) *comp.Destroyed {
			if destroyed == nil {
				fmt.Printf("error retrieving enemy destroyed component (projectile attack): ")
				return nil
			}
			destroyed.Destroyed = true
			return destroyed
		})
		return true
	})

	if err != nil {
		return fmt.Errorf("error retrieving projectile entities (projectile attack): %w", err)
	}

	return nil
}
