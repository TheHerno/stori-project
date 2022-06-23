package controller

import (
	"net/http"
	"stori-service/src/environments/common/resources/controller"
	"stori-service/src/libs/dto"
	"stori-service/src/utils/constant"
)

/*
ClientController composite with BaseController, extends all its public methods by fixing
the first argument (collection) to client.
*/
type ClientController struct {
	controller.BaseController
}

/*
MakePaginateResponse partial application for base method
*/
func (c *ClientController) MakePaginateResponse(response http.ResponseWriter, data interface{}, statusCode int, pagination *dto.Pagination) {
	c.BaseController.MakePaginateResponse(constant.ClientCollection, response, data, statusCode, pagination)
}

/*
MakeSuccessResponse partial application for base method
*/
func (c *ClientController) MakeSuccessResponse(response http.ResponseWriter, data interface{}, statusCode int, message string) {
	c.BaseController.MakeSuccessResponse(constant.ClientCollection, response, data, statusCode, message)
}

/*
MakeErrorResponse partial application for base method
*/
func (c *ClientController) MakeErrorResponse(response http.ResponseWriter, err error) {
	c.BaseController.MakeErrorResponse(constant.ClientCollection, response, err)
}
