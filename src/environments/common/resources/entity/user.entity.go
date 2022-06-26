package entity

import (
	"stori-service/src/libs/validator"
	"time"

	"gorm.io/gorm"
)

/*
	Customer model for Customer table
*/
type Customer struct {
	Customerid int            `json:"customer_id" gorm:"primaryKey" groups:"client"`
	Name       string         `json:"name" validate:"required,min=3,max=100" groups:"client"`
	Email      string         `json:"email" validate:"required,email,max=100" groups:"client"`
	CreatedAt  time.Time      `json:"created_at" groups:""`
	UpdatedAt  time.Time      `json:"updated_at" groups:""`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" groups:""`
}

/*
Validate returns an error if entity doesn't pass any of its own validations
*/
func (customer *Customer) Validate() error {
	if err := validator.ValidateStruct(customer); err != nil {
		return err
	}
	return nil
}
