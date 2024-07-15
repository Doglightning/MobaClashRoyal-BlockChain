package component

type Health struct {
	CurrentHP float32 `json:"currenthp"`
	MaxHP     float32 `json:"maxhp"`
}

func (Health) Name() string {
	return "Health"
}
