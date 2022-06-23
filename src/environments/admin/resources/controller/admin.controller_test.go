package controller

import (
	"net/http"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	"stori-service/src/utils"
	"stori-service/src/utils/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DData struct {
	Client  string `groups:"client"`
	Admin   string `groups:"admin"`
	Console string `groups:"console"`
}

var defaultController AdminController = AdminController{}
var defaultMessage string = "Test"
var defaultErrors []map[string]string = []map[string]string{
	{"test01": "Test01", "test02": "Test02"},
	{"test03": "Test03", "test04": "Test04"},
}
var defaultData = DData{Client: "Client", Admin: "Admin", Console: "Console"}

func TestBaseController_MakePaginateResponse(t *testing.T) {
	t.Run("Should success on", func(t *testing.T) {
		t.Run("Correct header data", func(t *testing.T) {
			// Fixture
			page := &dto.Pagination{}

			// Run Foo inside request
			response := mock.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						defaultController.MakePaginateResponse(
							response,
							defaultData,
							http.StatusOK,
							page)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := utils.GetBodyResponse(response, &DData{})

			// assert data
			assert.Equal(t, defaultData.Admin, body.Data.(*DData).Admin)
			assert.Zero(t, body.Data.(*DData).Client)
			assert.Zero(t, body.Data.(*DData).Console)
		})
	})
}

func TestBaseController_MakeSuccessResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct OK response", func(t *testing.T) {
			// Run Foo inside request
			response := mock.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						defaultController.MakeSuccessResponse(
							response,
							defaultData,
							http.StatusOK,
							defaultMessage,
						)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := utils.GetBodyResponse(response, &DData{})

			// Assert Data
			assert.Equal(t, defaultData.Admin, body.Data.(*DData).Admin)
			assert.Zero(t, body.Data.(*DData).Client)
			assert.Zero(t, body.Data.(*DData).Console)
		})
	})
}

func TestBaseController_MakeErrorResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct Error response", func(t *testing.T) {
			// Run Foo inside request
			response := mock.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						defaultController.MakeErrorResponse(
							response,
							errors.ErrNotFound,
						)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := utils.GetBodyResponse(response, &DData{})

			// Assert Data
			assert.Nil(t, body.Data)
		})
	})
}
