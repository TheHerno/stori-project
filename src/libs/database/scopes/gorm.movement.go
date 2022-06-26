package scopes

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"

	"gorm.io/gorm"
)

//StocksByCustomerid scope function to get stocks by user id
func StocksByCustomerid(customerid int, pagination *dto.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&entity.Movement{}).
			Where(&entity.Movement{Customerid: customerid}).
			Where("product.enabled = true AND product.deleted_at IS NULL").
			Select("DISTINCT ON (movement.product_id) product.product_id, movement.available as stock, product.name, product.slug, product.description").
			Joins("JOIN product ON product.product_id = movement.product_id").
			Order("movement.product_id ASC, movement_id DESC").
			Offset(pagination.Offset()).
			Limit(pagination.PageSize)
	}
}

//MovementByCustomerid scope function to get stock movement user id
func MovementByCustomerid(customerid int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&entity.Movement{}).
			Where(&entity.Movement{Customerid: customerid})
	}
}
