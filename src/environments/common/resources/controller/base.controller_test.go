package controller

import (
	"net/http"
	"stori-service/src/libs/dto"
	myErrors "stori-service/src/libs/errors"
	"stori-service/src/utils"
	customMock "stori-service/src/utils/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DData struct {
	ID   string `groups:"client"`
	Name string `groups:"admin"`
}

var defaultController = BaseController{}
var defaultCollection string = "client"
var defaultMessage string = "Test"
var defaultErrors []map[string]string = []map[string]string{
	{"test01": "Test01", "test02": "Test02"},
	{"test03": "Test03", "test04": "Test04"},
}
var defaultData = DData{ID: "1", Name: "Name"}

func TestBaseController_MakePaginateResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct header data", func(t *testing.T) {
			// Fixture
			page := &dto.Pagination{
				Page:       1,
				PageSize:   10,
				TotalCount: 20,
			}

			// Run Foo inside request
			response := customMock.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						defaultController.MakePaginateResponse(
							defaultCollection,
							response,
							dto.NewBodyResponse(defaultMessage, defaultErrors, defaultData),
							http.StatusOK,
							page,
						)
					})
				}, "", nil, nil)

			// Assert Data
			assert.Equal(t, "1", response.Header.Get("X-pagination-current-page"))
			assert.Equal(t, "10", response.Header.Get("X-pagination-page-size"))
			assert.Equal(t, "2", response.Header.Get("X-pagination-page-count"))
			assert.Equal(t, "20", response.Header.Get("X-pagination-total-count"))
		})
	})
}

func TestBaseController_MakeSuccessResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct OK response", func(t *testing.T) {
			// Run Foo inside request
			response := customMock.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						defaultController.MakeSuccessResponse(
							defaultCollection,
							response,
							defaultData,
							http.StatusOK,
							defaultData.Name,
						)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := utils.GetBodyResponse(response, &DData{})

			// Assert Data
			assert.Equal(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, defaultData.Name, body.Message)
			assert.Equal(t, []map[string]string{}, body.Errors)
			assert.Equal(t, defaultData.ID, body.Data.(*DData).ID)
		})
	})
}

func TestBaseController_MakeErrorResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct Error response", func(t *testing.T) {
			// Fixture
			expError := []map[string]string{
				{"error": myErrors.ErrNotFound.Error()},
			}

			// Run Foo inside request
			response := customMock.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						defaultController.MakeErrorResponse(
							defaultCollection,
							response,
							myErrors.ErrNotFound,
						)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := utils.GetBodyResponse(response, &DData{})

			// Assert Data
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
			assert.Equal(t, myErrors.ErrNotFound.Error(), body.Message)
			assert.Equal(t, expError, body.Errors)
			assert.Nil(t, body.Data)
		})
	})
}

func TestSheriffParse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Parse struct with 'client' group", func(t *testing.T) {
			data := sheriffParse(defaultCollection, defaultData)

			// Assert Data
			assert.Equal(t, map[string]interface{}{"ID": "1"}, data)
		})
	})
	t.Run("Should panic", func(t *testing.T) {
		t.Run("A struct with invalid since version", func(t *testing.T) {
			type foo struct {
				Bar string `groups:"client" since:"bad_version"`
			}
			data := foo{Bar: "test"}
			assert.Panics(t, func() {
				sheriffParse(defaultCollection, data)
			})
		})
	})
}
