package component

type Sp struct {
	DmgSp     int `json:"DmgSp"`
	SpRate    int `json:"rate"`
	CurrentSp int `json:"CurrentSp"`
	MaxSp     int `json:"MaxSp"`

	Animation string `json:"Animation"`
}

func (Sp) Name() string {
	return "Sp"
}
