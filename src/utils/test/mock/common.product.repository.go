package mock

import "stori-service/src/environments/common/resources/entity"

/*
CommonProductRepository is a IProductRepository mock
*/
type CommonProductRepository struct {
	TransactionalRepository
}

/*
FindAndLockByID mock method
*/
func (mock *CommonProductRepository) FindAndLockByID(productID int) (*entity.Product, error) {
	args := mock.Called(productID)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}
