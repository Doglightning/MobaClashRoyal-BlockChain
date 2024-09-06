package component

type Sp struct {
	DmgSp       int `json:"DmgSp"`
	SpRate      int `json:"sprate"`
	CurrentSp   int `json:"CurrentSp"`
	MaxSp       int `json:"MaxSp"`
	Rate        int `json:"rate"`        //tick based 5 Rate = 5 ticks (100ms tickrate = 500ms attack rate)
	DamageFrame int `json:"attackframe"` // Frame damage goes off (most animations have a wind down so the dmage goes off in the middle somewhere)

	Charged             bool `json:"Charged"`
	StructureTargetable bool `json:"StructureTargetable"`

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
