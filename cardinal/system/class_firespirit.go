package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

func fireSpiritSpawn(world cardinal.WorldContext, id types.EntityID) error {

	atk, err := cardinal.GetComponent[comp.Attack](world, id)
	if err != nil {
		return fmt.Errorf("error getting attack component (class fireSpirit.go): %v", err)
	}
	// reduce health by units attack damage
	err = cardinal.UpdateComponent(world, atk.Target, func(health *comp.Health) *comp.Health {
		if health == nil {
			fmt.Printf("error retrieving Health component (class fireSpirit.go) \n")
			return nil
		}
		health.CurrentHP -= float32(atk.Damage)
		if health.CurrentHP < 0 {
			health.CurrentHP = 0 //never have negative health
		}
		return health
	})
	if err != nil {
		return fmt.Errorf("error updating health (class fireSpirit.go): %v", err)
	}

	return nil
}
