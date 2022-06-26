package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovement(t *testing.T) {
	// fixture
	validMovementID := 1
	validUserID := 1
	validQty := 10
	validAvailable := 20
	validType := 1
	t.Run("Should success on", func(t *testing.T) {
		// fixture
		movement := &Movement{
			MovementID: validMovementID,
			UserID:     validUserID,
			Quantity:   validQty,
			Available:  validAvailable,
			Type:       validType,
		}
		// action
		err := movement.Validate()
		// assertion
		assert.NoError(t, err)
	})

	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			name  string
			input *Movement
		}{
			{
				name: "Without UserID",
				input: &Movement{
					MovementID: validMovementID,
					Quantity:   validQty,
					Available:  validAvailable,
					Type:       validType,
				},
			},
			{
				name: "Invalid UserID",
				input: &Movement{
					MovementID: validMovementID,
					UserID:     0,
					Quantity:   validQty,
					Available:  validAvailable,
					Type:       validType,
				},
			},
			{
				name: "Without Quantity",
				input: &Movement{
					MovementID: validMovementID,
					UserID:     validUserID,
					Available:  validAvailable,
					Type:       validType,
				},
			},
			{
				name: "Invalid Quantity",
				input: &Movement{
					MovementID: validMovementID,
					UserID:     validUserID,
					Quantity:   0,
					Available:  validAvailable,
					Type:       validType,
				},
			},
			{
				name: "Without Available",
				input: &Movement{
					MovementID: validMovementID,
					UserID:     validUserID,
					Quantity:   validQty,
					Type:       validType,
				},
			},

			{
				name: "Invalid Available",
				input: &Movement{
					MovementID: validMovementID,
					UserID:     validUserID,
					Quantity:   validQty,
					Available:  -1,
					Type:       validType,
				},
			},
			{
				name: "Without Type",
				input: &Movement{
					MovementID: validMovementID,
					UserID:     validUserID,
					Quantity:   validQty,
					Available:  validAvailable,
				},
			},
			{
				name: "Invalid Type",
				input: &Movement{
					MovementID: validMovementID,
					UserID:     validUserID,
					Quantity:   validQty,
					Available:  validAvailable,
					Type:       0,
				},
			},
		}

		for _, tC := range testCases {
			t.Run(tC.name, func(t *testing.T) {
				// action
				err := tC.input.Validate()
				// assertion
				assert.Error(t, err)
			})
		}
	})
}
