package scopes

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"

	"gorm.io/gorm"
)

//StocksByWarehouseID scope function to get stocks by warehouse id
func StocksByWarehouseID(warehouseID int, pagination *dto.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&entity.StockMovement{}).
			Where(&entity.StockMovement{WarehouseID: warehouseID}).
			Where("product.enabled = true AND product.deleted_at IS NULL").
			Select("DISTINCT ON (stock_movement.product_id) product.product_id, stock_movement.available as stock, product.name, product.slug, product.description").
			Joins("JOIN product ON product.product_id = stock_movement.product_id").
			Order("stock_movement.product_id ASC, stock_movement_id DESC").
			Offset(pagination.Offset()).
			Limit(pagination.PageSize)
	}
}

//StockMovementByProductAndWarehouseID scope function to get stock movement by product and warehouse id
func StockMovementByWarehouseAndProductID(warehouseID int, productID int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&entity.StockMovement{}).
			Where(&entity.StockMovement{WarehouseID: warehouseID, ProductID: productID})
	}
}

//StockMovementByProductID scope function to get stock movement by product and warehouse id
func StockMovementByProductID(productID int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&entity.StockMovement{}).
			Where(&entity.StockMovement{ProductID: productID})
	}
}

//StockMovementByWarehouseID scope function to get stock movement warehouse id
func StockMovementByWarehouseID(warehouseID int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&entity.StockMovement{}).
			Where(&entity.StockMovement{WarehouseID: warehouseID})
	}
}
