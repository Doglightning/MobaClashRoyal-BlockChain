package component

import "pkg.world.dev/world-engine/cardinal/types"

type Sp struct {
	DmgSp          int `json:"DmgSp"`
	SpRate         int `json:"sprate"`
	CurrentSp      int `json:"CurrentSp"`
	MaxSp          int `json:"MaxSp"`
	Rate           int `json:"rate"`        //tick based 5 Rate = 5 ticks (100ms tickrate = 500ms attack rate)
	DamageFrame    int `json:"damageframe"` // Frame damage goes off (most animations have a wind down so the dmage goes off in the middle somewhere)
	DamageEndFrame int `json:"damageendframe"`
	AttackRadius   int `json:"AttackRadius"`

	Target              types.EntityID `json:"target"`
	Combat              bool           `json:"Combat"`
	Charged             bool           `json:"Charged"`
	StructureTargetable bool           `json:"StructureTargetable"`
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
