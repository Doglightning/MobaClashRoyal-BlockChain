package system

import (
	comp "MobaClashRoyal/component"

	"pkg.world.dev/world-engine/cardinal/search/filter"
)

//filters to get specific unit types

func ProjectileFilter() (filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper, filter.ComponentWrapper) {
	return filter.Component[comp.MatchId](), filter.Component[comp.UID](), filter.Component[comp.UnitName](), filter.Component[comp.Movespeed](), filter.Component[comp.Position](), filter.Component[comp.MapName](), filter.Component[comp.Attack](), filter.Component[comp.Destroyed]()
}
