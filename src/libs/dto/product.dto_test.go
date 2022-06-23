package dto

import (
	"stori-service/src/utils/helpers"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var trueValue = true

func TestUpdateProduct(t *testing.T) {
	validName := "El Product válido"
	validDescription := helpers.PointerToString("Descripción de producto válida")
	shortString := "él"
	longString := strings.Repeat("a", 501)
	t.Run("Should success on", func(t *testing.T) {
		testCases := []struct {
			name string
			dto  *UpdateProduct
		}{
			{
				name: "With description",
				dto: &UpdateProduct{
					Name:        validName,
					Description: validDescription,
					Enabled:     &trueValue,
				},
			},
			{
				name: "Without description",
				dto: &UpdateProduct{
					Name:    validName,
					Enabled: &trueValue,
				},
			},
		}
		for _, tC := range testCases {
			t.Run(tC.name, func(t *testing.T) {
				err := tC.dto.Validate()
				assert.NoError(t, err)
			})
		}
	})
	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			TestName string
			Input    *UpdateProduct
		}{
			{
				TestName: "Without name",
				Input: &UpdateProduct{
					Description: validDescription,
					Enabled:     &trueValue,
				},
			},
			{
				TestName: "Without enabled",
				Input: &UpdateProduct{
					Name: validName,
				},
			},
			{
				TestName: "Short name",
				Input: &UpdateProduct{
					Name:        shortString,
					Description: validDescription,
				},
			},
			{
				TestName: "Long name",
				Input: &UpdateProduct{
					Name:        longString,
					Description: validDescription,
					Enabled:     &trueValue,
				},
			},
			{
				TestName: "Long description",
				Input: &UpdateProduct{
					Name:        validName,
					Description: &longString,
					Enabled:     &trueValue,
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

func TestCreateProduct(t *testing.T) {
	validName := "El Product válido"
	validDescription := helpers.PointerToString("La descripción válida")
	shortString := "él"
	longString := strings.Repeat("a", 501)
	t.Run("Should success on", func(t *testing.T) {
		testCases := []struct {
			name string
			dto  *CreateProduct
		}{
			{
				name: "With description",
				dto: &CreateProduct{
					Name:        validName,
					Description: validDescription,
					Enabled:     &trueValue,
				},
			},
			{
				name: "Without description",
				dto: &CreateProduct{
					Name:    validName,
					Enabled: &trueValue,
				},
			},
		}
		for _, tC := range testCases {
			t.Run(tC.name, func(t *testing.T) {
				err := tC.dto.Validate()
				assert.NoError(t, err)
			})
		}
	})
	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			TestName string
			Input    *CreateProduct
		}{
			{
				TestName: "Without name",
				Input: &CreateProduct{
					Description: validDescription,
					Enabled:     &trueValue,
				},
			},
			{
				TestName: "Short name",
				Input: &CreateProduct{
					Name:        shortString,
					Description: validDescription,
					Enabled:     &trueValue,
				},
			},
			{
				TestName: "Long Description",
				Input: &CreateProduct{
					Name:        validName,
					Description: &longString,
					Enabled:     &trueValue,
				},
			},
			{
				TestName: "Long name",
				Input: &CreateProduct{
					Name:        longString,
					Description: validDescription,
					Enabled:     &trueValue,
				},
			},
			{
				TestName: "Without Enabled",
				Input: &CreateProduct{
					Name:        validName,
					Description: validDescription,
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

func TestParseToProduct(t *testing.T) {
	validName := "Product Válido"
	validDescription := helpers.PointerToString("Descripción válida")
	t.Run("Should success on", func(t *testing.T) {
		dto := &CreateProduct{
			Name:        validName,
			Description: validDescription,
			Enabled:     &trueValue,
		}
		result := dto.ParseToProduct()
		assert.Equal(t, result.Name, validName)
		assert.Equal(t, result.Enabled, &trueValue)
		assert.Equal(t, result.Description, validDescription)
	})
}
