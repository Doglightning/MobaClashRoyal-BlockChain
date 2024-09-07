package system

import (
	comp "MobaClashRoyal/component"

	"pkg.world.dev/world-engine/cardinal/search/filter"
)

//filters to get specific unit types

func UnitFilters() (filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper) {
	return filter.Component[comp.MatchId](), filter.Component[comp.UID](), filter.Component[comp.Team](), filter.Component[comp.Health](), filter.Component[comp.Position](), filter.Component[comp.Movespeed](), filter.Component[comp.UnitName](), filter.Component[comp.MapName](), filter.Component[comp.Distance](), filter.Component[comp.UnitRadius](), filter.Component[comp.Attack](), filter.Component[comp.Sp](), filter.Component[comp.CenterOffset](), filter.Component[comp.CC](), filter.Component[comp.EffectsList]()
}

func GameStateFilters() (filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper) {
	return filter.Component[comp.MatchId](), filter.Component[comp.UID](), filter.Component[comp.Player1](), filter.Component[comp.Player2](), filter.Component[comp.SpatialHash]()
}

func StructureFilters() (filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper) {
	return filter.Component[comp.MatchId](), filter.Component[comp.UID](), filter.Component[comp.Team](), filter.Component[comp.Health](), filter.Component[comp.UnitName](), filter.Component[comp.MapName](), filter.Component[comp.Position](), filter.Component[comp.UnitRadius](), filter.Component[comp.State](), filter.Component[comp.Attack](), filter.Component[comp.CenterOffset]()
}

func ProjectileFilters() (filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper) {
	return filter.Component[comp.MatchId](), filter.Component[comp.UID](), filter.Component[comp.UnitName](), filter.Component[comp.Movespeed](), filter.Component[comp.Position](), filter.Component[comp.MapName](), filter.Component[comp.Attack](), filter.Component[comp.Destroyed]()
}
