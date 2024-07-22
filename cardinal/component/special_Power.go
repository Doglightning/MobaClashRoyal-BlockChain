package component

type Sp struct {
	DmgSp     int `json:"DmgSp"`
	SpRate    int `json:"rate"`
	CurrentSp int `json:"CurrentSp"`
	MaxSp     int `json:"MaxSp"`

	Animation string `json:"Animation"`
}

type SpEntity struct {
	SpName string `json:"SpName"`
}

func (Sp) Name() string {
	return "Sp"
}

func (SpEntity) Name() string {
	return "SpEntity"
}
