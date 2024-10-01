package system

import (
	"fmt"

	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/types"
)

// GetComponents2 fetches multiple components (2 in this case) from an entity.
func GetComponents2[T1, T2 types.Component](wCtx cardinal.WorldContext, unitID types.EntityID) (*T1, *T2, error) {
	// Get first component
	c1, err := cardinal.GetComponent[T1](wCtx, unitID)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T1), err)
	}
	// Get second component
	c2, err := cardinal.GetComponent[T2](wCtx, unitID)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T2), err)
	}

	return c1, c2, nil
}

// GetComponents3 fetches multiple components (3 in this case) from an entity.
func GetComponents3[T1, T2, T3 types.Component](wCtx cardinal.WorldContext, unitID types.EntityID) (*T1, *T2, *T3, error) {
	// Get first component
	c1, err := cardinal.GetComponent[T1](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T1), err)
	}
	// Get second component
	c2, err := cardinal.GetComponent[T2](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T2), err)
	}
	// Get third component
	c3, err := cardinal.GetComponent[T3](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T3), err)
	}

	return c1, c2, c3, nil
}

// GetComponents4 fetches multiple components (4 in this case) from an entity.
func GetComponents4[T1, T2, T3, T4 types.Component](wCtx cardinal.WorldContext, unitID types.EntityID) (*T1, *T2, *T3, *T4, error) {
	// Get first component
	c1, err := cardinal.GetComponent[T1](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T1), err)
	}
	// Get second component
	c2, err := cardinal.GetComponent[T2](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T2), err)
	}
	// Get third component
	c3, err := cardinal.GetComponent[T3](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T3), err)
	}
	// Get fourth component
	c4, err := cardinal.GetComponent[T4](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T4), err)
	}

	return c1, c2, c3, c4, nil
}

// GetComponents5 fetches multiple components (5 in this case) from an entity.
func GetComponents5[T1, T2, T3, T4, T5 types.Component](wCtx cardinal.WorldContext, unitID types.EntityID) (*T1, *T2, *T3, *T4, *T5, error) {
	// Get first component
	c1, err := cardinal.GetComponent[T1](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T1), err)
	}
	// Get second component
	c2, err := cardinal.GetComponent[T2](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T2), err)
	}
	// Get third component
	c3, err := cardinal.GetComponent[T3](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T3), err)
	}
	// Get fourth component
	c4, err := cardinal.GetComponent[T4](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T4), err)
	}
	// Get fifth component
	c5, err := cardinal.GetComponent[T5](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T5), err)
	}

	return c1, c2, c3, c4, c5, nil
}

// GetComponents6 fetches multiple components (6 in this case) from an entity.
func GetComponents6[T1, T2, T3, T4, T5, T6 types.Component](wCtx cardinal.WorldContext, unitID types.EntityID) (*T1, *T2, *T3, *T4, *T5, *T6, error) {
	// Get first component
	c1, err := cardinal.GetComponent[T1](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T1), err)
	}
	// Get second component
	c2, err := cardinal.GetComponent[T2](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T2), err)
	}
	// Get third component
	c3, err := cardinal.GetComponent[T3](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T3), err)
	}
	// Get fourth component
	c4, err := cardinal.GetComponent[T4](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T4), err)
	}
	// Get fifth component
	c5, err := cardinal.GetComponent[T5](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T5), err)
	}
	// Get sixth component
	c6, err := cardinal.GetComponent[T6](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T6), err)
	}

	return c1, c2, c3, c4, c5, c6, nil
}

// GetMultipleComponents fetches multiple components (7 in this case) from an entity.
func GetComponents7[T1, T2, T3, T4, T5, T6, T7 types.Component](wCtx cardinal.WorldContext, unitID types.EntityID) (*T1, *T2, *T3, *T4, *T5, *T6, *T7, error) {
	// Get first component
	c1, err := cardinal.GetComponent[T1](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T1), err)
	}
	// Get second component
	c2, err := cardinal.GetComponent[T2](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T2), err)
	}
	// Get third component
	c3, err := cardinal.GetComponent[T3](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T3), err)
	}
	// Get fourth component
	c4, err := cardinal.GetComponent[T4](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T4), err)
	}
	// Get fifth component
	c5, err := cardinal.GetComponent[T5](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T5), err)
	}
	// Get sixth component
	c6, err := cardinal.GetComponent[T6](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T6), err)
	}
	// Get seventh component
	c7, err := cardinal.GetComponent[T7](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T7), err)
	}

	return c1, c2, c3, c4, c5, c6, c7, nil
}

