package component

type CC struct {
	Stun int `json:"Stun"`
}

func (CC) Name() string {
	return "CC"
}
