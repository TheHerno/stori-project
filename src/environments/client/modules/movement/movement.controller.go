package stock

import (
	"stori-service/src/environments/client/resources/controller"
	"stori-service/src/environments/client/resources/interfaces"
)

// struct that implements IMovementController
type movementController struct {
	controller.ClientController
	sMovement interfaces.IMovementService
}

/*
NewMovementController creates a new controller, receives service by dependency injection
and returns IMovementController, so needs to implement all its methods
*/
func NewMovementController(sMovement interfaces.IMovementService) interfaces.IMovementController {
	return &movementController{sMovement: sMovement}
}
