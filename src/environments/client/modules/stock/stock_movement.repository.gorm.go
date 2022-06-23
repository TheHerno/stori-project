package stock

import (
	goerrors "errors"
	"stori-service/src/environments/client/resources/interfaces"
	database "stori-service/src/environments/common/resources/database/transaction"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database/scopes"
	"stori-service/src/libs/dto"
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
Create receives the stockmovement to be created and creates it
If there is an error, returns it as a second result
*/
func (r *stockmovementGormRepo) Create(stockmovement *entity.StockMovement) (*entity.StockMovement, error) {
	err := r.DB.Create(stockmovement).Error
	if err != nil {
		return nil, err
	}

	return stockmovement, nil
}

/*
FindLastStockMovement finds the last stock movement of a product and a warehouse
*/
func (r *stockmovementGormRepo) FindLastStockMovement(warehouseID int, productID int) (*entity.StockMovement, error) {
	var stockMovement entity.StockMovement
	err := r.DB.Scopes(scopes.StockMovementByWarehouseAndProductID(warehouseID, productID)).
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
getStockCountByWarehouse returns the stock count of a product in a warehouse
*/
func (r *stockmovementGormRepo) getStockCountByWarehouse(warehouseID int) (int64, error) {
	var count int64
	// No encontramos como hacer que gorm nos deje hacer esto con subquery.
	err := r.DB.Raw(`SELECT COUNT(*) FROM (
		SELECT DISTINCT ON
		(stock_movement.product_id) product.product_id,
		stock_movement.available as stock,
		product.name,
		product.slug,
		product.description FROM "stock_movement"
		JOIN product
		ON product.product_id = stock_movement.product_id
		WHERE stock_movement.warehouse_id = ?
		AND (product.enabled = true AND product.deleted_at IS NULL)
		AND "stock_movement"."deleted_at" IS NULL
		ORDER BY stock_movement.product_id ASC,
		stock_movement_id DESC
		) AS count`, warehouseID).Scan(&count).Error
	if err != nil {
		return int64(0), err
	}
	return count, nil
}

/*
FindStockByWarehouse finds the stock of all products in a warehouse
*/
func (r *stockmovementGormRepo) FindStocksByWarehouse(warehouseID int, pagination *dto.Pagination) ([]dto.ProductWithStock, error) {
	var stock []dto.ProductWithStock
	count, err := r.getStockCountByWarehouse(warehouseID)
	if err != nil {
		return nil, err
	}
	pagination.TotalCount = count
	err = r.DB.
		Scopes(scopes.StocksByWarehouseID(warehouseID, pagination)).
		Scan(&stock).Error
	if err != nil {
		return nil, err
	}
	return stock, nil
}

/*
Clone returns a new instance of the repository
*/
func (r *stockmovementGormRepo) Clone() interface{} {
	return NewStockMovementGormRepo(r.DB)
}
