package component

type UID struct {
	UID int `json:"UID"`
}

func (UID) Name() string {
	return "UID"
}
