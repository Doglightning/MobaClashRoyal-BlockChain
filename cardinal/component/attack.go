package component

import "pkg.world.dev/world-engine/cardinal/types"

type Attack struct {
	Combat      bool           `json:"combat"`
	Damage      int            `json:"damage"`
	Rate        int            `json:"rate"`        //tick based 5 Rate = 5 ticks (100ms tickrate = 500ms attack rate)
	Frame       int            `json:"frame"`       //current attack frame ex. frame 0-5 for a 5 rate
	DamageFrame int            `json:"attackframe"` // Frame damage goes off (most animations have a wind down so the dmage goes off in the middle somewhere)
	Target      types.EntityID `json:"target"`
}

func (Attack) Name() string {
	return "Attack"
}
