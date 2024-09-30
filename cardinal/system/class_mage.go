package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// update struct
type mageUpdateSP struct {
	frameCount int
}

// update vars
func NewMageUpdateSP() *mageUpdateSP {
	return &mageUpdateSP{
		frameCount: 25,
	}
}

// Stuns mages target
func MageSpawnSP(world cardinal.WorldContext, id types.EntityID) error {
	// get unit attack component
	unitAtk, err := cardinal.GetComponent[comp.Attack](world, id)
	if err != nil {
		return fmt.Errorf("error retrieving unit Attack component (class mage.go): %w", err)
	}

	//Add stun effect to target effects list
	err = cardinal.UpdateComponent(world, unitAtk.Target, func(effect *comp.EffectsList) *comp.EffectsList {
		if effect == nil {
			fmt.Printf("error retrieving effect list (class mage.go) \n")
			return nil
		}

		effect.EffectsList["MageStun"]++ //increase mage stun count on target

		return effect
	})
	if err != nil {
		return fmt.Errorf("error on effect list (class mage.go): %v", err)
	}

	//get matchid component
	matchID, err := cardinal.GetComponent[comp.MatchId](world, id)
	if err != nil {
		return fmt.Errorf("error getting matchID comp (class mage.go): %w", err)
	}
	//get new uid
	UID, err := getNextUID(world, matchID.MatchId)
	if err != nil {
		return fmt.Errorf("(class mage.go): %v - ", err)
	}
	//create stun entity attached to target
	_, err = cardinal.Create(world,
		comp.MatchId{MatchId: matchID.MatchId},
		comp.UID{UID: UID},
		comp.SpEntity{SpName: "MageSP"},
		comp.IntTracker{Num: 0}, //tracks duration
		comp.Target{Target: unitAtk.Target},
	)
	if err != nil {
		return fmt.Errorf("error creating stun entity (class mage.go): %v", err)
	}

	return err
}

// called every tick to update the active stuns to see if they have surpased duration limit
func MageUpdate(world cardinal.WorldContext, id types.EntityID) error {
	mage := NewMageUpdateSP() // get mage vars
	// get target id
	tarID, err := cardinal.GetComponent[comp.Target](world, id)
	if err != nil {
		return fmt.Errorf("error getting target comp (class mage.go): %w", err)
	}

	//get targets cc component
	cc, err := cardinal.GetComponent[comp.CC](world, tarID.Target)
	if err != nil {
		return fmt.Errorf("error getting target cc component(class mage.go): %w", err)
	}
	//if target not stunned then stun it (should not happen but just to be careful)
	if cc.Stun == 0 {
		cc.Stun++
	}

	// get tracker holding number of frames stun has been active for (stun count)
	stunCount, err := cardinal.GetComponent[comp.IntTracker](world, id)
	if err != nil {
		return fmt.Errorf("error retrieving int tracker component (class mage.go): %w", err)
	}
	stunCount.Num += 1                    // increase stun frame count
	if stunCount.Num >= mage.frameCount { //if stun count is greater than mage max stun count

		//remove stun of target
		err := cardinal.UpdateComponent(world, tarID.Target, func(effect *comp.EffectsList) *comp.EffectsList {
			if effect == nil {
				fmt.Printf("error retrieving effect list (class mage.go) \n")
				return nil
			}

			if list, ok := effect.EffectsList["MageStun"]; ok { // if key exists
				if list <= 1 { // if 1 or less of this buff active remove
					delete(effect.EffectsList, "MageStun")
				} else { // if more then 1 active reduce by 1
					effect.EffectsList["MageStun"] -= 1
				}
			}
			return effect
		})
		if err != nil {
			return err
		}

		// delete entity
		if err := cardinal.Remove(world, id); err != nil {
			return fmt.Errorf("error removing entity sp (class mage.go): %w", err)
		}

		cc.Stun-- //reduce stun count by 1
		if cc.Stun < 0 {
			cc.Stun = 0
		}

	} else { // else update stun count component
		if err := cardinal.SetComponent(world, id, stunCount); err != nil {
			return fmt.Errorf("error setting stun count (class mage.go): %w", err)
		}
	}
	//set cc component
	if err = cardinal.SetComponent(world, tarID.Target, cc); err != nil {
		return fmt.Errorf("error setting stun component on target (class mage.go): %w", err)
	}

	return err
}

// spawns projectile for mage basic attack
func mageAttack(world cardinal.WorldContext, id types.EntityID, atk *comp.Attack) error {
	//get units component
	unitPosition, matchID, mapName, unitName, err := GetComponents4[comp.Position, comp.MatchId, comp.MapName, comp.UnitName](world, id)
	if err != nil {
		return fmt.Errorf("unit components (class mageAttack.go): %v ", err)
	}
	//get next uid
	UID, err := getNextUID(world, matchID.MatchId)
	if err != nil {
		return fmt.Errorf("(class mageAttack.go): %v ", err)
	}
	// set offset to units mesh in client
	newX, newY := RelativeOffsetXY(unitPosition.PositionVectorX, unitPosition.PositionVectorY, unitPosition.RotationVectorX, unitPosition.RotationVectorY, ProjectileRegistry[unitName.UnitName].offSetX, ProjectileRegistry[unitName.UnitName].offSetY)
	//create projectile entity
	_, err = cardinal.Create(world,
		comp.MatchId{MatchId: matchID.MatchId},
		comp.UID{UID: UID},
		comp.UnitName{UnitName: ProjectileRegistry[unitName.UnitName].Name},
		comp.Movespeed{CurrentMS: ProjectileRegistry[unitName.UnitName].Speed},
		comp.Position{
			PositionVectorX: newX,
			PositionVectorY: newY,
			PositionVectorZ: unitPosition.PositionVectorZ + ProjectileRegistry[unitName.UnitName].offSetZ,
			RotationVectorX: unitPosition.RotationVectorX,
			RotationVectorY: unitPosition.RotationVectorY,
			RotationVectorZ: unitPosition.RotationVectorZ,
		},
		comp.MapName{MapName: mapName.MapName},
		comp.Attack{Target: atk.Target, Class: "projectile", Damage: UnitRegistry[unitName.UnitName].Damage},
		comp.Destroyed{Destroyed: false},
	)

	if err != nil {
		return fmt.Errorf("error spawning mage basic attack (class mageAttack.go): %v ", err)
	}

	return nil
}
