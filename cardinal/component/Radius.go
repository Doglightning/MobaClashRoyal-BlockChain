package component

type UnitRadius struct {
	UnitRadius int `json:"UnitRadius"`
}

type AggroRadius struct {
	AggroRadius int `json:"AggroRadius"`
}

type AttackRadius struct {
	AttackRadius int `json:"AttackRadius"`
}

func (UnitRadius) Name() string {
	return "UnitRadius"
}

func (AggroRadius) Name() string {
	return "AggroRadius"
}

func (AttackRadius) Name() string {
	return "AttackRadius"
}
