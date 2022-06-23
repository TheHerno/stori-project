package dto

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/validator"
)

/*
	UpdateWarehouse is a DTO to update a warehouse
*/
type UpdateWarehouse struct {
	WarehouseID int    `json:"-" validate:"omitempty,gt=0"`
	Name        string `json:"name" validate:"required,min=3,max=300"`
	Address     string `json:"address" validate:"required,min=3,max=300"`
}

/*
	Validate returns an error if DTO doesn't pass any of its validations
*/
func (dto *UpdateWarehouse) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}
	return nil
}

/*
	CreateWarehouse is a DTO to create a warehouse
*/
type CreateWarehouse struct {
	Name    string `json:"name" validate:"required,min=3,max=300"`
	Address string `json:"address" validate:"required,min=3,max=300"`
	UserID  int    `json:"user_id" validate:"required,gt=0"`
}

/*
	Validate returns an error if DTO doesn't pass any of its validations
*/
func (dto *CreateWarehouse) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}
	return nil
}

/*
	ParseToWarehouse returns a Warehouse entity with the values of the DTO
*/
func (dto *CreateWarehouse) ParseToWarehouse() *entity.Warehouse {
	return &entity.Warehouse{
		Name:    dto.Name,
		Address: dto.Address,
		UserID:  dto.UserID,
	}
}
