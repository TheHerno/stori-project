package product

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"stori-service/src/utils/test/mock"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestNewProductRouter(t *testing.T) {
	t.Run("Routes with controller mock", func(t *testing.T) {
		t.Run("Should success on", func(t *testing.T) {
			testCases := []struct {
				Path    string
				Method  string
				Handler string
			}{
				{
					Path:    "",
					Method:  http.MethodGet,
					Handler: "GetStocks",
				},
			}

			for _, testCase := range testCases {
				t.Run(fmt.Sprintf("Method: %s Path: %s Handler: %s", testCase.Method, testCase.Path, testCase.Handler), func(t *testing.T) {
					muxRouter := mux.NewRouter()
					subRouterPath := "/test"
					subRouter := muxRouter.PathPrefix(subRouterPath).Subrouter()
					mockProductC := new(mock.ClientProductController)
					NewProductRouter(subRouter, mockProductC)
					mockProductC.On(
						testCase.Handler,
						testifyMock.AnythingOfType("*http.response"),
						testifyMock.AnythingOfType("*http.Request"),
					).Run(func(args testifyMock.Arguments) {
						firstArgument := args[0]
						response := firstArgument.(http.ResponseWriter)
						response.WriteHeader(http.StatusTeapot) //using teapot status, to ensure it in assertions
					})
					ts := httptest.NewServer(muxRouter)
					URL := fmt.Sprint(ts.URL, subRouterPath, testCase.Path)
					req, _ := http.NewRequest(testCase.Method, URL, nil)
					res, err := ts.Client().Do(req)

					// mock assertion: Behavioural
					mockProductC.AssertExpectations(t)
					mockProductC.AssertNumberOfCalls(t, testCase.Handler, 1)

					// data assertion
					assert.NoError(t, err)
					assert.NotNil(t, res)
					assert.Equal(t, http.StatusTeapot, res.StatusCode)
				})
			}
		})
	})
}
