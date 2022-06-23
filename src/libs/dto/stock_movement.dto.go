package dto

import "stori-service/src/libs/validator"

/*
NewStockMovement is a DTO to create a Stock Movement
*/
type NewStockMovement struct {
	ProductID   int    `json:"product_id" validate:"required,gt=0"`
	UserID      int    `json:"-" validate:"required,gt=0"`
	WarehouseID int    `json:"warehouse_id" validate:"omitempty,gt=0"`
	Quantity    int    `json:"quantity" validate:"required,gt=0"`
	Concept     string `json:"concept" validate:"required,max=100"`
	Type        int    `json:"type" validate:"required,eq=1|eq=-1"`
}

/*
Validate returns an error if entity doesn't pass any of its own validations
*/
func (newStockMovement *NewStockMovement) Validate() error {
	if err := validator.ValidateStruct(newStockMovement); err != nil {
		return err
	}
	return nil
}
