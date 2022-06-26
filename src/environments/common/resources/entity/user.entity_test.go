package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	validId := 5
	validName := "Pepe pepito"
	shortString := "El"
	longString := strings.Repeat("E", 301)
	t.Run("Should success on", func(t *testing.T) {
		w := User{
			UserID: 1,
			Name:   validName,
		}
		err := w.Validate()
		assert.NoError(t, err)
	})

	t.Run("Should fail on", func(t *testing.T) {
		testCases := map[string]*User{
			"Without name": {
				UserID: validId,
			},
			"Without adress": {
				UserID: validId,
				Name:   validName,
			},
			"Short adress": {
				UserID: validId,
				Name:   validName,
			},
			"Short name": {
				UserID: validId,
				Name:   shortString,
			},
			"Long name": {
				UserID: validId,
				Name:   longString,
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
