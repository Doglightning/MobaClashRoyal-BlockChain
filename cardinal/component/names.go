package component

type UnitName struct {
	UnitName string `json:"UnitName"`
}

type MapName struct {
	MapName string `json:"MapName"`
}

func (UnitName) Name() string {
	return "UnitName"
}

func (MapName) Name() string {
	return "MapName"
}
