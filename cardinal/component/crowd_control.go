package component

type CC struct {
	Stun      int  `json:"Stun"`
	KnockBack bool `json:"KnockBack"`
}

func (CC) Name() string {
	return "CC"
}
