package product

import (
	"stori-service/src/environments/client/resources/interfaces"
	"stori-service/src/libs/dto"
)

/*
Struct that implements IProductService
*/
type ProductService struct {
	rWarehouse     interfaces.IWarehouseRepository
	rStockMovement interfaces.IStockMovementRepository
}

/*
	NewProductService creates a new service, receives repository by dependency injection
	and returns IRepositoryService, so it needs to implement all its methods
*/
func NewProductService(rWarehouse interfaces.IWarehouseRepository, rStockMovement interfaces.IStockMovementRepository) interfaces.IProductService {
	return &ProductService{rWarehouse, rStockMovement}
}

/*
	GetStockList takes the userID and returns a list of products with stock
*/
func (s *ProductService) GetStockList(userID int, pagination *dto.Pagination) ([]dto.ProductWithStock, error) {
	warehouse, err := s.rWarehouse.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return s.rStockMovement.FindStocksByWarehouse(warehouse.WarehouseID, pagination)
}
