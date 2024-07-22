package component

import "pkg.world.dev/world-engine/cardinal/types"

type Target struct {
	Target types.EntityID `json:"target"`
}

func (Target) Name() string {
	return "Target"
}
