package system

import (
	comp "MobaClashRoyal/component"

	"pkg.world.dev/world-engine/cardinal/search/filter"
)

func UnitFilters() (filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper) {
	return filter.Component[comp.MatchId](), filter.Component[comp.UID](), filter.Component[comp.Team](), filter.Component[comp.UnitHealth](), filter.Component[comp.Position](), filter.Component[comp.Movespeed](), filter.Component[comp.UnitName](), filter.Component[comp.MapName](), filter.Component[comp.Distance](), filter.Component[comp.UnitRadius](), filter.Component[comp.AttackRadius](), filter.Component[comp.Attack]()
}

func GameStateFilters() (filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper) {
	return filter.Component[comp.MatchId](), filter.Component[comp.UID](), filter.Component[comp.Player1](), filter.Component[comp.Player2](), filter.Component[comp.SpatialHash]()
}

func StructureFilters() (filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper) {
	return filter.Component[comp.MatchId](), filter.Component[comp.UID](), filter.Component[comp.Team](), filter.Component[comp.UnitHealth](), filter.Component[comp.UnitName](), filter.Component[comp.MapName](), filter.Component[comp.Position](), filter.Component[comp.UnitRadius]()
}
