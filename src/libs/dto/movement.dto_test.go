package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMovementDTO(t *testing.T) {
	// fixture
	validProductID := 1
	validUserID := 1
	validQty := 10
	validType := 1
	t.Run("Should success on", func(t *testing.T) {
		// fixture
		dto := &NewMovement{
			ProductID: validProductID,
			UserID:    validUserID,
			Quantity:  validQty,
			Type:      validType,
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
				name: "Without ProductID",
				input: &NewMovement{
					UserID:   validUserID,
					Quantity: validQty,
					Type:     validType,
				},
			},
			{
				name: "Invalid ProductID",
				input: &NewMovement{
					ProductID: 0,
					UserID:    validUserID,
					Quantity:  validQty,
					Type:      validType,
				},
			},
			{
				name: "Without UserID",
				input: &NewMovement{
					ProductID: validProductID,
					Quantity:  validQty,
					Type:      validType,
				},
			},
			{
				name: "Invalid UserID",
				input: &NewMovement{
					ProductID: validProductID,
					UserID:    0,
					Quantity:  validQty,
					Type:      validType,
				},
			},
			{
				name: "Without Quantity",
				input: &NewMovement{
					ProductID: validProductID,
					UserID:    validUserID,
					Type:      validType,
				},
			},
			{
				name: "Invalid Quantity",
				input: &NewMovement{
					ProductID: validProductID,
					UserID:    validUserID,
					Quantity:  0,
					Type:      validType,
				},
			},
			{
				name: "Without Type",
				input: &NewMovement{
					ProductID: validProductID,
					UserID:    validUserID,
					Quantity:  validQty,
				},
			},
			{
				name: "Invalid Type",
				input: &NewMovement{
					ProductID: validProductID,
					UserID:    validUserID,
					Quantity:  validQty,
					Type:      0,
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