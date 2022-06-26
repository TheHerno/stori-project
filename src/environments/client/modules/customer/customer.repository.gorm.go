package customer

import (
	goerrors "errors"
	"stori-service/src/environments/client/resources/interfaces"
	database "stori-service/src/environments/common/resources/database/transaction"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
struct that implements ICustomerRepository
*/
type customerGormRepo struct {
	database.TransactionalGORMRepository
}

/*
NewCustomerGormRepo creates a new repo and returns IMovementRepository,
so it needs to implement all its methods
*/
func NewCustomerGormRepo(gormDb *gorm.DB) interfaces.ICustomerRepository {
	rCustomer := &customerGormRepo{}
	rCustomer.DB = gormDb
	return rCustomer
}

/*
findByCustomeridAndMayLock finds a customer by its Customerid and locks it if lock is true
*/
func (r *customerGormRepo) findByCustomeridAndMayLock(customerid int, lock bool) (*entity.Customer, error) {
	db := r.DB
	if lock {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var customer entity.Customer
	err := db.
		Where("customer_id", customerid).
		Take(&customer).Error
	if goerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &customer, err
}

/*
FindAndLockByCustomerid finds and locks a customer by its ID
*/
func (r *customerGormRepo) FindAndLockByCustomerid(customerid int) (*entity.Customer, error) {
	return r.findByCustomeridAndMayLock(customerid, true)
}

/*
FindByCustomerid finds a customer by its ID
*/
func (r *customerGormRepo) FindByCustomerid(customerid int) (*entity.Customer, error) {
	return r.findByCustomeridAndMayLock(customerid, false)
}

/*
Clone returns a new instance of the repository
*/
func (r *customerGormRepo) Clone() interface{} {
	return NewCustomerGormRepo(r.DB)
}
