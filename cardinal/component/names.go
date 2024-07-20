package component

type UnitName struct {
	UnitName string `json:"UnitName"`
}

type MapName struct {
	MapName string `json:"MapName"`
}

type SpName struct {
	SpName string `json:"SpName"`
}

func (UnitName) Name() string {
	return "UnitName"
}

func (MapName) Name() string {
	return "MapName"
}

func (SpName) Name() string {
	return "SpName"
}
