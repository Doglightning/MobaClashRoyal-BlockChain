package component

import "pkg.world.dev/world-engine/cardinal/types"

type SpatialCell struct {
	UnitIDs    []types.EntityID `json:"UnitIDs"`
	PositionsX []float32        `json:"PositionsX"`
	PositionsY []float32        `json:"PositionsY"`
	Radii      []int            `json:"Radius"`
	Team       []string         `json:"Team"`
	Type       []string         `json:"Type"`
}

// DirectionMap is the component that holds mapping from coordinates to vectors directly within.
type SpatialHash struct {
	Cells    map[string]SpatialCell `json:"Cells"`
	CellSize int                    `json:"CellSize"`
	StartX   float32                `json:"StartX"`
	StartY   float32                `json:"StartY"`
}

func (SpatialHash) Name() string {
	return "SpatialHash"
}

func (SpatialCell) Name() string {
	return "SpatialCell"
}
