package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var trueValue = true

func TestProduct(t *testing.T) {
	validId := 5
	validName := "Product 1"
	validSlug := "product-1"
	validDescription := "Product 1 description"
	shortString := "El"
	longString := strings.Repeat("E", 501)
	t.Run("Should success on", func(t *testing.T) {
		t.Run("With description", func(t *testing.T) {
			p := Product{
				ProductID:   1,
				Name:        validName,
				Slug:        validSlug,
				Enabled:     &trueValue,
				Description: &validDescription,
			}
			err := p.Validate()
			assert.NoError(t, err)
		})
		t.Run("Without description", func(t *testing.T) {
			p := Product{
				ProductID: 1,
				Name:      validName,
				Slug:      validSlug,
				Enabled:   &trueValue,
			}
			err := p.Validate()
			assert.NoError(t, err)
		})
	})

	t.Run("Should fail on", func(t *testing.T) {
		testCases := map[string]*Product{
			"Without name": {
				ProductID:   validId,
				Slug:        validSlug,
				Enabled:     &trueValue,
				Description: &validDescription,
			},
			"Without slug": {
				ProductID:   validId,
				Name:        validName,
				Enabled:     &trueValue,
				Description: &validDescription,
			},
			"Short slug": {
				ProductID: validId,
				Name:      validName,
				Slug:      shortString,
				Enabled:   &trueValue,
			},
			"Short name": {
				ProductID: validId,
				Name:      shortString,
				Slug:      validSlug,
				Enabled:   &trueValue,
			},
			"Long slug": {
				ProductID: validId,
				Name:      validName,
				Slug:      longString,
				Enabled:   &trueValue,
			},
			"Long name": {
				ProductID: validId,
				Name:      longString,
				Slug:      validSlug,
				Enabled:   &trueValue,
			},
			"Long description": {
				ProductID:   validId,
				Name:        validName,
				Description: &longString,
				Slug:        validSlug,
				Enabled:     &trueValue,
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
