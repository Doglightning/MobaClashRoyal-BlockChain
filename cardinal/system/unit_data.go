package system

import "pkg.world.dev/world-engine/cardinal/types"

type UnitType struct {
	Name         string
	Class        string
	Health       float32
	Damage       float32
	AttackRate   int //tick based 5 = 5 ticks (100ms tickrate = 500ms attack rate)
	DamageFrame  int
	Speed        float32
	CenterOffset float32
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
	"ArcherLady": {Class: "ground", Health: 75, Damage: 30, AttackRate: 30, DamageFrame: 20, Speed: 50, Cost: 3, Radius: 50, AggroRadius: 1400, AttackRadius: 1200, CenterOffset: 150, DmgSp: 25, SpRate: 50, CurrentSP: 0, MaxSP: 100},
	"FireSpirit": {Class: "ground", Health: 100, Damage: 2.5, AttackRate: 20, DamageFrame: 13, Speed: 40, Cost: 2, Radius: 100, AggroRadius: 1400, AttackRadius: 350, CenterOffset: 150, DmgSp: 10, SpRate: 34, CurrentSP: 0, MaxSP: 100},
	"Mage":       {Class: "ground", Health: 75, Damage: 15, AttackRate: 20, DamageFrame: 8, Speed: 30, Cost: 3, Radius: 130, AggroRadius: 1400, AttackRadius: 1000, CenterOffset: 150, DmgSp: 25, SpRate: 50, CurrentSP: 0, MaxSP: 100},
	"Vampire":    {Class: "ground", Health: 100, Damage: 10, AttackRate: 10, DamageFrame: 4, Speed: 50, Cost: 2, Radius: 80, AggroRadius: 1400, AttackRadius: 10, CenterOffset: 150, DmgSp: 10, SpRate: 25, CurrentSP: 0, MaxSP: 100},
}

type SpType struct {
	AttackRate          int  `json:"attackrate"`  //tick based 5 Rate = 5 ticks (100ms tickrate = 500ms attack rate)
	DamageFrame         int  `json:"damageframe"` // Frame damage goes off (most animations have a wind down so the dmage goes off in the middle somewhere)
	DamageEndFrame      int  `json:"damageendframe"`
	StructureTargetable bool `json:"StructureTargetable"`
}

var SpRegistry = map[string]SpType{
	"ArcherLady": {AttackRate: 30, DamageFrame: 20, DamageEndFrame: 20, StructureTargetable: true},
	"FireSpirit": {AttackRate: 41, DamageFrame: 14, DamageEndFrame: 27, StructureTargetable: true},
	"Mage":       {AttackRate: 15, DamageFrame: 8, DamageEndFrame: 8, StructureTargetable: false},
	"Vampire":    {AttackRate: 10, DamageFrame: 4, DamageEndFrame: 4, StructureTargetable: true},
}

type ProjectileType struct {
	Name    string
	Speed   float32
	offSetX float32
	offSetY float32
	offSetZ float32
}

// registry of all projectiles in game
var ProjectileRegistry = map[string]ProjectileType{
	"ArcherLady": {Name: "ArcherLadyArrow", Speed: 150, offSetX: 20, offSetY: 28, offSetZ: 190},
	"Mage":       {Name: "MageBolt", Speed: 80, offSetX: 45, offSetY: 80, offSetZ: 307},
	"Base":       {Name: "BaseBolt", Speed: 150, offSetX: 0, offSetY: 0, offSetZ: 1000},
	"Tower":      {Name: "TowerBolt", Speed: 150, offSetX: 0, offSetY: 0, offSetZ: 1000},
}

type StructureData struct {
	Health       float32        `json:"health"`
	Radius       int            `json:"radius"`
	Damage       float32        `json:"damage"`
	AttackRate   int            `json:"attackrate"`  //tick based 5 Rate = 5 ticks (100ms tickrate = 500ms attack rate)
	DamageFrame  int            `json:"attackframe"` // Frame damage goes off (most animations have a wind down so the dmage goes off in the middle somewhere)
	AttackRadius int            `json:"AttackRadius"`
	AggroRadius  int            `json:"AggroRadius"`
	Target       types.EntityID `json:"target"`
	Class        string         `json:"class"`

	CenterOffset float32
}

// structures
var StructureDataRegistry = map[string]StructureData{
	"Base":  {Class: "structure", Health: 200, Radius: 240, Damage: 15, AttackRate: 20, DamageFrame: 10, AttackRadius: 1700, AggroRadius: 1700, CenterOffset: 230},
	"Tower": {Class: "structure", Health: 150, Radius: 150, Damage: 15, AttackRate: 20, DamageFrame: 10, AttackRadius: 1700, AggroRadius: 1700, CenterOffset: 230},
}
