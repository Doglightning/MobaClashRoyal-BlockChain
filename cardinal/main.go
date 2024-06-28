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
		cardinal.RegisterComponent[component.DirectionMap](w),
		cardinal.RegisterComponent[component.Distance](w),
		cardinal.RegisterComponent[component.GridUtils](w),
		cardinal.RegisterComponent[component.MapName](w),
		cardinal.RegisterComponent[component.MatchId](w),
		cardinal.RegisterComponent[component.Movespeed](w),
		cardinal.RegisterComponent[component.Player1](w),
		cardinal.RegisterComponent[component.Player2](w),
		cardinal.RegisterComponent[component.Position](w),
		cardinal.RegisterComponent[component.SizeCircle](w),
		cardinal.RegisterComponent[component.SpatialCell](w),
		cardinal.RegisterComponent[component.SpatialHash](w),
		cardinal.RegisterComponent[component.Team](w),
		cardinal.RegisterComponent[component.UID](w),
		cardinal.RegisterComponent[component.UnitHealth](w),
		cardinal.RegisterComponent[component.UnitName](w),
	)

	// Register messages (user action)
	// NOTE: You must register your transactions here for it to be executed.
	Must(
		cardinal.RegisterMessage[msg.CreateMatchMsg, msg.CreateMatchResult](w, "create-match"),
		cardinal.RegisterMessage[msg.CreateUnitMsg, msg.CreateUnitResult](w, "create-unit"),
		cardinal.RegisterMessage[msg.RemoveAllEntitiesMsg, msg.RemoveAllEntitiesResult](w, "remove-all-entities"),
	)

	// Register queries
	// NOTE: You must register your queries here for it to be accessible.
	Must(
		cardinal.RegisterQuery[query.MatchIdRequest, query.TeamStateResponse](w, "team-state", query.TeamState),
		cardinal.RegisterQuery[query.UnitMatchIdRequest, query.UnitStateResponse](w, "unit-state", query.UnitState),
	)

	// Each system executes deterministically in the order they are added.
	// This is a neat feature that can be strategically used for systems that depends on the order of execution.
	// For example, you may want to run the attack system before the regen system
	// so that the player's HP is subtracted (and player killed if it reaches 0) before HP is regenerated.
	Must(cardinal.RegisterSystems(w,
		system.RemoveAllEntitiesSystem,
		system.MatchSpawnerSystem,
		system.UnitSpawnerSystem,
		system.UnitMovementSystem,
	))

	Must(cardinal.RegisterInitSystems(w,
		system.SpawnMaps,
	))
}

func Must(err ...error) {
	e := errors.Join(err...)
	if e != nil {
		log.Fatal().Err(e).Msg("")
	}
}
