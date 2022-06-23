package mock

import (
	"github.com/stretchr/testify/mock"
)

/*
TransactionalRepository is a ITransactionalRepository mock
*/
type TransactionalRepository struct {
	mock.Mock
}

/*
Begin mock method
*/
func (mock *TransactionalRepository) Begin(initialTx interface{}) interface{} {
	args := mock.Called(initialTx)
	return args.Get(0)
}

/*
Commit mock method
*/
func (mock *TransactionalRepository) Commit() error {
	args := mock.Called()
	return args.Error(0)
}

/*
SavePoint mock method
*/
func (mock *TransactionalRepository) SavePoint() error {
	args := mock.Called()
	return args.Error(0)
}

/*
Rollback mock method
*/
func (mock *TransactionalRepository) Rollback() error {
	args := mock.Called()
	return args.Error(0)
}

/*
Clone mock method
*/
func (mock *TransactionalRepository) Clone() interface{} {
	args := mock.Called()
	return args.Get(0)
}
