package entity

import (
	"stori-service/src/libs/validator"
	"time"

	"gorm.io/gorm"
)

/*
	Warehouse model for Warehouse table
*/
type Warehouse struct {
	WarehouseID int            `json:"warehouse_id" gorm:"primaryKey" groups:"admin"`
	Name        string         `json:"name" validate:"required,min=3,max=300" groups:"admin"`
	Address     string         `json:"address" validate:"required,min=3,max=300" groups:"admin"`
	UserID      int            `json:"user_id" validate:"required,gt=0" groups:"admin"`
	CreatedAt   time.Time      `json:"created_at" groups:""`
	UpdatedAt   time.Time      `json:"updated_at" groups:""`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" groups:""`
}

/*
Validate returns an error if entity doesn't pass any of its own validations
*/
func (warehouse *Warehouse) Validate() error {
	if err := validator.ValidateStruct(warehouse); err != nil {
		return err
	}
	return nil
}
