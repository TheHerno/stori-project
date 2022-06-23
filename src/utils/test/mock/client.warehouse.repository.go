package mock

import "stori-service/src/environments/common/resources/entity"

/*
ClientWarehouseRepository is a IWarehouseRepository mock
*/
type ClientWarehouseRepository struct {
	TransactionalRepository
}

/*
FindAndLockByUserID mock method
*/
func (mock *ClientWarehouseRepository) FindAndLockByUserID(userID int) (*entity.Warehouse, error) {
	args := mock.Called(userID)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Warehouse), args.Error(1)
	}
	return nil, args.Error(1)
}

/*
FindByUserID mock method
*/
func (mock *ClientWarehouseRepository) FindByUserID(userID int) (*entity.Warehouse, error) {
	args := mock.Called(userID)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Warehouse), args.Error(1)
	}
	return nil, args.Error(1)
}
