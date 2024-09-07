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
		healAmount: 1.2,
	}
}

// updates SP entity per tick
func vampireUpdateSP(world cardinal.WorldContext, id types.EntityID) error {
	vampire := NewVampireSP() // get vampire vars
	// get target id
	targetID, err := cardinal.GetComponent[comp.Target](world, id)
	if err != nil {
		return fmt.Errorf("error getting attack comp (sp_vampire.go): %w", err)
	}
	// get targets health component
	targetHP, err := cardinal.GetComponent[comp.Health](world, targetID.Target)
	if err != nil {
		return fmt.Errorf("error getting health comp (sp_vampire.go): %w", err)
	}

	if targetHP.CurrentHP != 0 { // do not heal because unit will never die if its always healing at 0
		targetHP.CurrentHP += vampire.healAmount //heal unit

		//check if unit being spawned exsists in the unit registry
		unitType, exsist := UnitRegistry["Vampire"]
		if !exsist {
			return fmt.Errorf("vampire type not found in registry (sp_vampire.go)")
		}
		if targetHP.CurrentHP > unitType.Health { //cap healing at vampire max hp
			targetHP.CurrentHP = unitType.Health
		}
		//update health component
		if err := cardinal.SetComponent(world, targetID.Target, targetHP); err != nil {
			return fmt.Errorf("error setting target health comp (sp_vampire.go): %w", err)
		}
	}
	// get tracker holding number of frames heal has gone off (heal count)
	healCount, err := cardinal.GetComponent[comp.IntTracker](world, id)
	if err != nil {
		return fmt.Errorf("error retrieving int tracker component (sp_vampire.go): %w", err)
	}
	healCount.Num += 1                      // increase heal frame count
	if healCount.Num >= vampire.healCount { //if heal count is greater than vampire max heal count

		//remove heal spiral effect to the effects list
		err := cardinal.UpdateComponent(world, targetID.Target, func(effect *comp.EffectsList) *comp.EffectsList {
			if effect == nil {
				fmt.Printf("error retrieving effect list (sp_vampire.go) \n")
				return nil
			}

			if list, ok := effect.EffectsList["HealSpiral"]; ok { // if key exists
				if list <= 1 { // if 1 or less of this buff active remove
					delete(effect.EffectsList, "HealSpiral")
				} else { // if more then 1 active reduce by 1
					effect.EffectsList["HealSpiral"] -= 1
				}
			}
			return effect
		})
		if err != nil {
			return err
		}

		// remove entity
		if err := cardinal.Remove(world, id); err != nil {
			return fmt.Errorf("error removing entity sp (sp_vampire.go): %w", err)
		}
	} else { // else update heal count component
		if err := cardinal.SetComponent(world, id, healCount); err != nil {
			return fmt.Errorf("error setting heal count (sp_vampire.go): %w", err)
		}
	}

	return err
}

// spawning the vampire special power
func vampireSpawnSP(world cardinal.WorldContext, id types.EntityID) error {

	//Add heal spiral effect to the effects list
	err := cardinal.UpdateComponent(world, id, func(effect *comp.EffectsList) *comp.EffectsList {
		if effect == nil {
			fmt.Printf("error retrieving effect list (sp_vampire.go) \n")
			return nil
		}

		effect.EffectsList["HealSpiral"]++

		return effect
	})
	if err != nil {
		return fmt.Errorf("error on effect list (class vampire.go): %v", err)
	}

	// get unit attack component
	unitAtk, err := cardinal.GetComponent[comp.Attack](world, id)
	if err != nil {
		return fmt.Errorf("error retrieving unit Attack component (sp_vampire.go): %w", err)
	}

	err = vampireAttack(world, unitAtk)

	if err != nil {
		return err
	}
	//get matchid component
	matchID, err := cardinal.GetComponent[comp.MatchId](world, id)
	if err != nil {
		return fmt.Errorf("error getting matchID comp (sp_vampire.go): %w", err)
	}
	//get new uid
	UID, err := getNextUID(world, matchID.MatchId)
	if err != nil {
		return fmt.Errorf("(sp_vampire.go): %v - ", err)
	}
	//create healing buff entity
	_, err = cardinal.Create(world,
		comp.MatchId{MatchId: matchID.MatchId},
		comp.UID{UID: UID},
		comp.SpEntity{SpName: "VampireSP"},
		comp.IntTracker{Num: 0},
		comp.Target{Target: id},
	)
	if err != nil {
		return fmt.Errorf("error creating healing entity (sp_vampire.go): %v", err)
	}

	return err
}

func vampireAttack(world cardinal.WorldContext, atk *comp.Attack) error {
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
