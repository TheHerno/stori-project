package movement

import (
	"net/http"
	"stori-service/src/environments/client/resources/controller"
	"stori-service/src/environments/client/resources/interfaces"
	"stori-service/src/libs/i18n"
	"stori-service/src/utils/helpers"
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

/*
ProcessFile takes the customerID from params and calls the service to process the file
*/
func (c *movementController) ProcessFile(response http.ResponseWriter, request *http.Request) {
	customerID, err := helpers.IDFromRequestToInt(request)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	movementList, err := c.sMovement.ProcessFile(customerID)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}

	c.MakeSuccessResponse(response, movementList, http.StatusOK, i18n.T(i18n.Message{MessageID: "MOVEMENT_LIST.CREATED"}))
}
