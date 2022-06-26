package mock

import (
	"stori-service/src/environments/common/resources/entity"
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
func (mock *ClientMovementRepository) FindLastMovement(customerid int, productID int) (*entity.Movement, error) {
	args := mock.Called(customerid, productID)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Movement), args.Error(1)
	}
	return nil, args.Error(1)
}
