package mock

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"

	"github.com/stretchr/testify/mock"
)

type AdminWarehouseService struct {
	mock.Mock
}

// Index mock method
func (mock *AdminWarehouseService) Index(pagination *dto.Pagination) (*[]entity.Warehouse, error) {
	args := mock.Called(pagination)
	result := args.Get(0)
	if result != nil {
		return result.(*[]entity.Warehouse), args.Error(1)
	}
	return nil, args.Error(1)
}

// FindByID mock method
func (mock *AdminWarehouseService) FindByID(id int) (*entity.Warehouse, error) {
	args := mock.Called(id)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Warehouse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Update mock method
func (mock *AdminWarehouseService) Update(warehouse *dto.UpdateWarehouse) (*entity.Warehouse, error) {
	args := mock.Called(warehouse)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Warehouse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Create mock method
func (mock *AdminWarehouseService) Create(warehouse *dto.CreateWarehouse) (*entity.Warehouse, error) {
	args := mock.Called(warehouse)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Warehouse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Delete mock method
func (mock *AdminWarehouseService) Delete(id int) error {
	args := mock.Called(id)
	return args.Error(0)
}
