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

// BulkCreate mock method
func (mock *ClientMovementRepository) BulkCreate(movements []entity.Movement) error {
	args := mock.Called(movements)
	return args.Error(0)
}

// GetLastMovementByCustomerID mock method
func (mock *ClientMovementRepository) GetLastMovementByCustomerID(customerID int) (*entity.Movement, error) {
	args := mock.Called(customerID)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Movement), args.Error(1)
	}
	return nil, args.Error(1)
}
