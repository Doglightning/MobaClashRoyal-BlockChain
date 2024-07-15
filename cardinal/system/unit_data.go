package system

type UnitType struct {
	Name         string
	Class        string
	Health       float32
	Damage       int
	AttackRate   int //tick based 5 = 5 ticks (100ms tickrate = 500ms attack rate)
	DamageFrame  int
	Target       int
	Speed        float32
	Cost         int
	Radius       int
	AggroRadius  int
	AttackRadius int
}

var UnitRegistry = map[string]UnitType{
	"Vampire":    {Class: "melee", Health: 100, Damage: 10, AttackRate: 10, DamageFrame: 4, Speed: 50, Cost: 50, Radius: 80, AggroRadius: 1400, AttackRadius: 10},
	"ArcherLady": {Class: "range", Health: 75, Damage: 35, AttackRate: 30, DamageFrame: 20, Speed: 50, Cost: 50, Radius: 50, AggroRadius: 1400, AttackRadius: 1200},
}

type ProjectileType struct {
	Name  string
	Speed float32
}

var ProjectileRegistry = map[string]ProjectileType{
	"ArcherLady": {Name: "ArcherLadyArrow", Speed: 150},
}
