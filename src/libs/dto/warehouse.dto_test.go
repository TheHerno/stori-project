package dto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateWarehouse(t *testing.T) {
	validName := "El Warehouse válido"
	validAddress := "Calle falsa 123"
	shortString := "él"
	longString := strings.Repeat("a", 301)
	t.Run("Should success on", func(t *testing.T) {
		dto := UpdateWarehouse{
			Name:    validName,
			Address: validAddress,
		}
		err := dto.Validate()
		assert.NoError(t, err)
	})
	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			TestName string
			Input    *UpdateWarehouse
		}{
			{
				TestName: "Without name",
				Input: &UpdateWarehouse{
					Address: validAddress,
				},
			},
			{
				TestName: "Without address",
				Input: &UpdateWarehouse{
					Name: validName,
				},
			},
			{
				TestName: "Short name",
				Input: &UpdateWarehouse{
					Name:    shortString,
					Address: validAddress,
				},
			},
			{
				TestName: "Short Address",
				Input: &UpdateWarehouse{
					Name:    validName,
					Address: shortString,
				},
			},
			{
				TestName: "Long name",
				Input: &UpdateWarehouse{
					Name:    longString,
					Address: validAddress,
				},
			},
			{
				TestName: "Long address",
				Input: &UpdateWarehouse{
					Name:    validName,
					Address: longString,
				},
			},
		}
		for _, tC := range testCases {
			t.Run(tC.TestName, func(t *testing.T) {
				err := tC.Input.Validate()
				assert.Error(t, err)
			})
		}
	})
}

func TestCreateWarehouse(t *testing.T) {
	validName := "El Warehouse válido"
	validAddress := "Calle falsa 123"
	shortString := "él"
	validUserID := 1
	longString := strings.Repeat("a", 301)
	t.Run("Should success on", func(t *testing.T) {
		dto := CreateWarehouse{
			Name:    validName,
			Address: validAddress,
			UserID:  validUserID,
		}
		err := dto.Validate()
		assert.NoError(t, err)
	})
	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			TestName string
			Input    *CreateWarehouse
		}{
			{
				TestName: "Without name",
				Input: &CreateWarehouse{
					Address: validAddress,
					UserID:  validUserID,
				},
			},
			{
				TestName: "Without UserID",
				Input: &CreateWarehouse{
					Address: validAddress,
					Name:    validName,
				},
			},
			{
				TestName: "Invalid userID",
				Input: &CreateWarehouse{
					Name:    validName,
					Address: validAddress,
					UserID:  0,
				},
			},
			{
				TestName: "Without address",
				Input: &CreateWarehouse{
					Name:   validName,
					UserID: validUserID,
				},
			},
			{
				TestName: "Short name",
				Input: &CreateWarehouse{
					Name:    shortString,
					Address: validAddress,
					UserID:  validUserID,
				},
			},
			{
				TestName: "Short Address",
				Input: &CreateWarehouse{
					Name:    validName,
					Address: shortString,
					UserID:  validUserID,
				},
			},
			{
				TestName: "Long name",
				Input: &CreateWarehouse{
					Name:    longString,
					Address: validAddress,
					UserID:  validUserID,
				},
			},
			{
				TestName: "Long address",
				Input: &CreateWarehouse{
					Name:    validName,
					Address: longString,
					UserID:  validUserID,
				},
			},
		}
		for _, tC := range testCases {
			t.Run(tC.TestName, func(t *testing.T) {
				err := tC.Input.Validate()
				assert.Error(t, err)
			})
		}
	})
}

func TestParseToWarehouse(t *testing.T) {
	validName := "Warehouse Válido"
	validAddress := "Calle falsa 123"
	t.Run("Should success on", func(t *testing.T) {
		dto := &CreateWarehouse{
			Name:    validName,
			Address: validAddress,
		}
		result := dto.ParseToWarehouse()
		assert.Equal(t, result.Name, validName)
		assert.Equal(t, result.Address, validAddress)
	})
}
