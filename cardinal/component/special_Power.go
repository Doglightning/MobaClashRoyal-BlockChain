package component

type Sp struct {
	SpRate    int `json:"rate"`
	CurrentSp int `json:"CurrentSp"`
	MaxSp     int `json:"MaxSp"`
}

func (Sp) Name() string {
	return "Sp"
}
