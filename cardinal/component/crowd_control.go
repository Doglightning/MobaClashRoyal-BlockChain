package component

type CC struct {
	Stun bool `json:"Stun"`
}

func (CC) Name() string {
	return "CC"
}
