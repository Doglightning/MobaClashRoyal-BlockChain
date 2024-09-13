package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// fireSpiritSpawnSP struct contains configuration for an fire spirit in terms of their shooting properties.
type fireSpiritSpawnSP struct {
	Hieght    float32
	BaseWidth float32
	Damage    float32
}

// NewArcherLadySP creates a new instance of archerLadySP with default settings.
func NewFireSpiritSpawnSP() *fireSpiritSpawnSP {
	return &fireSpiritSpawnSP{
		Hieght:    570,
		BaseWidth: 385,
		Damage:    3.5,
	}
}

func fireSpiritSpawn(world cardinal.WorldContext, id types.EntityID) error {
	//get fire spirit vars
	fireSprit := NewFireSpiritSpawnSP()

	//get team comp
	team, err := cardinal.GetComponent[comp.Team](world, id)
	if err != nil {
		return fmt.Errorf("error getting team component (class fireSpirit.go): %v", err)
	}

	matchID, err := cardinal.GetComponent[comp.MatchId](world, id)
	if err != nil {
		return fmt.Errorf("error getting position component (class fireSpirit.go): %v", err)
	}

	gameState, err := getGameStateGSS(world, matchID)
	if err != nil {
		return fmt.Errorf("(class fireSpirit.go): %v", err)
	}

	//get position comp
	hash, err := cardinal.GetComponent[comp.SpatialHash](world, gameState)
	if err != nil {
		return fmt.Errorf("error getting spatial hash compoenent(class fireSpirit.go): %v", err)
	}

	//get position comp
	pos, err := cardinal.GetComponent[comp.Position](world, id)
	if err != nil {
		return fmt.Errorf("error getting position component (class fireSpirit.go): %v", err)
	}

	//find the 3 points of the fire spirit AoE triangle attack
	apex, baseLeft, baseRight := CreateIsoscelesTriangle(Point{X: pos.PositionVectorX, Y: pos.PositionVectorY}, Point{X: pos.RotationVectorX, Y: pos.RotationVectorY}, fireSprit.Hieght, fireSprit.BaseWidth)

	//list of every point within the triangle
	points := RasterizeIsoscelesTriangle(apex, baseLeft, baseRight)

	// Define a map to track unique collisions
	collidedEntities := make(map[types.EntityID]bool)

	for _, p := range points {

		collList := CheckCollisionSpatialHashList(hash, p.X, p.Y, 1)
		for _, collID := range collList {
			collidedEntities[collID] = true
		}
	}

	// Iterate over each key in the map
	for collID := range collidedEntities {
		//get targets team
		targetTeam, err := cardinal.GetComponent[comp.Team](world, collID)
		if err != nil {
			fmt.Printf("error getting targets team compoenent (class fireSpirit.go): %v \n", err)
			continue
		}

		if team.Team != targetTeam.Team {

			// reduce health by units attack damage
			err = cardinal.UpdateComponent(world, collID, func(health *comp.Health) *comp.Health {
				if health == nil {
					fmt.Printf("error retrieving Health component (class fireSpirit.go) \n")
					return nil
				}
				health.CurrentHP -= fireSprit.Damage
				if health.CurrentHP < 0 {
					health.CurrentHP = 0 //never have negative health
				}
				return health
			})
			if err != nil {
				fmt.Printf("error updating health (class fireSpirit.go): %v \n", err)
				continue
			}

		}
	}

	//reset attack component
	err = cardinal.UpdateComponent(world, id, func(attack *comp.Attack) *comp.Attack {
		if attack == nil {
			fmt.Printf("error retrieving enemy attack component (class fireSpirit.go): \n")
			return nil
		}

		sp, err := cardinal.GetComponent[comp.Sp](world, id)
		if err != nil {
			fmt.Printf("error retrieving special power comp (class fireSpirit.go): ")
			return nil
		}

		if attack.Frame == sp.DamageEndFrame+1 {

			attack.State = "Default"
			if attack.Target == id {
				attack.Combat = false
			}

		}

		return attack
	})
	if err != nil {
		return fmt.Errorf("error updating attack comp (class fireSpirit.go): %v", err)
	}

	return nil
}

