package mock

import "stori-service/src/environments/common/resources/entity"

/*
ClientUserRepository is a IUserRepository mock
*/
type ClientUserRepository struct {
	TransactionalRepository
}

/*
FindAndLockByUserID mock method
*/
func (mock *ClientUserRepository) FindAndLockByUserID(userID int) (*entity.User, error) {
	args := mock.Called(userID)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

/*
FindByUserID mock method
*/
func (mock *ClientUserRepository) FindByUserID(userID int) (*entity.User, error) {
	args := mock.Called(userID)
	result := args.Get(0)
	if result != nil {
		return result.(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}