// GetComponents8 fetches multiple components (8 in this case) from an entity.
func GetComponents8[T1, T2, T3, T4, T5, T6, T7, T8 types.Component](wCtx cardinal.WorldContext, unitID types.EntityID) (*T1, *T2, *T3, *T4, *T5, *T6, *T7, *T8, error) {
	// Get first component
	c1, err := cardinal.GetComponent[T1](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T1), err)
	}
	// Get second component
	c2, err := cardinal.GetComponent[T2](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T2), err)
	}
	// Get third component
	c3, err := cardinal.GetComponent[T3](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T3), err)
	}
	// Get fourth component
	c4, err := cardinal.GetComponent[T4](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T4), err)
	}
	// Get fifth component
	c5, err := cardinal.GetComponent[T5](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T5), err)
	}
	// Get sixth component
	c6, err := cardinal.GetComponent[T6](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T6), err)
	}
	// Get seventh component
	c7, err := cardinal.GetComponent[T7](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T7), err)
	}
	// Get eighth component
	c8, err := cardinal.GetComponent[T8](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T8), err)
	}

	return c1, c2, c3, c4, c5, c6, c7, c8, nil
}

// GetComponents9 fetches multiple components (9 in this case) from an entity.
func GetComponents9[T1, T2, T3, T4, T5, T6, T7, T8, T9 types.Component](wCtx cardinal.WorldContext, unitID types.EntityID) (*T1, *T2, *T3, *T4, *T5, *T6, *T7, *T8, *T9, error) {
	// Get first component
	c1, err := cardinal.GetComponent[T1](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T1), err)
	}
	// Get second component
	c2, err := cardinal.GetComponent[T2](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T2), err)
	}
	// Get third component
	c3, err := cardinal.GetComponent[T3](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T3), err)
	}
	// Get fourth component
	c4, err := cardinal.GetComponent[T4](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T4), err)
	}
	// Get fifth component
	c5, err := cardinal.GetComponent[T5](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T5), err)
	}
	// Get sixth component
	c6, err := cardinal.GetComponent[T6](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T6), err)
	}
	// Get seventh component
	c7, err := cardinal.GetComponent[T7](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T7), err)
	}
	// Get eighth component
	c8, err := cardinal.GetComponent[T8](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T8), err)
	}
	// Get ninth component
	c9, err := cardinal.GetComponent[T9](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T9), err)
	}

	return c1, c2, c3, c4, c5, c6, c7, c8, c9, nil
}

// GetComponents10 fetches multiple components (10 in this case) from an entity.
func GetComponents10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 types.Component](wCtx cardinal.WorldContext, unitID types.EntityID) (*T1, *T2, *T3, *T4, *T5, *T6, *T7, *T8, *T9, *T10, error) {
	// Get first component
	c1, err := cardinal.GetComponent[T1](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T1), err)
	}
	// Get second component
	c2, err := cardinal.GetComponent[T2](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T2), err)
	}
	// Get third component
	c3, err := cardinal.GetComponent[T3](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T3), err)
	}
	// Get fourth component
	c4, err := cardinal.GetComponent[T4](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T4), err)
	}
	// Get fifth component
	c5, err := cardinal.GetComponent[T5](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T5), err)
	}
	// Get sixth component
	c6, err := cardinal.GetComponent[T6](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T6), err)
	}
	// Get seventh component
	c7, err := cardinal.GetComponent[T7](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T7), err)
	}
	// Get eighth component
	c8, err := cardinal.GetComponent[T8](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T8), err)
	}
	// Get ninth component
	c9, err := cardinal.GetComponent[T9](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T9), err)
	}
	// Get tenth component
	c10, err := cardinal.GetComponent[T10](wCtx, unitID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("error retrieving component of type %T: %v", *new(T10), err)
	}

	return c1, c2, c3, c4, c5, c6, c7, c8, c9, c10, nil
}
