package stock

import (
	"net/http"
	"stori-service/src/environments/client/resources/controller"
	"stori-service/src/environments/client/resources/interfaces"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/i18n"
	"stori-service/src/libs/middleware"
	"stori-service/src/utils"
)

// struct that implements IStockMovementController
type stockMovementController struct {
	controller.ClientController
	sStockMovement interfaces.IStockMovementService
}

/*
NewStockMovementController creates a new controller, receives service by dependency injection
and returns IStockMovementController, so needs to implement all its methods
*/
func NewStockMovementController(sStockMovement interfaces.IStockMovementService) interfaces.IStockMovementController {
	return &stockMovementController{sStockMovement: sStockMovement}
}

/*
Create gets the body, then calls Create service
*/
func (c *stockMovementController) Create(response http.ResponseWriter, request *http.Request) {
	stockmovementToCreate := &dto.NewStockMovement{}
	if err := utils.GetBodyRequest(request, stockmovementToCreate); err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	user := request.Context().Value(middleware.ContextKeyUser).(*middleware.User)
	stockmovementToCreate.UserID = user.UserID
	stockmovementCreated, err := c.sStockMovement.Create(stockmovementToCreate)
	if err != nil {
		c.MakeErrorResponse(response, err)
		return
	}
	c.MakeSuccessResponse(response, stockmovementCreated, http.StatusCreated, i18n.T(i18n.Message{MessageID: "STOCKMOVEMENT.CREATED"}))
}
