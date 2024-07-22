package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// vampireSP struct contains configuration for an vampires special properties.
type vampireSP struct {
	healCount  int
	healAmount float32
}

// NewVampireSP creates a new instance of vampireSP with default settings.
func NewVampireSP() *vampireSP {
	return &vampireSP{
		healCount:  25,
		healAmount: 1,
	}
}

func vampireUpdateSP(world cardinal.WorldContext, id types.EntityID) error {
	vampire := NewVampireSP()
	healCount, err := cardinal.GetComponent[comp.IntTracker](world, id)
	if err != nil {
		return fmt.Errorf("error retrieving int tracker component (sp_vampire.go): %w", err)
	}
	fmt.Println("hi")
	if healCount.Num < vampire.healCount {
		healCount.Num += 1
		err = cardinal.SetComponent(world, id, healCount)
		if err != nil {
			return fmt.Errorf("error setting int tracker component (sp_vampire.go): %w", err)
		}

		//reduce health by units attack damage
		err = cardinal.UpdateComponent(world, id, func(health *comp.Health) *comp.Health {
			if health == nil {
				fmt.Printf("error retrieving Health component (sp_vampire.go)")
				return nil
			}
			if health.CurrentHP == 0 { // do not heal because unit will never die if its always healing at 0
				return nil
			}

			health.CurrentHP += vampire.healAmount
			if health.CurrentHP > UnitRegistry["vampire"].Health {
				health.CurrentHP = UnitRegistry["vampire"].Health
			}
			return health
		})

		if err != nil {
			return err
		}

		err := cardinal.SetComponent[comp.IntTracker](world, id, healCount)
		if err != nil {
			return fmt.Errorf("error setting int tracker component (sp_vampire.go): %w", err)
		}
	}

	return err
}

func vampireSpawnSP(world cardinal.WorldContext, id types.EntityID) error {
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

	err = cardinal.UpdateComponent(world, id, func(tracker *comp.IntTracker) *comp.IntTracker {
		if tracker == nil {
			fmt.Printf("error getting int tracker component (sp_vampire.go): ")
			return nil
		}
		tracker.Num = 0
		return tracker
	})
	if err != nil {
		return err
	}

	return err
}

// add Sp component to vampire unit
func vampireInitSP(world cardinal.WorldContext, id types.EntityID) error {

	err := cardinal.AddComponentTo[comp.IntTracker](world, id)
	if err != nil {
		return fmt.Errorf("error adding init component (sp_vampire.go): %w", err)
	}

	err = cardinal.SetComponent(world, id, &comp.IntTracker{Num: 0})
	if err != nil {
		return fmt.Errorf("error setting init component (sp_vampire.go): %w", err)
	}
	return nil
}
