package mock

import "stori-service/src/environments/common/resources/entity"

/*
ClientCustomerRepository is a ICustomerRepository mock
*/
type ClientCustomerRepository struct {
	TransactionalRepository
}

/*
FindAndLockByCustomerid mock method
*/
func (mock *ClientCustomerRepository) FindAndLockByCustomerid(customerid int) (*entity.Customer, error) {
	args := mock.Called(customerid)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

/*
FindByCustomerid mock method
*/
func (mock *ClientCustomerRepository) FindByCustomerid(customerid int) (*entity.Customer, error) {
	args := mock.Called(customerid)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}
