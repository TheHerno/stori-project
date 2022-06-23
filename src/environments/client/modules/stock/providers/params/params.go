package params

import (
	"encoding/json"
	"errors"
	"stori-service/src/libs/adapters/external_call/fetch"
	"stori-service/src/libs/dto"
)

var (
	GetParameter = fetch.GetParameterResponse
)

/*
findParam function gets a param by attr and returns the value
*/
func findParam(paramName string) (interface{}, error) {
	resp, err := GetParameter(paramName)
	if err != nil {
		return nil, err
	}

	var bodyResult dto.BodyResponse
	err = json.NewDecoder(resp.Body).Decode(&bodyResult)
	if err != nil {
		return nil, err
	}

	if len(bodyResult.Errors) != 0 {
		// Me dijeron que no me enrosque así que no haré un parse para los customErrors
		return nil, errors.New("error in response")
	}

	if resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}

	return bodyResult.Data, nil
}

/*
GetParamInt receives a param attribute, calls GetParamFloat64 because
a number param always cames as Float64, then parses to int and returns it
If there is an error, returns it or nil
*/
func GetParamInt(attribute string) (int, error) {
	param, err := findParam(attribute)
	if err != nil {
		return 0.0, err
	}
	return int(param.(float64)), nil
}
