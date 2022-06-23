package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStockMovement(t *testing.T) {
	// fixture
	validStockMovementID := 1
	validProductID := 1
	validWarehouseID := 1
	validQty := 10
	validAvailable := 20
	validConcept := "Venta"
	validType := 1
	longString := strings.Repeat("a", 101)
	t.Run("Should success on", func(t *testing.T) {
		// fixture
		stockMovement := &StockMovement{
			StockMovementID: validStockMovementID,
			ProductID:       validProductID,
			WarehouseID:     validWarehouseID,
			Quantity:        validQty,
			Available:       validAvailable,
			Concept:         validConcept,
			Type:            validType,
		}
		// action
		err := stockMovement.Validate()
		// assertion
		assert.NoError(t, err)
	})

	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			name  string
			input *StockMovement
		}{
			{
				name: "Without ProductID",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					WarehouseID:     validWarehouseID,
					Quantity:        validQty,
					Available:       validAvailable,
					Concept:         validConcept,
					Type:            validType,
				},
			},
			{
				name: "Invalid ProductID",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					ProductID:       0,
					WarehouseID:     validWarehouseID,
					Quantity:        validQty,
					Available:       validAvailable,
					Concept:         validConcept,
					Type:            validType,
				},
			},
			{
				name: "Without WarehouseID",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					ProductID:       validProductID,
					Quantity:        validQty,
					Available:       validAvailable,
					Concept:         validConcept,
					Type:            validType,
				},
			},
			{
				name: "Invalid WarehouseID",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					ProductID:       validProductID,
					WarehouseID:     0,
					Quantity:        validQty,
					Available:       validAvailable,
					Concept:         validConcept,
					Type:            validType,
				},
			},
			{
				name: "Without Quantity",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					ProductID:       validProductID,
					WarehouseID:     validWarehouseID,
					Available:       validAvailable,
					Concept:         validConcept,
					Type:            validType,
				},
			},
			{
				name: "Invalid Quantity",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					ProductID:       validProductID,
					WarehouseID:     validWarehouseID,
					Quantity:        0,
					Available:       validAvailable,
					Concept:         validConcept,
					Type:            validType,
				},
			},
			{
				name: "Without Available",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					ProductID:       validProductID,
					WarehouseID:     validWarehouseID,
					Quantity:        validQty,
					Concept:         validConcept,
					Type:            validType,
				},
			},

			{
				name: "Invalid Available",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					ProductID:       validProductID,
					WarehouseID:     validWarehouseID,
					Quantity:        validQty,
					Available:       -1,
					Concept:         validConcept,
					Type:            validType,
				},
			},
			{
				name: "Without Concept",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					ProductID:       validProductID,
					WarehouseID:     validWarehouseID,
					Quantity:        validQty,
					Available:       validAvailable,
					Type:            validType,
				},
			},
			{
				name: "Invalid Concept",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					ProductID:       validProductID,
					WarehouseID:     validWarehouseID,
					Quantity:        validQty,
					Available:       validAvailable,
					Concept:         longString,
					Type:            validType,
				},
			},
			{
				name: "Without Type",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					ProductID:       validProductID,
					WarehouseID:     validWarehouseID,
					Quantity:        validQty,
					Available:       validAvailable,
					Concept:         validConcept,
				},
			},
			{
				name: "Invalid Type",
				input: &StockMovement{
					StockMovementID: validStockMovementID,
					ProductID:       validProductID,
					WarehouseID:     validWarehouseID,
					Quantity:        validQty,
					Available:       validAvailable,
					Concept:         validConcept,
					Type:            0,
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
