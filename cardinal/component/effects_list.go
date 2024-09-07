package component

type EffectsList struct {
	EffectsList map[string]int `json:"EffectsList"`
}

func (EffectsList) Name() string {
	return "EffectsList"
}
