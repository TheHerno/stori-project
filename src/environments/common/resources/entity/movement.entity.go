package entity

import (
	"stori-service/src/libs/validator"
	"time"

	"gorm.io/gorm"
)

/*
Movement model for movement table
*/
type Movement struct {
	MovementID int            `json:"movement_id" gorm:"primaryKey" groups:""`
	CustomerID int            `json:"customer_id" groups:"" validate:"required,gte=1"`
	Quantity   float64        `json:"quantity" groups:"client" validate:"required,gt=0"`
	Available  float64        `json:"available" groups:"client" validate:"required,gte=0"`
	Type       int            `json:"type" groups:"client" validate:"required,eq=1|eq=-1"`
	Date       time.Time      `json:"date" groups:"client" validate:"required"`
	CreatedAt  time.Time      `json:"created_at" groups:""`
	UpdatedAt  time.Time      `json:"updated_at" groups:""`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" groups:""`
}

/*
Validate returns an error if entity doesn't pass any of its own validations
*/
func (movement *Movement) Validate() error {
	if err := validator.ValidateStruct(movement); err != nil {
		return err
	}
	return nil
}
