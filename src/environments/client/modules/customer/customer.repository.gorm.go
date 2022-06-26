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
FindAndLockByCustomerID it's a partial application of findAndMayLockByCustomerID with lock argument set to true
*/
func (r *customerGormRepo) FindAndLockByCustomerID(customerId int) (*entity.Customer, error) {
	return r.findAndMayLockByCustomerID(customerId, true)
}

/*
FindByCustomerID it's a partial application of findAndMayLockByCustomerID with lock argument set to true
*/
func (r *customerGormRepo) FindByCustomerID(customerId int) (*entity.Customer, error) {
	return r.findAndMayLockByCustomerID(customerId, false)
}

/*
findAndMayLockByCustomerID returns a customer by its id and locks it if the second argument is true
*/
func (r *customerGormRepo) findAndMayLockByCustomerID(customerId int, lock bool) (*entity.Customer, error) {
	var customer entity.Customer
	db := r.DB
	if lock {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	err := db.Where("customer_id = ?", customerId).First(&customer).Error
	if goerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

/*
Clone returns a new instance of the repository
*/
func (r *customerGormRepo) Clone() interface{} {
	return NewCustomerGormRepo(r.DB)
}
