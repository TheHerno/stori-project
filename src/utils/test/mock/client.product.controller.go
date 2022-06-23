package mock

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

/*
ClientProductController is a IProductController mock
*/
type ClientProductController struct {
	mock.Mock
}

// GetStocks mock method
func (mock *ClientProductController) GetStocks(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}
