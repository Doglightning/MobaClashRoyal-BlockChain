package component

type Damage struct {
	Damage int `json:"Damage"`
}

func (Damage) Name() string {
	return "Damage"
}
