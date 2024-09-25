package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
)

func leafBirdAttack(world cardinal.WorldContext, atk *comp.Attack) error {
	// reduce health by units attack damage
	err := cardinal.UpdateComponent(world, atk.Target, func(health *comp.Health) *comp.Health {
		if health == nil {
			fmt.Printf("error retrieving Health component (Unit_Attack.go) \n")
			return nil
		}
		health.CurrentHP -= float32(atk.Damage)
		if health.CurrentHP < 0 {
			health.CurrentHP = 0 //never have negative health
		}
		return health
	})
	if err != nil {
		return fmt.Errorf("error on vampire attack (class vampire.go): %v", err)
	}

	return nil
}
