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

	targetID, err := cardinal.GetComponent[comp.Target](world, id)
	if err != nil {
		return fmt.Errorf("error getting attack comp (sp_vampire.go): %w", err)
	}

	targetHP, err := cardinal.GetComponent[comp.Health](world, targetID.Target)
	if err != nil {
		return fmt.Errorf("error getting health comp (sp_vampire.go): %w", err)
	}

	if targetHP.CurrentHP != 0 { // do not heal because unit will never die if its always healing at 0
		targetHP.CurrentHP += vampire.healAmount
		if targetHP.CurrentHP > UnitRegistry["Vampire"].Health {
			targetHP.CurrentHP = UnitRegistry["Vampire"].Health
		}

		if err := cardinal.SetComponent(world, targetID.Target, targetHP); err != nil {
			return fmt.Errorf("error setting target health comp (sp_vampire.go): %w", err)
		}
	}

	healCount, err := cardinal.GetComponent[comp.IntTracker](world, id)
	if err != nil {
		return fmt.Errorf("error retrieving int tracker component (sp_vampire.go): %w", err)
	}
	healCount.Num += 1
	if healCount.Num >= vampire.healCount {
		// remove entity
		if err := cardinal.Remove(world, id); err != nil {
			return fmt.Errorf("error removing entity sp (sp_vampire.go): %w", err)
		}
	} else {
		if err := cardinal.SetComponent(world, id, healCount); err != nil {
			return fmt.Errorf("error setting heal count (sp_vampire.go): %w", err)
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

	matchID, err := cardinal.GetComponent[comp.MatchId](world, id)
	if err != nil {
		return fmt.Errorf("error getting matchID comp (sp_vampire.go): %w", err)
	}

	UID, err := getNextUID(world, matchID.MatchId)
	if err != nil {
		return fmt.Errorf("(sp_vampire.go): %v - ", err)
	}
	//create projectile entity
	cardinal.Create(world,
		comp.MatchId{MatchId: matchID.MatchId},
		comp.UID{UID: UID},
		comp.SpEntity{SpName: "VampireSP"},
		comp.IntTracker{Num: 0},
		comp.Target{Target: id},
	)

	return err
}
