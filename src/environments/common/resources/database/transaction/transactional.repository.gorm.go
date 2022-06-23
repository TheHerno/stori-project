package database

import (
	"fmt"

	"gorm.io/gorm"
)

/*
TransactionalGORMRepository is a struct that implements ITransactionalRepository
*/
type TransactionalGORMRepository struct {
	DB        *gorm.DB
	savePoint int
}

/*
Begin begins transaction for current connection
*/
func (r *TransactionalGORMRepository) Begin(initialTx interface{}) interface{} {
	if initialTx == nil {
		r.DB = r.DB.Begin()
	} else {
		r.DB = initialTx.(*gorm.DB)
	}
	return r.DB
}

/*
Commit commits transaction for current connection
*/
func (r *TransactionalGORMRepository) Commit() error {
	r.savePoint = 0
	return r.DB.Commit().Error
}

/*
Rollback rollbacks transaction for current connection
*/
func (r *TransactionalGORMRepository) Rollback() error {
	if r.savePoint > 0 {
		err := r.DB.RollbackTo(fmt.Sprint("sp", r.savePoint)).Error
		if err != nil {
			return err
		}
		r.savePoint--
		return nil
	}
	return r.DB.Rollback().Error
}

/*
SavePoint creates a savepoint of the transaction for current connection
*/
func (r *TransactionalGORMRepository) SavePoint() error {
	err := r.DB.SavePoint(fmt.Sprint("sp", r.savePoint+1)).Error
	if err != nil {
		return err
	}
	r.savePoint++
	return nil
}
