package system

import (
	comp "MobaClashRoyal/component"
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/search/filter"
	"pkg.world.dev/world-engine/cardinal/types"
)

// update struct
type Tower struct {
	healing float32
}

// update vars
func NewTower() *Tower {
	return &Tower{
		healing: 2,
	}
}

// Heal towers while converting
func TowerConverterSystem(world cardinal.WorldContext) error {
	// Filter for no HP
	stateFilter := cardinal.ComponentFilter(func(m comp.State) bool {
		return m.State == "Converting"
	})
	//for each tower still converting teams
	err := cardinal.NewSearch().Entity(
		filter.Contains(filter.Component[comp.StructureTag]())).
		Where(stateFilter).Each(world, func(id types.EntityID) bool {

		tower := NewTower()

		// increase tower hp until full
		err := cardinal.UpdateComponent(world, id, func(health *comp.Health) *comp.Health {
			if health == nil {
				fmt.Printf("error retrieving health component (tower conversion.go)")
				return nil
			}
			health.CurrentHP += tower.healing
			if health.CurrentHP >= health.MaxHP {
				health.CurrentHP = health.MaxHP

				// set tower state to Default
				err := cardinal.UpdateComponent(world, id, func(state *comp.State) *comp.State {
					if state == nil {
						fmt.Printf("error retrieving state component (tower conversion.go)")
						return nil
					}
					state.State = "Default"
					return state
				})

				if err != nil {
					fmt.Printf("error updating state component (tower conversion.go): %s", err)
					return health
				}

			}
			return health
		})

		if err != nil {
			fmt.Printf("error updating health component (tower conversion.go): %s", err)
			return false
		}

		return true
	})
	return err
}

// spawns projectile for tower basic attack
func towerAttack(world cardinal.WorldContext, id types.EntityID, atk *comp.Attack) error {
	//get units component
	unitPosition, matchID, mapName, unitName, err := GetComponents4[comp.Position, comp.MatchId, comp.MapName, comp.UnitName](world, id) //reusing
	if err != nil {
		return fmt.Errorf("tower components (class towerAttack.go): %v ", err)
	}
	//get next uid
	UID, err := getNextUID(world, matchID.MatchId)
	if err != nil {
		return fmt.Errorf("(class towerAttack.go): %v ", err)
	}
	//create projectile entity
	_, err = cardinal.Create(world,
		comp.MatchId{MatchId: matchID.MatchId},
		comp.UID{UID: UID},
		comp.Class{Class: "projectile"},
		comp.UnitName{UnitName: ProjectileRegistry[unitName.UnitName].Name},
		comp.Movespeed{CurrentMS: ProjectileRegistry[unitName.UnitName].Speed},
		comp.Position{
			PositionVectorX: unitPosition.PositionVectorX,
			PositionVectorY: unitPosition.PositionVectorY,
			PositionVectorZ: unitPosition.PositionVectorZ + ProjectileRegistry[unitName.UnitName].offSetZ,
			RotationVectorX: unitPosition.RotationVectorX,
			RotationVectorY: unitPosition.RotationVectorY,
			RotationVectorZ: unitPosition.RotationVectorZ},
		comp.MapName{MapName: mapName.MapName},
		comp.Attack{Target: atk.Target, Damage: StructureDataRegistry[unitName.UnitName].Damage},
		comp.Destroyed{Destroyed: false},
		comp.ProjectileTag{},
	)

	if err != nil {
		return fmt.Errorf("error spawning tower basic attack (class towerAttack.go): %v ", err)
	}

	return nil
}
