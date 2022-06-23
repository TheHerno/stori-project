package warehouse

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
struct that implements IWarehouseRepository
*/
type warehouseGormRepo struct {
	database.TransactionalGORMRepository
}

/*
NewWarehouseGormRepo creates a new repo and returns IStockMovementRepository,
so it needs to implement all its methods
*/
func NewWarehouseGormRepo(gormDb *gorm.DB) interfaces.IWarehouseRepository {
	rWarehouse := &warehouseGormRepo{}
	rWarehouse.DB = gormDb
	return rWarehouse
}

/*
findByUserIDAndMayLock finds a warehouse by its UserID and locks it if lock is true
*/
func (r *warehouseGormRepo) findByUserIDAndMayLock(userID int, lock bool) (*entity.Warehouse, error) {
	db := r.DB
	if lock {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var warehouse entity.Warehouse
	err := db.
		Where("user_id", userID).
		Take(&warehouse).Error
	if goerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &warehouse, err
}

/*
FindAndLockByUserID finds and locks a warehouse by its ID
*/
func (r *warehouseGormRepo) FindAndLockByUserID(userID int) (*entity.Warehouse, error) {
	return r.findByUserIDAndMayLock(userID, true)
}

/*
FindByUserID finds and locks a warehouse by its ID
*/
func (r *warehouseGormRepo) FindByUserID(userID int) (*entity.Warehouse, error) {
	return r.findByUserIDAndMayLock(userID, false)
}

/*
Clone returns a new instance of the repository
*/
func (r *warehouseGormRepo) Clone() interface{} {
	return NewWarehouseGormRepo(r.DB)
}
