package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMovement(t *testing.T) {
	// fixture
	validMovementID := 1
	validCustomerID := 1
	validQty := 10.0
	validDate := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	validAvailable := 20.00
	validType := 1
	t.Run("Should success on", func(t *testing.T) {
		// fixture
		movement := &Movement{
			MovementID: validMovementID,
			CustomerID: validCustomerID,
			Quantity:   validQty,
			Available:  validAvailable,
			Type:       validType,
			Date:       validDate,
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
				name: "Without CustomerID",
				input: &Movement{
					MovementID: validMovementID,
					Quantity:   validQty,
					Available:  validAvailable,
					Type:       validType,
					Date:       validDate,
				},
			},
			{
				name: "Invalid CustomerID",
				input: &Movement{
					MovementID: validMovementID,
					CustomerID: 0,
					Quantity:   validQty,
					Available:  validAvailable,
					Type:       validType,
					Date:       validDate,
				},
			},
			{
				name: "Without Quantity",
				input: &Movement{
					MovementID: validMovementID,
					CustomerID: validCustomerID,
					Available:  validAvailable,
					Type:       validType,
					Date:       validDate,
				},
			},
			{
				name: "Invalid Quantity",
				input: &Movement{
					MovementID: validMovementID,
					CustomerID: validCustomerID,
					Quantity:   0,
					Available:  validAvailable,
					Type:       validType,
					Date:       validDate,
				},
			},
			{
				name: "Without Available",
				input: &Movement{
					MovementID: validMovementID,
					CustomerID: validCustomerID,
					Quantity:   validQty,
					Type:       validType,
					Date:       validDate,
				},
			},

			{
				name: "Invalid Available",
				input: &Movement{
					MovementID: validMovementID,
					CustomerID: validCustomerID,
					Quantity:   validQty,
					Available:  -1,
					Type:       validType,
					Date:       validDate,
				},
			},
			{
				name: "Without Type",
				input: &Movement{
					MovementID: validMovementID,
					CustomerID: validCustomerID,
					Quantity:   validQty,
					Available:  validAvailable,
					Date:       validDate,
				},
			},
			{
				name: "Invalid Type",
				input: &Movement{
					MovementID: validMovementID,
					CustomerID: validCustomerID,
					Quantity:   validQty,
					Available:  validAvailable,
					Type:       0,
					Date:       validDate,
				},
			},
			{
				name: "Without date",
				input: &Movement{
					MovementID: validMovementID,
					CustomerID: validCustomerID,
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
