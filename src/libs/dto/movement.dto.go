package dto

import (
	"stori-service/src/environments/common/resources/entity"
)

/*
MovementList is a DTO to list all Movements of a customer
*/
type MovementList struct {
	Customer  *entity.Customer  `json:"customer"`
	Movements []entity.Movement `json:"movements"`
}
