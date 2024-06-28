package component

type UnitHealth struct {
	CurrentHP float32 `json:"currenthp"`
	MaxHP     float32 `json:"maxhp"`
}

func (UnitHealth) Name() string {
	return "UnitHealth"
}
