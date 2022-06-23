package stock

import (
	goerrors "errors"
	"stori-service/src/environments/admin/resources/interfaces"
	database "stori-service/src/environments/common/resources/database/transaction"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database/scopes"
	"stori-service/src/libs/errors"

	"gorm.io/gorm"
)

/*
struct that implements IStockMovementRepository
*/
type stockmovementGormRepo struct {
	database.TransactionalGORMRepository
}

/*
NewStockMovementGormRepo creates a new repo and returns IStockMovementRepository,
so it needs to implement all its methods
*/
func NewStockMovementGormRepo(gormDb *gorm.DB) interfaces.IStockMovementRepository {
	rStockMovement := &stockmovementGormRepo{}
	rStockMovement.DB = gormDb
	return rStockMovement
}

/*
FindLastMovementBySomeID find the last movement of a produdct or warehouse.
First arg means warehouse if 0 or product if 1. Second arg is the id
*/
func (r *stockmovementGormRepo) findLastMovementBySomeID(typee int, id int) (*entity.StockMovement, error) {
	var stockMovement entity.StockMovement
	if id <= 0 {
		return nil, errors.ErrNotFound
	}
	scope := scopes.StockMovementByWarehouseID
	if typee == 1 {
		scope = scopes.StockMovementByProductID
	}
	err := r.DB.Scopes(scope(id)).
		Order("stock_movement_id DESC").
		Take(&stockMovement).Error

	if goerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &stockMovement, nil
}

/*
FindLastMovementByProductID finds the last movement of a product
*/
func (r *stockmovementGormRepo) FindLastMovementByProductID(productID int) (*entity.StockMovement, error) {
	return r.findLastMovementBySomeID(1, productID)
}

/*
FindLastMovementByWarehouseID finds the last movement of a warehouse
*/
func (r *stockmovementGormRepo) FindLastMovementByWarehouseID(warehouseID int) (*entity.StockMovement, error) {
	return r.findLastMovementBySomeID(0, warehouseID)
}

/*
Clone returns a new instance of the repository
*/
func (r *stockmovementGormRepo) Clone() interface{} {
	return NewStockMovementGormRepo(r.DB)
}
