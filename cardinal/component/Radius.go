package component

type UnitRadius struct {
	UnitRadius int `json:"UnitRadius"`
}

func (UnitRadius) Name() string {
	return "UnitRadius"
}