func fireSpiritResetCombat(world cardinal.WorldContext, id types.EntityID) error {

	//reset attack component
	err := cardinal.UpdateComponent(world, id, func(attack *comp.Attack) *comp.Attack {
		if attack == nil {
			fmt.Printf("error retrieving enemy attack component (fireSpiritResetCombat): \n")
			return nil
		}

		sp, err := cardinal.GetComponent[comp.Sp](world, id)
		if err != nil {
			fmt.Printf("error retrieving special power comp (fireSpiritResetCombat): ")
			return nil
		}

		if (attack.Frame < sp.DamageFrame && sp.Charged) || (attack.Frame > sp.DamageEndFrame && sp.Charged) {

			attack.Frame = 0
			attack.Combat = false
			attack.State = "Default"
			attack.Target = id
		} else {
			attack.State = "Channeling"
			attack.Target = id
		}

		return attack
	})
	if err != nil {
		return fmt.Errorf("error updating attack comp (fireSpiritResetCombat): %v", err)
	}

	return nil
}

// custom phase_attack.go logic
func FireSpiritAttack(world cardinal.WorldContext, id types.EntityID, atk *comp.Attack) error {

	//get Unit CC component
	cc, err := cardinal.GetComponent[comp.CC](world, id)
	if err != nil {
		fmt.Printf("error getting unit cc component (Fire Spirit Attack): %v", err)
	}

	if cc.Stun > 0 { //if unit stunned cannot attack
		return nil
	}

	//get special power component
	unitSp, err := cardinal.GetComponent[comp.Sp](world, id)
	if err != nil {
		return fmt.Errorf("error retrieving special power component (Fire Spirit Attack): %v", err)
	}

	//check if in a SP animation or a regular attack
	if atk.Frame == 0 && unitSp.CurrentSp >= unitSp.MaxSp { //In special power

		unitSp.Charged = true

	} else if atk.Frame == 0 && unitSp.CurrentSp < unitSp.MaxSp { // in regular attack
		unitSp.Charged = false
	}

	//if unit is in its damage frame and not charged
	if atk.Frame == atk.DamageFrame && !unitSp.Charged {
		unitSp.CurrentSp += unitSp.SpRate //increase sp after attack
		// make sure we are not over MaxSp
		if unitSp.CurrentSp >= unitSp.MaxSp {
			unitSp.CurrentSp = unitSp.MaxSp
		}
	}

	//if unit is in damage frames when charged
	if unitSp.DamageFrame <= atk.Frame && atk.Frame <= unitSp.DamageEndFrame && unitSp.Charged {

		//Shoot Fire >:D
		err = fireSpiritSpawn(world, id)
		if err != nil {
			return err
		}

		//return Sp to 0
		unitSp.CurrentSp = 0

	}

	//if target died in cast (self target) and attack frame is at end of animation or start (don't interupt the fire strike once its going even if target died)
	if (atk.Target == id && atk.Frame == unitSp.Rate-6) || atk.Target == id && atk.Frame < unitSp.DamageFrame {

		atk.State = "Default"
		atk.Combat = false

	}

	//if attack frame is at max and not sp charged  OR attack fram at sp max and charged
	if (atk.Frame >= atk.Rate && !unitSp.Charged) || (atk.Frame >= unitSp.Rate && unitSp.Charged) {
		atk.Frame = -1
	}

	atk.Frame++
	// set updated attack component
	if err := cardinal.SetComponent(world, id, atk); err != nil {
		return fmt.Errorf("error updating attack component (Fire Spirit Attack): %s ", err)
	}
	// set updated sp component
	if err := cardinal.SetComponent(world, id, unitSp); err != nil {
		return fmt.Errorf("error updating special power component (Fire Spirit Attack): %s ", err)
	}

	return nil
}
