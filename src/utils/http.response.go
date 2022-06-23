package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"stori-service/src/libs/dto"
	myErrors "stori-service/src/libs/errors"
	"stori-service/src/utils/constant"
	"strconv"
)

/*
MakePaginateResponse Set Data array, pagination headers and Errors to empty array
*/
func MakePaginateResponse(response http.ResponseWriter, data interface{}, statusCode int, pagination *dto.Pagination) {
	response.Header().Set("X-pagination-total-count", strconv.FormatInt(pagination.TotalCount, 10))
	response.Header().Set("X-pagination-page-count", strconv.Itoa(pagination.PageCount()))
	response.Header().Set("X-pagination-current-page", strconv.Itoa(pagination.Page))
	response.Header().Set("X-pagination-page-size", strconv.Itoa(pagination.PageSize))
	body := dto.NewBodyResponse("Success", make([]map[string]string, 0), data)
	makeResponse(response, body, statusCode)
}

/*
MakeSuccessResponse Set Message, Data object and Errors to empty array
*/
func MakeSuccessResponse(response http.ResponseWriter, data interface{}, statusCode int, message string) {
	body := dto.NewBodyResponse(message, make([]map[string]string, 0), data)
	makeResponse(response, body, statusCode)
}

/*
MakeErrorResponse Set Message, Errors to an Array of objects (JSON) and Data to null
*/
func MakeErrorResponse(response http.ResponseWriter, err error) {
	errorMessage := myErrors.GetErrorMessage(err)
	errors := []map[string]string{{"error": errorMessage}}
	body := dto.NewBodyResponse(errorMessage, errors, nil)
	SetActionNeeded(err, response)
	makeResponse(response, body, myErrors.GetStatusCode(err))
}

/*
makeResponse Serialize and send the JSON body to client. Above methods end here
*/
func makeResponse(response http.ResponseWriter, body *dto.BodyResponse, statusCode int) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	json.NewEncoder(response).Encode(&body)
}

/*
GetBodyRequest parses request body to received variable
*/
func GetBodyRequest(req *http.Request, data interface{}) error {
	bodyBytes, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	err := json.Unmarshal(bodyBytes, data)
	if err != nil {
		return err
	}
	return nil
}

/*
GetBodyResponse parses body to received variable
*/
func GetBodyResponse(res *http.Response, data interface{}) (*dto.BodyResponse, error) {
	bodyResponse := dto.NewBodyResponse("", nil, data)
	err := takeBodyResponse(res, bodyResponse)
	if err != nil {
		return nil, err
	}
	return bodyResponse, nil
}

/*
takeBodyResponse parses response body to received variable and closes the body
*/
func takeBodyResponse(res *http.Response, data interface{}) error {
	body := res.Body
	defer body.Close()
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return err
	}
	return nil
}

/*
SetActionNeeded receives error and response, then calls GetAction from errors package
finally it sets on header the action needed
*/
func SetActionNeeded(err error, response http.ResponseWriter) {
	action := myErrors.GetAction(err)
	if action != nil {
		response.Header().Set(constant.HeaderNeedsAction, *action)
	}
}
