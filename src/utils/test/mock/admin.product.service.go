package mock

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"

	"github.com/stretchr/testify/mock"
)

type AdminProductService struct {
	mock.Mock
}

// Index mock method
func (mock *AdminProductService) Index(pagination *dto.Pagination) ([]entity.Product, error) {
	args := mock.Called(pagination)
	result := args.Get(0)
	if result != nil {
		return result.([]entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// FindByID mock method
func (mock *AdminProductService) FindByID(id int) (*entity.Product, error) {
	args := mock.Called(id)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// Update mock method
func (mock *AdminProductService) Update(product *dto.UpdateProduct) (*entity.Product, error) {
	args := mock.Called(product)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// Create mock method
func (mock *AdminProductService) Create(product *dto.CreateProduct) (*entity.Product, error) {
	args := mock.Called(product)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// Delete mock method
func (mock *AdminProductService) Delete(id int) error {
	args := mock.Called(id)
	return args.Error(0)
}
