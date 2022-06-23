package dto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStockMovementDTO(t *testing.T) {
	// fixture
	validProductID := 1
	validUserID := 1
	validQty := 10
	validConcept := "Venta"
	validType := 1
	longString := strings.Repeat("a", 101)
	t.Run("Should success on", func(t *testing.T) {
		// fixture
		dto := &NewStockMovement{
			ProductID: validProductID,
			UserID:    validUserID,
			Quantity:  validQty,
			Concept:   validConcept,
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
			input *NewStockMovement
		}{
			{
				name: "Without ProductID",
				input: &NewStockMovement{
					UserID:   validUserID,
					Quantity: validQty,
					Concept:  validConcept,
					Type:     validType,
				},
			},
			{
				name: "Invalid ProductID",
				input: &NewStockMovement{
					ProductID: 0,
					UserID:    validUserID,
					Quantity:  validQty,
					Concept:   validConcept,
					Type:      validType,
				},
			},
			{
				name: "Without UserID",
				input: &NewStockMovement{
					ProductID: validProductID,
					Quantity:  validQty,
					Concept:   validConcept,
					Type:      validType,
				},
			},
			{
				name: "Invalid UserID",
				input: &NewStockMovement{
					ProductID: validProductID,
					UserID:    0,
					Quantity:  validQty,
					Concept:   validConcept,
					Type:      validType,
				},
			},
			{
				name: "Without Quantity",
				input: &NewStockMovement{
					ProductID: validProductID,
					UserID:    validUserID,
					Concept:   validConcept,
					Type:      validType,
				},
			},
			{
				name: "Invalid Quantity",
				input: &NewStockMovement{
					ProductID: validProductID,
					UserID:    validUserID,
					Quantity:  0,
					Concept:   validConcept,
					Type:      validType,
				},
			},
			{
				name: "Without Concept",
				input: &NewStockMovement{
					ProductID: validProductID,
					UserID:    validUserID,
					Quantity:  validQty,
					Type:      validType,
				},
			},
			{
				name: "Invalid Concept",
				input: &NewStockMovement{
					ProductID: validProductID,
					UserID:    validUserID,
					Quantity:  validQty,
					Concept:   longString,
					Type:      validType,
				},
			},
			{
				name: "Without Type",
				input: &NewStockMovement{
					ProductID: validProductID,
					UserID:    validUserID,
					Quantity:  validQty,
					Concept:   validConcept,
				},
			},
			{
				name: "Invalid Type",
				input: &NewStockMovement{
					ProductID: validProductID,
					UserID:    validUserID,
					Quantity:  validQty,
					Concept:   validConcept,
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
