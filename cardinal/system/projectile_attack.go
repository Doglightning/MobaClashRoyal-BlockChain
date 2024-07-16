package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// System to make projectiles inflict damage to enemies they are in combat with
func ProjectileAttackSystem(world cardinal.WorldContext) error {

	// Filter for projectile in combat (only set in combat when in range)
	combatFilter := cardinal.ComponentFilter(func(m comp.Attack) bool {
		return m.Combat
	})
	// go through each projectile id
	err := cardinal.NewSearch().Entity(
		filter.Exact(ProjectileFilters())).
		Where(combatFilter).Each(world, func(id types.EntityID) bool {

		//get projectile attack component
		projectileAttack, err := cardinal.GetComponent[comp.Attack](world, id)
		if err != nil {
			fmt.Printf("error retrieving projectile Attack component (projectile_Attack.go): %v", err)
			return false
		}

		//get targets health compoenent from the projectiles attack target
		enemyHealth, err := cardinal.GetComponent[comp.Health](world, projectileAttack.Target)
		if err != nil {
			fmt.Printf("error getting enemy Health component (projectile_Attack.go): %v", err)
			return false
		}

		//reduce enemy HP
		enemyHealth.CurrentHP -= float32(projectileAttack.Damage)
		if enemyHealth.CurrentHP < 0 {
			enemyHealth.CurrentHP = 0
		}
		//set enemy HP compoenent
		err = cardinal.SetComponent(world, projectileAttack.Target, enemyHealth)
		if err != nil {
			fmt.Printf("error setting Health component (projectile_Attack.go): %v", err)
			return false
		}
		//set projectime combat to false
		projectileAttack.Combat = false
		//set attack component
		if err := cardinal.SetComponent(world, id, projectileAttack); err != nil {
			fmt.Printf("error updating attack component (projectile_Attack.go): %s", err)
			return false
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
		return true
	})

	if err != nil {
		return fmt.Errorf("error retrieving projectile entities (projectile_Attack.go): %w", err)
	}

	return nil
}
