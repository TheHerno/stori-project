package dto

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/validator"
)

/*
	UpdateUser is a DTO to update a user
*/
type UpdateUser struct {
	UserID int    `json:"-" validate:"omitempty,gt=0"`
	Name   string `json:"name" validate:"required,min=3,max=300"`
}

/*
	Validate returns an error if DTO doesn't pass any of its validations
*/
func (dto *UpdateUser) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}
	return nil
}

/*
	CreateUser is a DTO to create a user
*/
type CreateUser struct {
	Name   string `json:"name" validate:"required,min=3,max=300"`
	UserID int    `json:"user_id" validate:"required,gt=0"`
}

/*
	Validate returns an error if DTO doesn't pass any of its validations
*/
func (dto *CreateUser) Validate() error {
	//Checks for struct validations (see validate tag)
	if err := validator.ValidateStruct(dto); err != nil {
		return err
	}
	return nil
}

/*
	ParseToUser returns a User entity with the values of the DTO
*/
func (dto *CreateUser) ParseToUser() *entity.User {
	return &entity.User{
		Name:   dto.Name,
		UserID: dto.UserID,
	}
}
