package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomer(t *testing.T) {
	validId := 5
	validName := "Pepe pepito"
	validEmail := "pepepe@hotmail.com"
	shortString := "El"
	longString := strings.Repeat("E", 301)
	t.Run("Should success on", func(t *testing.T) {
		w := Customer{
			CustomerID: 1,
			Name:       validName,
			Email:      validEmail,
		}
		err := w.Validate()
		assert.NoError(t, err)
	})

	t.Run("Should fail on", func(t *testing.T) {
		testCases := map[string]*Customer{
			"Without name": {
				CustomerID: validId,
				Email:      validEmail,
			},
			"Short name": {
				CustomerID: validId,
				Name:       shortString,
				Email:      validEmail,
			},
			"Long name": {
				CustomerID: validId,
				Name:       longString,
				Email:      validEmail,
			},
			"Without email": {
				CustomerID: validId,
				Name:       validName,
			},
			"Invalid email": {
				CustomerID: validId,
				Name:       validName,
				Email:      "invalid email",
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
