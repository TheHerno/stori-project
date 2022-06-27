package mock

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

/*
ClientMovementController is a IMovementController mock
*/
type ClientMovementController struct {
	mock.Mock
}

// ProcessFile mock method
func (mock *ClientMovementController) ProcessFile(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}
