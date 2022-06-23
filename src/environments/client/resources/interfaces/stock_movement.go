package interfaces

import (
	"net/http"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/environments/common/resources/interfaces"
	"stori-service/src/libs/dto"
)

/*
IStockMovementRepository to interact with entity and database
*/
type IStockMovementRepository interface {
	interfaces.ITransactionalRepository
	Create(stockMovement *entity.StockMovement) (*entity.StockMovement, error)
	FindLastStockMovement(warehouseID int, productID int) (*entity.StockMovement, error)
	FindStocksByWarehouse(warehouseID int, pagination *dto.Pagination) ([]dto.ProductWithStock, error)
}

/*
	IStockMovementService methods with bussiness logic
*/
type IStockMovementService interface {
	Create(stockMovement *dto.NewStockMovement) (*entity.StockMovement, error)
}

/*
	IStockMovementController methods to handle requests and responses
*/
type IStockMovementController interface {
	Create(response http.ResponseWriter, request *http.Request)
}
