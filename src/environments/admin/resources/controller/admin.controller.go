package controller

import (
	"net/http"
	"stori-service/src/environments/common/resources/controller"
	"stori-service/src/libs/dto"
	"stori-service/src/utils/constant"
)

/*
AdminController composite with BaseController, extends all its public methods by fixing
the first argument (collection) to admin.
*/
type AdminController struct {
	controller.BaseController
}

/*
MakePaginateResponse partial application for base method
*/
func (c *AdminController) MakePaginateResponse(response http.ResponseWriter, data interface{}, statusCode int, pagination *dto.Pagination) {
	c.BaseController.MakePaginateResponse(constant.AdminCollection, response, data, statusCode, pagination)
}

/*
MakeSuccessResponse partial application for base method
*/
func (c *AdminController) MakeSuccessResponse(response http.ResponseWriter, data interface{}, statusCode int, message string) {
	c.BaseController.MakeSuccessResponse(constant.AdminCollection, response, data, statusCode, message)
}

/*
MakeErrorResponse partial application for base method
*/
func (c *AdminController) MakeErrorResponse(response http.ResponseWriter, err error) {
	c.BaseController.MakeErrorResponse(constant.AdminCollection, response, err)
}
