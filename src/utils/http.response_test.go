package utils

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"stori-service/src/libs/dto"
	myErrors "stori-service/src/libs/errors"
	"stori-service/src/utils/constant"
	customMock "stori-service/src/utils/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DData struct {
	ID   string `groups:"client"`
	Name string `groups:"admin"`
}

var defaultMessage string = "Test"
var defaultErrors []map[string]string = []map[string]string{
	{"test01": "Test01", "test02": "Test02"},
	{"test03": "Test03", "test04": "Test04"},
}
var defaultData = DData{ID: "1", Name: "Name"}

func TestMakePaginateResponse(t *testing.T) {
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
						MakePaginateResponse(
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

func TestMakeSuccessResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Correct OK response", func(t *testing.T) {
			// Run Foo inside request
			response := customMock.MHTTPHandle("GET", "/",
				func(response http.ResponseWriter, request *http.Request) {
					assert.NotPanics(t, func() {
						MakeSuccessResponse(
							response,
							defaultData,
							http.StatusOK,
							defaultData.Name,
						)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := GetBodyResponse(response, &DData{})

			// Assert Data
			assert.Equal(t, http.StatusOK, response.StatusCode)
			assert.Equal(t, defaultData.Name, body.Message)
			assert.Equal(t, []map[string]string{}, body.Errors)
			assert.Equal(t, defaultData.ID, body.Data.(*DData).ID)
		})
	})
}

func TestMakeErrorResponse(t *testing.T) {
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
						MakeErrorResponse(
							response,
							myErrors.ErrNotFound,
						)
					})
				}, "", nil, nil)
			defer response.Body.Close()
			body, _ := GetBodyResponse(response, &DData{})

			// Assert Data
			assert.Equal(t, http.StatusNotFound, response.StatusCode)
			assert.Equal(t, myErrors.ErrNotFound.Error(), body.Message)
			assert.Equal(t, expError, body.Errors)
			assert.Nil(t, body.Data)
		})
	})
}

func TestGetBodyRequest(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		t.Run("Valid JSON struct", func(t *testing.T) {
			// Fixture
			type Body struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}
			reqBody := []byte(`{
					"name":"Name","age":20
				}`)

			// Prepare Request
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
			body := Body{}
			err := GetBodyRequest(req, &body)

			// Response Assert
			assert.NoError(t, err)
			assert.Equal(t, "Name", body.Name)
			assert.Equal(t, 20, body.Age)
		})
	})
	t.Run("Should Fail", func(t *testing.T) {
		t.Run("Invalid JSON struct", func(t *testing.T) {
			// Fixture
			type Body struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}
			reqBody := []byte(`{
					"name":"Nam
				}`)

			// Prepare Request
			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(reqBody))
			body := Body{}
			err := GetBodyRequest(req, &body)

			// Response Assert
			assert.Error(t, err)
			assert.Equal(t, Body{}, body)
		})
	})
}

func TestGetBodyResponse(t *testing.T) {
	t.Run("Should Succeed", func(t *testing.T) {
		// Fixture
		type Body struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		reqBody := []byte(`{
				"Message":"Message",
				"Errors":[{"a":"aa","b":"bb"}],
				"Data":{"name":"Name","age":20}
			}`)

		// Prepare request
		response := customMock.MHTTPHandle("GET", "/", func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusOK)
			response.Write(reqBody)
		}, "", nil, nil)

		// Get response
		data := &Body{}
		bodyResponse, err := GetBodyResponse(response, data)

		// Data assertion
		assert.NoError(t, err)
		assert.Equal(t, "Message", bodyResponse.Message)
		assert.Equal(t, []map[string]string{{"a": "aa", "b": "bb"}}, bodyResponse.Errors)
		assert.Equal(t, "Name", data.Name)
		assert.Equal(t, 20, data.Age)
	})
	t.Run("Should Fail", func(t *testing.T) {
		// Fixture
		type Body struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		reqBody := []byte(`{
				"Message":"Messa
				"Errors":[{"a":"aa"
				"Data":{"name":"Na
			}`)

		// Prepare request
		response := customMock.MHTTPHandle("GET", "/", func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusOK)
			response.Write(reqBody)
		}, "", nil, nil)

		// Get response
		data := &Body{}
		bodyResponse, err := GetBodyResponse(response, data)

		// Data assertion
		assert.Error(t, err)
		assert.Nil(t, bodyResponse)
	})
}

func TesttakeBodyResponse(t *testing.T) {

	// Fixture
	type Body struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	reqBody := []byte(`{
		"name": "Name",
		"age": 20
	}`)

	t.Run("Should Succeed", func(t *testing.T) {

		// Prepare request
		response := customMock.MHTTPHandle("GET", "/", func(response http.ResponseWriter, request *http.Request) {
			response.WriteHeader(http.StatusOK)
			response.Write(reqBody)
		}, "", nil, nil)

		// Get response
		data := &Body{}
		err := takeBodyResponse(response, data)

		// Data assertion
		assert.NoError(t, err)
		assert.Equal(t, "Name", data.Name)
		assert.Equal(t, 20, data.Age)
	})
	t.Run("Should Fail", func(t *testing.T) {
		t.Run("Bad JSON", func(t *testing.T) {
			// Fixture
			reqBody := []byte(`{"name": "Name",`)

			// Prepare request
			response := customMock.MHTTPHandle("GET", "/", func(response http.ResponseWriter, request *http.Request) {
				response.WriteHeader(http.StatusOK)
				response.Write(reqBody)
			}, "", nil, nil)

			// Get response
			data := &Body{}
			err := takeBodyResponse(response, data)

			// Data assertion
			assert.Error(t, err)
			assert.Empty(t, data)
		})
		t.Run("Invalid destination variable type", func(t *testing.T) {

			// Prepare request
			response := customMock.MHTTPHandle("GET", "/", func(response http.ResponseWriter, request *http.Request) {
				response.WriteHeader(http.StatusOK)
				response.Write(reqBody)
			}, "", nil, nil)

			// Get response
			data := make(chan bool)
			err := takeBodyResponse(response, data)

			// Data assertion
			assert.Error(t, err)
		})
	})
}

func TestSetAction(t *testing.T) {
	testCases := []struct {
		testName             string
		err                  error
		expectedHeaderAction string
	}{
		{
			testName:             "Custom error without action",
			err:                  myErrors.ErrNotFound,
			expectedHeaderAction: "", //Empty
		},
		{
			testName:             "Generic error",
			err:                  errors.New("Generic error"),
			expectedHeaderAction: "", //Empty
		},
	}
	for _, tC := range testCases {
		t.Run(tC.testName, func(t *testing.T) {
			writer := httptest.NewRecorder()
			SetActionNeeded(tC.err, writer)
			headerAction := writer.Header().Get(constant.HeaderNeedsAction)

			//Data Assertion
			assert.Equal(t, tC.expectedHeaderAction, headerAction)
		})
	}
}
