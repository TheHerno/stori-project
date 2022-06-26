package dto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	validName := "El User válido"
	shortString := "él"
	longString := strings.Repeat("a", 301)
	t.Run("Should success on", func(t *testing.T) {
		dto := UpdateUser{
			Name: validName,
		}
		err := dto.Validate()
		assert.NoError(t, err)
	})
	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			TestName string
			Input    *UpdateUser
		}{
			{
				TestName: "Without name",
				Input:    &UpdateUser{},
			},
			{
				TestName: "Without address",
				Input: &UpdateUser{
					Name: validName,
				},
			},
			{
				TestName: "Short name",
				Input: &UpdateUser{
					Name: shortString,
				},
			},
			{
				TestName: "Short Address",
				Input: &UpdateUser{
					Name: validName,
				},
			},
			{
				TestName: "Long name",
				Input: &UpdateUser{
					Name: longString,
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

func TestCreateUser(t *testing.T) {
	validName := "El User válido"
	shortString := "él"
	validUserID := 1
	longString := strings.Repeat("a", 301)
	t.Run("Should success on", func(t *testing.T) {
		dto := CreateUser{
			Name:   validName,
			UserID: validUserID,
		}
		err := dto.Validate()
		assert.NoError(t, err)
	})
	t.Run("Should fail on", func(t *testing.T) {
		testCases := []struct {
			TestName string
			Input    *CreateUser
		}{
			{
				TestName: "Without name",
				Input: &CreateUser{
					UserID: validUserID,
				},
			},
			{
				TestName: "Without UserID",
				Input: &CreateUser{
					Name: validName,
				},
			},
			{
				TestName: "Invalid userID",
				Input: &CreateUser{
					Name:   validName,
					UserID: 0,
				},
			},
			{
				TestName: "Without address",
				Input: &CreateUser{
					Name:   validName,
					UserID: validUserID,
				},
			},
			{
				TestName: "Short name",
				Input: &CreateUser{
					Name:   shortString,
					UserID: validUserID,
				},
			},
			{
				TestName: "Short Address",
				Input: &CreateUser{
					Name:   validName,
					UserID: validUserID,
				},
			},
			{
				TestName: "Long name",
				Input: &CreateUser{
					Name:   longString,
					UserID: validUserID,
				},
			},
			{
				TestName: "Long address",
				Input: &CreateUser{
					Name:   validName,
					UserID: validUserID,
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

func TestParseToUser(t *testing.T) {
	validName := "User Válido"
	t.Run("Should success on", func(t *testing.T) {
		dto := &CreateUser{
			Name: validName,
		}
		result := dto.ParseToUser()
		assert.Equal(t, result.Name, validName)
	})
}
