package mock

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
)

/*
AdminWarehouseRepository is a IWarehouseRepository mock
*/
type AdminWarehouseRepository struct {
	TransactionalRepository
}

// Index mock method
func (mock *AdminWarehouseRepository) Index(pagination *dto.Pagination) (*[]entity.Warehouse, error) {
	args := mock.Called(pagination)
	result := args.Get(0)
	if result != nil {
		return result.(*[]entity.Warehouse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Update mock method
func (mock *AdminWarehouseRepository) Update(warehouse *entity.Warehouse) (*entity.Warehouse, error) {
	args := mock.Called(warehouse)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Warehouse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Create mock method
func (mock *AdminWarehouseRepository) Create(warehouse *entity.Warehouse) (*entity.Warehouse, error) {
	args := mock.Called(warehouse)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Warehouse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Delete mock method
func (mock *AdminWarehouseRepository) Delete(id int) error {
	args := mock.Called(id)
	return args.Error(0)
}
