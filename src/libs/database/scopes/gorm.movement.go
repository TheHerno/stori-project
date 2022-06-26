package scopes

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"

	"gorm.io/gorm"
)

//StocksByUserID scope function to get stocks by user id
func StocksByUserID(userID int, pagination *dto.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&entity.Movement{}).
			Where(&entity.Movement{UserID: userID}).
			Where("product.enabled = true AND product.deleted_at IS NULL").
			Select("DISTINCT ON (movement.product_id) product.product_id, movement.available as stock, product.name, product.slug, product.description").
			Joins("JOIN product ON product.product_id = movement.product_id").
			Order("movement.product_id ASC, movement_id DESC").
			Offset(pagination.Offset()).
			Limit(pagination.PageSize)
	}
}

//MovementByUserID scope function to get stock movement user id
func MovementByUserID(userID int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&entity.Movement{}).
			Where(&entity.Movement{UserID: userID})
	}
}
