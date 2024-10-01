package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// SetComponents2 sets multiple components (2 in this case) to an entity.
func SetComponents2[T1, T2 types.Component](wCtx cardinal.WorldContext, id types.EntityID, c1 *T1, c2 *T2) error {
	if err := cardinal.SetComponent(wCtx, id, c1); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T1), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c2); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T2), err)
	}

	return nil
}

// SetComponents3 sets multiple components (3 in this case) to an entity.
func SetComponents3[T1, T2, T3 types.Component](wCtx cardinal.WorldContext, id types.EntityID, c1 *T1, c2 *T2, c3 *T3) error {
	if err := cardinal.SetComponent(wCtx, id, c1); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T1), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c2); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T2), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c3); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T3), err)
	}

	return nil
}

// SetComponents4 sets multiple components (4 in this case) to an entity.
func SetComponents4[T1, T2, T3, T4 types.Component](wCtx cardinal.WorldContext, id types.EntityID, c1 *T1, c2 *T2, c3 *T3, c4 *T4) error {
	if err := cardinal.SetComponent(wCtx, id, c1); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T1), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c2); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T2), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c3); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T3), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c4); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T4), err)
	}

	return nil
}

// SetComponents5 sets multiple components (5 in this case) to an entity.
func SetComponents5[T1, T2, T3, T4, T5 types.Component](wCtx cardinal.WorldContext, id types.EntityID, c1 *T1, c2 *T2, c3 *T3, c4 *T4, c5 *T5) error {
	if err := cardinal.SetComponent(wCtx, id, c1); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T1), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c2); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T2), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c3); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T3), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c4); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T4), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c5); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T5), err)
	}

	return nil
}

// SetComponents6 sets multiple components (6 in this case) to an entity.
func SetComponents6[T1, T2, T3, T4, T5, T6 types.Component](wCtx cardinal.WorldContext, id types.EntityID, c1 *T1, c2 *T2, c3 *T3, c4 *T4, c5 *T5, c6 *T6) error {
	if err := cardinal.SetComponent(wCtx, id, c1); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T1), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c2); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T2), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c3); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T3), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c4); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T4), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c5); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T5), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c6); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T6), err)
	}

	return nil
}

// SetComponents7 sets multiple components (7 in this case) to an entity.
func SetComponents7[T1, T2, T3, T4, T5, T6, T7 types.Component](wCtx cardinal.WorldContext, id types.EntityID, c1 *T1, c2 *T2, c3 *T3, c4 *T4, c5 *T5, c6 *T6, c7 *T7) error {
	if err := cardinal.SetComponent(wCtx, id, c1); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T1), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c2); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T2), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c3); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T3), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c4); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T4), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c5); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T5), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c6); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T6), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c7); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T7), err)
	}

	return nil
}

// SetComponents8 sets multiple components (8 in this case) to an entity.
func SetComponents8[T1, T2, T3, T4, T5, T6, T7, T8 types.Component](wCtx cardinal.WorldContext, id types.EntityID, c1 *T1, c2 *T2, c3 *T3, c4 *T4, c5 *T5, c6 *T6, c7 *T7, c8 *T8) error {
	if err := cardinal.SetComponent(wCtx, id, c1); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T1), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c2); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T2), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c3); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T3), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c4); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T4), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c5); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T5), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c6); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T6), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c7); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T7), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c8); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T8), err)
	}

	return nil
}

// SetComponents9 sets multiple components (9 in this case) to an entity.
func SetComponents9[T1, T2, T3, T4, T5, T6, T7, T8, T9 types.Component](wCtx cardinal.WorldContext, id types.EntityID, c1 *T1, c2 *T2, c3 *T3, c4 *T4, c5 *T5, c6 *T6, c7 *T7, c8 *T8, c9 *T9) error {
	if err := cardinal.SetComponent(wCtx, id, c1); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T1), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c2); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T2), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c3); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T3), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c4); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T4), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c5); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T5), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c6); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T6), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c7); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T7), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c8); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T8), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c9); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T9), err)
	}

	return nil
}

// SetComponents10 sets multiple components (10 in this case) to an entity.
func SetComponents10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 types.Component](wCtx cardinal.WorldContext, id types.EntityID, c1 *T1, c2 *T2, c3 *T3, c4 *T4, c5 *T5, c6 *T6, c7 *T7, c8 *T8, c9 *T9, c10 *T10) error {
	if err := cardinal.SetComponent(wCtx, id, c1); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T1), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c2); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T2), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c3); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T3), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c4); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T4), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c5); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T5), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c6); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T6), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c7); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T7), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c8); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T8), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c9); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T9), err)
	}
	if err := cardinal.SetComponent(wCtx, id, c10); err != nil {
		return fmt.Errorf("error setting component of type %T: %v", *new(T10), err)
	}

	return nil
}
