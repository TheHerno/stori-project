package warehouse

import (
	goerrors "errors"
	database "stori-service/src/environments/common/resources/database/transaction"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/environments/common/resources/interfaces"
	"stori-service/src/libs/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
struct that implements IWarehouseRepository
*/
type warehouseGormRepo struct {
	database.TransactionalGORMRepository
}

/*
NewWarehouseGormRepo creates a new repo and returns IWarehouseRepository,
so it needs to implement all its methods
*/
func NewWarehouseGormRepo(gormDb *gorm.DB) interfaces.IWarehouseRepository {
	rWarehouse := &warehouseGormRepo{}
	rWarehouse.DB = gormDb
	return rWarehouse
}

/*
findByIDAndMayLock finds a warehouse by its ID and locks it if lock is true
*/
func (r *warehouseGormRepo) findByIDAndMayLock(id int, lock bool) (*entity.Warehouse, error) {
	db := r.DB
	if lock {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	warehouse := &entity.Warehouse{}
	err := db.
		Take(warehouse, id).Error
	if goerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return warehouse, nil
}

/*
FindByID receives and id and finds the warehouse
*/
func (r *warehouseGormRepo) FindByID(id int) (*entity.Warehouse, error) {
	return r.findByIDAndMayLock(id, false)
}

/*
FindAndLockByID finds and locks a warehouse by its ID
*/
func (r *warehouseGormRepo) FindAndLockByID(id int) (*entity.Warehouse, error) {
	return r.findByIDAndMayLock(id, true)
}

/*
Clone returns a new instance of the repository
*/
func (r *warehouseGormRepo) Clone() interface{} {
	return NewWarehouseGormRepo(r.DB)
}
