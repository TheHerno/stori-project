package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMovementDTO(t *testing.T) {
	// fixture
	validCustomerid := 1
	validQty := 10
	validType := 1
	t.Run("Should success on", func(t *testing.T) {
		// fixture
		dto := &NewMovement{
			Customerid: validCustomerid,
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
				name: "Without Customerid",
				input: &NewMovement{
					Quantity: validQty,
					Type:     validType,
				},
			},
			{
				name: "Invalid Customerid",
				input: &NewMovement{
					Customerid: 0,
					Quantity:   validQty,
					Type:       validType,
				},
			},
			{
				name: "Without Quantity",
				input: &NewMovement{
					Customerid: validCustomerid,
					Type:       validType,
				},
			},
			{
				name: "Invalid Quantity",
				input: &NewMovement{
					Customerid: validCustomerid,
					Quantity:   0,
					Type:       validType,
				},
			},
			{
				name: "Without Type",
				input: &NewMovement{
					Customerid: validCustomerid,
					Quantity:   validQty,
				},
			},
			{
				name: "Invalid Type",
				input: &NewMovement{
					Customerid: validCustomerid,
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
