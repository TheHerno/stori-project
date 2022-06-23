package mock

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

/*
AdminWarehouseController is a IWarehouseController mock
*/
type AdminWarehouseController struct {
	mock.Mock
}

// Index mock method
func (mock *AdminWarehouseController) Index(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}

// FindByID mock method
func (mock *AdminWarehouseController) FindByID(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}

// Update mock method
func (mock *AdminWarehouseController) Update(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}

// Create mock method
func (mock *AdminWarehouseController) Create(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}

// Delete mock method
func (mock *AdminWarehouseController) Delete(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}
