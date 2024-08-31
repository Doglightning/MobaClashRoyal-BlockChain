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

	DmgSp     int
	SpRate    int
	CurrentSP int
	MaxSP     int
}

// registry of all units in game
var UnitRegistry = map[string]UnitType{
	"ArcherLady": {Class: "range", Health: 75, Damage: 30, AttackRate: 30, DamageFrame: 20, Speed: 50, Cost: 3, Radius: 50, AggroRadius: 1400, AttackRadius: 1200, DmgSp: 25, SpRate: 50, CurrentSP: 0, MaxSP: 100},
	"Vampire":    {Class: "melee", Health: 100, Damage: 10, AttackRate: 10, DamageFrame: 4, Speed: 50, Cost: 2, Radius: 80, AggroRadius: 1400, AttackRadius: 10, DmgSp: 10, SpRate: 25, CurrentSP: 0, MaxSP: 100},
}

type ProjectileType struct {
	Name  string
	Speed float32
}

// registry of all projectiles in game
var ProjectileRegistry = map[string]ProjectileType{
	"ArcherLady": {Name: "ArcherLadyArrow", Speed: 150},
}
