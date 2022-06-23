package entity

import (
	"stori-service/src/libs/validator"
	"time"

	"gorm.io/gorm"
)

/*
	Product model for Product table
*/
type Product struct {
	ProductID   int            `json:"product_id" gorm:"primaryKey" groups:"admin"`
	Name        string         `json:"name" validate:"required,min=3,max=300" groups:"admin"`
	Slug        string         `json:"slug" validate:"required,min=8,max=305" groups:"admin"`
	Enabled     *bool          `json:"enabled" validate:"required" groups:"admin"`
	Description *string        `json:"description" validate:"omitempty,min=3,max=500" groups:"admin"`
	CreatedAt   time.Time      `json:"created_at" groups:""`
	UpdatedAt   time.Time      `json:"updated_at" groups:""`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" groups:""`
}

/*
Validate returns an error if entity doesn't pass any of its own validations
*/
func (product *Product) Validate() error {
	if err := validator.ValidateStruct(product); err != nil {
		return err
	}
	return nil
}
