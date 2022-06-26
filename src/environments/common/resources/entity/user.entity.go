package entity

import (
	"stori-service/src/libs/validator"
	"time"

	"gorm.io/gorm"
)

/*
	User model for User table
*/
type User struct {
	UserID    int            `json:"user_id" gorm:"primaryKey" groups:"admin"`
	Name      string         `json:"name" validate:"required,min=3,max=300" groups:"admin"`
	CreatedAt time.Time      `json:"created_at" groups:""`
	UpdatedAt time.Time      `json:"updated_at" groups:""`
	DeletedAt gorm.DeletedAt `json:"deleted_at" groups:""`
}

/*
Validate returns an error if entity doesn't pass any of its own validations
*/
func (user *User) Validate() error {
	if err := validator.ValidateStruct(user); err != nil {
		return err
	}
	return nil
}