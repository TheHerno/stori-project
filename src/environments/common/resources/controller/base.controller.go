package controller

import (
	"net/http"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/logger"
	"stori-service/src/utils"

	"github.com/liip/sheriff"
)

//BaseController defines methods to make responses
type BaseController struct{}

/*
MakePaginateResponse Set Data array, pagination headers and Errors to empty array
*/
func (b *BaseController) MakePaginateResponse(collection string, response http.ResponseWriter, data interface{}, statusCode int, pagination *dto.Pagination) {
	parsedData := sheriffParse(collection, data)
	utils.MakePaginateResponse(response, parsedData, statusCode, pagination)
}

/*
MakeSuccessResponse Set Message, Data object and Errors to empty array
*/
func (b *BaseController) MakeSuccessResponse(collection string, response http.ResponseWriter, data interface{}, statusCode int, message string) {
	parsedData := sheriffParse(collection, data)
	utils.MakeSuccessResponse(response, parsedData, statusCode, message)
}

/*
MakeErrorResponse Set Message, Errors to an Array of objects (JSON) and Data to null
*/
func (b *BaseController) MakeErrorResponse(collection string, response http.ResponseWriter, err error) {
	utils.MakeErrorResponse(response, err)
}

/*
sherifParse struct to json filtering by "groups" tag
*/
func sheriffParse(collection string, data interface{}) interface{} {
	option := sheriff.Options{
		Groups: []string{collection}, // Matchs `groups:""` struct tag
	}
	parsedData, err := sheriff.Marshal(&option, data)
	if err != nil {
		logger.GetInstance().Panic(err)
	}
	return parsedData
}
