package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWarehouse(t *testing.T) {
	validId := 5
	validName := "Warehouse del Centro"
	validAddress := "San Luis 301, Mendoza"
	shortString := "El"
	validUserId := 1
	longString := strings.Repeat("E", 301)
	t.Run("Should success on", func(t *testing.T) {
		w := Warehouse{
			WarehouseID: 1,
			Name:        validName,
			Address:     validAddress,
			UserID:      validUserId,
		}
		err := w.Validate()
		assert.NoError(t, err)
	})

	t.Run("Should fail on", func(t *testing.T) {
		testCases := map[string]*Warehouse{
			"Without name": {
				WarehouseID: validId,
				Address:     validAddress,
				UserID:      validUserId,
			},
			"Without adress": {
				WarehouseID: validId,
				Name:        validName,
				UserID:      validUserId,
			},
			"Short adress": {
				WarehouseID: validId,
				Name:        validName,
				UserID:      validUserId,
				Address:     shortString,
			},
			"Short name": {
				WarehouseID: validId,
				Name:        shortString,
				Address:     validAddress,
				UserID:      validUserId,
			},
			"Long adress": {
				WarehouseID: validId,
				Name:        validName,
				UserID:      validUserId,
				Address:     longString,
			},
			"Long name": {
				WarehouseID: validId,
				Name:        longString,
				UserID:      validUserId,
				Address:     validAddress,
			},
			"Invalid UserID": {
				WarehouseID: validId,
				Name:        validName,
				UserID:      0,
				Address:     validAddress,
			},
			"No UserID": {
				WarehouseID: validId,
				Name:        validName,
				Address:     validAddress,
			},
		}

		for name, input := range testCases {
			t.Run(name, func(t *testing.T) {
				err := input.Validate()
				assert.Error(t, err)
			})
		}
	})
}
