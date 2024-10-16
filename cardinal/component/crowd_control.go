package component

type CC struct {
	Stun      int  `json:"Stun"`
	KnockBack bool `json:"KnockBack"`
}

func (CC) Name() string {
	return "CC"
}

type KnockUp struct {
	CurrentHieght float32 `json:"CurrentHieght"`
	TargetHieght  float32 `json:"TargetHieght"`
	Speed         float32 `json:"Speed"`
	Damage        float32 `json:"Damage"`
	ApexReached   bool    `json:"ApexReached"`
}

func (KnockUp) Name() string {
	return "KnockUp"
}
