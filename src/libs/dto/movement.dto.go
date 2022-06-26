package dto

import "stori-service/src/libs/validator"

/*
NewMovement is a DTO to create a Stock Movement
*/
type NewMovement struct {
	Customerid int `json:"customer_id" validate:"required,gt=0"`
	Quantity   int `json:"quantity" validate:"required,gt=0"`
	Type       int `json:"type" validate:"required,eq=1|eq=-1"`
}

/*
Validate returns an error if entity doesn't pass any of its own validations
*/
func (newMovement *NewMovement) Validate() error {
	if err := validator.ValidateStruct(newMovement); err != nil {
		return err
	}
	return nil
}
