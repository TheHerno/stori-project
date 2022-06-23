package entity

import (
	"stori-service/src/libs/validator"
	"time"

	"gorm.io/gorm"
)

/*
StockMovement model for stock_movement table
*/
type StockMovement struct {
	StockMovementID int            `json:"stock_movement_id" gorm:"primaryKey" groups:"client"`
	ProductID       int            `json:"product_id" groups:"client" validate:"required,gt=0"`
	WarehouseID     int            `json:"warehouse_id" groups:"client" validate:"required,gt=0"`
	Quantity        int            `json:"quantity" groups:"client" validate:"required,gt=0"`
	Available       int            `json:"available" groups:"client" validate:"required,gte=0"`
	Concept         string         `json:"concept" groups:"client" validate:"required,max=100"`
	Type            int            `json:"type" groups:"client" validate:"required,eq=1|eq=-1"`
	CreatedAt       time.Time      `json:"created_at" groups:""`
	UpdatedAt       time.Time      `json:"updated_at" groups:""`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" groups:""`
}

/*
Validate returns an error if entity doesn't pass any of its own validations
*/
func (stockMovement *StockMovement) Validate() error {
	if err := validator.ValidateStruct(stockMovement); err != nil {
		return err
	}
	return nil
}
