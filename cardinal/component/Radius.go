package component

type UnitRadius struct {
	UnitRadius int `json:"UnitRadius"`
}

type AttackRadius struct {
	AttackRadius int `json:"AttackRadius"`
}

func (UnitRadius) Name() string {
	return "UnitRadius"
}

func (AttackRadius) Name() string {
	return "AttackRadius"
}
