package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

func vampireSpawn(world cardinal.WorldContext, id types.EntityID) error {

	// get unit attack component
	unitAtk, err := cardinal.GetComponent[comp.Attack](world, id)
	if err != nil {
		return fmt.Errorf("error retrieving unit Attack component (sp_vampire.go): %w", err)
	}

	//reduce health by units attack damage
	err = cardinal.UpdateComponent(world, unitAtk.Target, func(health *comp.Health) *comp.Health {
		if health == nil {
			fmt.Printf("error retrieving Health component (sp_vampire.go)")
			return nil
		}
		health.CurrentHP -= float32(unitAtk.Damage)
		if health.CurrentHP < 0 {
			health.CurrentHP = 0 //never have negative health
		}
		return health
	})

	if err != nil {
		return err
	}

	err = cardinal.UpdateComponent(world, id, func(health *comp.Health) *comp.Health {
		if health == nil {
			fmt.Printf("error geting health component (sp_vampire.go): ")
			return nil
		}
		health.CurrentHP += 30
		if health.CurrentHP > health.MaxHP {
			health.CurrentHP = health.MaxHP
		}
		return health
	})

	return err
}
