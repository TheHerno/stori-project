package interfaces

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/environments/common/resources/interfaces"
)

/*
IStockMovementRepository to interact with entity and database
*/
type IStockMovementRepository interface {
	interfaces.ITransactionalRepository
	FindLastMovementByProductID(productID int) (*entity.StockMovement, error)
	FindLastMovementByWarehouseID(warehouseID int) (*entity.StockMovement, error)
}
