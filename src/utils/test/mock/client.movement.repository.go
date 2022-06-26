package mock

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
)

/*
ClientMovementRepository is a IMovementRepository mock
*/
type ClientMovementRepository struct {
	TransactionalRepository
}

// Create mock method
func (mock *ClientMovementRepository) Create(movement *entity.Movement) (*entity.Movement, error) {
	args := mock.Called(movement)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Movement), args.Error(1)
	}
	return nil, args.Error(1)
}

// FindLastMovement mock method
func (mock *ClientMovementRepository) FindLastMovement(userID int, productID int) (*entity.Movement, error) {
	args := mock.Called(userID, productID)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Movement), args.Error(1)
	}
	return nil, args.Error(1)
}

//FindStockByUser mock method
func (mock *ClientMovementRepository) FindStocksByUser(userID int, pagination *dto.Pagination) ([]dto.ProductWithStock, error) {
	args := mock.Called(userID, pagination)
	result := args.Get(0)
	if result != nil {
		return result.([]dto.ProductWithStock), args.Error(1)
	}
	return nil, args.Error(1)
}
