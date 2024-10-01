package component

type UnitTag struct {
}

func (UnitTag) Name() string {
	return "UnitTag"
}

type StructureTag struct {
}

func (StructureTag) Name() string {
	return "StructureTag"
}

type ProjectileTag struct {
}

func (ProjectileTag) Name() string {
	return "ProjectileTag"
}

type GameStateTag struct {
}

func (GameStateTag) Name() string {
	return "GameStateTag"
}
