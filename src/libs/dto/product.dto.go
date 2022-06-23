package dto

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/validator"
)

/*
	ProductWithStock is a DTO that contains a Product and its stock
*/
type ProductWithStock struct {
	ProductID   int     `json:"-" groups:"client"`
	Name        string  `json:"name" groups:"client"`
	Description *string `json:"description" groups:"client"`
	Slug        string  `json:"slug" groups:"client"`
	Stock       int     `json:"stock" groups:"client"`
}

/*
	UpdateProduct is a DTO to update a product
*/
type UpdateProduct struct {
	ProductID   int     `json:"-" validate:"omitempty,gt=0"`
	Name        string  `json:"name" validate:"required,min=3,max=300"`
	Description *string `json:"description" validate:"omitempty,min=3,max=500"`
	Enabled     *bool   `json:"enabled" validate:"required"`
}

/*
	Validate returns an error if DTO doesn't pass any of its validations
*/
func (dto *UpdateProduct) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}
	return nil
}

/*
	CreateProduct is a DTO to create a product
*/
type CreateProduct struct {
	Name        string  `json:"name" validate:"required,min=3,max=300"`
	Description *string `json:"description" validate:"omitempty,min=3,max=500"`
	Enabled     *bool   `json:"enabled" validate:"required"`
}

/*
	Validate returns an error if DTO doesn't pass any of its validations
*/
func (dto *CreateProduct) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}
	return nil
}

/*
	ParseToProduct returns a Product entity with the values of the DTO
*/
func (dto *CreateProduct) ParseToProduct() *entity.Product {
	return &entity.Product{
		Name:        dto.Name,
		Description: dto.Description,
		Enabled:     dto.Enabled,
	}
}
