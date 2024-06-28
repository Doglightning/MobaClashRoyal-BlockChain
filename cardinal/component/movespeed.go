package component

type Movespeed struct {
	CurrentMS float32 `json:"currentms"`
}

func (Movespeed) Name() string {
	return "Movespeed"
}
