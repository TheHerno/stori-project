package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMovementDTO(t *testing.T) {
	// fixture
	validCustomerID := 1
	validQty := 10
	validType := 1
	t.Run("Should success on", func(t *testing.T) {
		// fixture
		dto := &NewMovement{
			CustomerID: validCustomerID,
			Quantity:   validQty,
			Type:       validType,
		}
		// action
		err := dto.Validate()
		// assertion
		assert.NoError(t, err)
	})
	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			name  string
			input *NewMovement
		}{
			{
				name: "Without CustomerID",
				input: &NewMovement{
					Quantity: validQty,
					Type:     validType,
				},
			},
			{
				name: "Invalid CustomerID",
				input: &NewMovement{
					CustomerID: 0,
					Quantity:   validQty,
					Type:       validType,
				},
			},
			{
				name: "Without Quantity",
				input: &NewMovement{
					CustomerID: validCustomerID,
					Type:       validType,
				},
			},
			{
				name: "Invalid Quantity",
				input: &NewMovement{
					CustomerID: validCustomerID,
					Quantity:   0,
					Type:       validType,
				},
			},
			{
				name: "Without Type",
				input: &NewMovement{
					CustomerID: validCustomerID,
					Quantity:   validQty,
				},
			},
			{
				name: "Invalid Type",
				input: &NewMovement{
					CustomerID: validCustomerID,
					Quantity:   validQty,
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
