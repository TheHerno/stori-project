package mock

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
)

/*
ClientStockMovementRepository is a IStockMovementRepository mock
*/
type ClientStockMovementRepository struct {
	TransactionalRepository
}

// Create mock method
func (mock *ClientStockMovementRepository) Create(stockMovement *entity.StockMovement) (*entity.StockMovement, error) {
	args := mock.Called(stockMovement)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.StockMovement), args.Error(1)
	}
	return nil, args.Error(1)
}

// FindLastStockMovement mock method
func (mock *ClientStockMovementRepository) FindLastStockMovement(warehouseID int, productID int) (*entity.StockMovement, error) {
	args := mock.Called(warehouseID, productID)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.StockMovement), args.Error(1)
	}
	return nil, args.Error(1)
}

//FindStockByWarehouse mock method
func (mock *ClientStockMovementRepository) FindStocksByWarehouse(warehouseID int, pagination *dto.Pagination) ([]dto.ProductWithStock, error) {
	args := mock.Called(warehouseID, pagination)
	result := args.Get(0)
	if result != nil {
		return result.([]dto.ProductWithStock), args.Error(1)
	}
	return nil, args.Error(1)
}
