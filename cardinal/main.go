package main

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"pkg.world.dev/world-engine/cardinal"

	"MobaClashRoyal/component"
	"MobaClashRoyal/msg"
	"MobaClashRoyal/query"
	"MobaClashRoyal/system"
)

func main() {
	w, err := cardinal.NewWorld(cardinal.WithDisableSignatureVerification(), cardinal.WithTickChannel(time.Tick(100*time.Millisecond)))
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	MustInitWorld(w)

	Must(w.StartGame())
}

// MustInitWorld registers all components, messages, queries, and systems. This initialization happens in a helper
// function so that this can be used directly in tests.
func MustInitWorld(w *cardinal.World) {
	// Register components
	// NOTE: You must register your components here for it to be accessible.
	Must(
		cardinal.RegisterComponent[component.Attack](w),
		cardinal.RegisterComponent[component.Class](w),
		cardinal.RegisterComponent[component.Damage](w),
		cardinal.RegisterComponent[component.Destroyed](w),
		cardinal.RegisterComponent[component.DirectionMap](w),
		cardinal.RegisterComponent[component.Distance](w),
		cardinal.RegisterComponent[component.MapName](w),
		cardinal.RegisterComponent[component.MatchId](w),
		cardinal.RegisterComponent[component.Movespeed](w),
		cardinal.RegisterComponent[component.Player1](w),
		cardinal.RegisterComponent[component.Player2](w),
		cardinal.RegisterComponent[component.Position](w),
		cardinal.RegisterComponent[component.SpatialCell](w),
		cardinal.RegisterComponent[component.SpatialHash](w),
		cardinal.RegisterComponent[component.Sp](w),
		cardinal.RegisterComponent[component.SpEntity](w),
		cardinal.RegisterComponent[component.SpName](w),
		cardinal.RegisterComponent[component.Target](w),
		cardinal.RegisterComponent[component.Team](w),
		cardinal.RegisterComponent[component.UID](w),
		cardinal.RegisterComponent[component.Health](w),
		cardinal.RegisterComponent[component.UnitName](w),
		cardinal.RegisterComponent[component.UnitRadius](w),
		cardinal.RegisterComponent[component.IntTracker](w),
		cardinal.RegisterComponent[component.State](w),
		cardinal.RegisterComponent[component.CenterOffset](w),
		cardinal.RegisterComponent[component.CC](w),
		cardinal.RegisterComponent[component.EffectsList](w),
		cardinal.RegisterComponent[component.UnitTag](w),
		cardinal.RegisterComponent[component.StructureTag](w),
		cardinal.RegisterComponent[component.ProjectileTag](w),
	)

	// Register messages (user action)
	// NOTE: You must register your transactions here for it to be executed.
	Must(
		cardinal.RegisterMessage[msg.CreateMatchMsg, msg.CreateMatchResult](w, "create-match"),
		cardinal.RegisterMessage[msg.CreateUnitMsg, msg.CreateUnitResult](w, "create-unit"),
		cardinal.RegisterMessage[msg.RemoveAllEntitiesMsg, msg.RemoveAllEntitiesResult](w, "remove-all-entities"),
		cardinal.RegisterMessage[msg.RemoveUnitMsg, msg.RemoveUnitResult](w, "remove-list"),
	)

	// Register queries
	// NOTE: You must register your queries here for it to be accessible.
	Must(
		cardinal.RegisterQuery[query.MatchIdRequest, query.TeamStateResponse](w, "team-state", query.TeamState),
		cardinal.RegisterQuery[query.UnitMatchIdRequest, query.UnitStateResponse](w, "game-state", query.GameState),
		cardinal.RegisterQuery[query.PSMatchIdRequest, query.PlayerStateResponse](w, "player-state", query.PlayerState),
	)

	// Each system executes deterministically in the order they are added.
	// This is a neat feature that can be strategically used for systems that depends on the order of execution.
	// For example, you may want to run the attack system before the regen system
	// so that the player's HP is subtracted (and player killed if it reaches 0) before HP is regenerated.
	Must(cardinal.RegisterSystems(w,
		system.RemoveAllEntitiesMsgSystem,
		system.GameStateSpawnerSystem,

		system.GoldGeneration, //prespawn phase
		system.TowerConverterSystem,
		system.UnitSpawnerSystem,  //spawn phase
		system.UnitMovementSystem, //move phase
		system.ProjectileMovementSystem,
		system.CombatCheckSystem, //pre attack phase
		system.AttackPhaseSystem,
		system.SpUpdater,
		system.DestroyerSystem,   //destroy phase
		system.RemovalListSystem, //client replication
		system.WinCondition,      // game over
	))

	// Must(cardinal.RegisterInitSystems(w,
	// 	system.SpawnMaps,
	// ))
}

func Must(err ...error) {
	e := errors.Join(err...)
	if e != nil {
		log.Fatal().Err(e).Msg("")
	}
}
