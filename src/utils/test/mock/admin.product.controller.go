package mock

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

/*
AdminProductController is a IProductController mock
*/
type AdminProductController struct {
	mock.Mock
}

// Index mock method
func (mock *AdminProductController) Index(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}

// FindByID mock method
func (mock *AdminProductController) FindByID(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}

// Update mock method
func (mock *AdminProductController) Update(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}

// Create mock method
func (mock *AdminProductController) Create(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}

// Delete mock method
func (mock *AdminProductController) Delete(response http.ResponseWriter, request *http.Request) {
	mock.Called(response, request)
}
