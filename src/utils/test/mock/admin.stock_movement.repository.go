package mock

import (
	"stori-service/src/environments/common/resources/entity"
)

/*
AdminStockMovementRepository is a IStockMovementRepository mock
*/
type AdminStockMovementRepository struct {
	TransactionalRepository
}

/*
FindLastMovementByProductID mock method
*/
func (mock *AdminStockMovementRepository) FindLastMovementByProductID(productID int) (*entity.StockMovement, error) {
	args := mock.Called(productID)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.StockMovement), args.Error(1)
	}
	return nil, args.Error(1)
}

/*
FindLastMovementByWarehouseID mock method
*/
func (mock *AdminStockMovementRepository) FindLastMovementByWarehouseID(warehouseID int) (*entity.StockMovement, error) {
	args := mock.Called(warehouseID)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.StockMovement), args.Error(1)
	}
	return nil, args.Error(1)
}
