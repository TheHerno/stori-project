package mock

import "stori-service/src/environments/common/resources/entity"

/*
CommonWarehouseRepository is a IWarehouseRepository mock
*/
type CommonWarehouseRepository struct {
	TransactionalRepository
}

// FindByID mock method
func (mock *CommonWarehouseRepository) FindByID(id int) (*entity.Warehouse, error) {
	args := mock.Called(id)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Warehouse), args.Error(1)
	}
	return nil, args.Error(1)
}

// FindAndLockByID mock method
func (mock *CommonWarehouseRepository) FindAndLockByID(id int) (*entity.Warehouse, error) {
	args := mock.Called(id)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Warehouse), args.Error(1)
	}
	return nil, args.Error(1)
}
