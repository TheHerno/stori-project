package mock

import "stori-service/src/environments/common/resources/entity"

/*
ClientCustomerRepository is a ICustomerRepository mock
*/
type ClientCustomerRepository struct {
	TransactionalRepository
}

/*
FindAndLockByCustomerID mock method
*/
func (mock *ClientCustomerRepository) FindAndLockByCustomerID(customerid int) (*entity.Customer, error) {
	args := mock.Called(customerid)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

/*
FindByCustomerID mock method
*/
func (mock *ClientCustomerRepository) FindByCustomerID(customerid int) (*entity.Customer, error) {
	args := mock.Called(customerid)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}
