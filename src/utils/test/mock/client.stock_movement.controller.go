package mock

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

/*
ClientStockMovementController is a IStockMovementController mock
*/
type ClientStockMovementController struct {
	mock.Mock
}

// Create mock method
func (mock *ClientStockMovementController) Create(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}
