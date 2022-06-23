package stock

import (
	goerrors "errors"
	"stori-service/src/environments/client/modules/stock/providers/params"
	"stori-service/src/environments/client/resources/interfaces"
	"stori-service/src/environments/common/resources/entity"
	commonInterfaces "stori-service/src/environments/common/resources/interfaces"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	"stori-service/src/utils/constant"
)

var getParamInt = params.GetParamInt

/*
Struct that implements IStockMovementService
*/
type StockMovementService struct {
	rStockMovement interfaces.IStockMovementRepository
	rWarehouse     interfaces.IWarehouseRepository
	rCWarehouse    commonInterfaces.IWarehouseRepository
	rProduct       commonInterfaces.IProductRepository
}

/*
	NewStockMovementService creates a new service, receives repository by dependency injection
	and returns IRepositoryService, so it needs to implement all its methods
*/
func NewStockMovementService(rStockMovement interfaces.IStockMovementRepository, rWarehouse interfaces.IWarehouseRepository, rCWarehouse commonInterfaces.IWarehouseRepository, rProduct commonInterfaces.IProductRepository) interfaces.IStockMovementService {
	return &StockMovementService{rStockMovement, rWarehouse, rCWarehouse, rProduct}
}

/*
getAvailableStock finds the last stock movement and returns the available stock
*/
func (s *StockMovementService) getAvailableStock(warehouseID int, productID int) (int, error) {
	lastMovement, err := s.rStockMovement.FindLastStockMovement(warehouseID, productID)
	if goerrors.Is(err, errors.ErrNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return lastMovement.Available, nil
}

/*
save calculates the new available, validates it,  and saves the new stock movement
*/
func (s *StockMovementService) save(rStockMovement interfaces.IStockMovementRepository, newStockMovement *dto.NewStockMovement, warehouse *entity.Warehouse) (*entity.StockMovement, error) {
	lastAvailable, err := s.getAvailableStock(warehouse.WarehouseID, newStockMovement.ProductID)
	if err != nil {
		return nil, err
	}

	newAvailable := lastAvailable + newStockMovement.Quantity*newStockMovement.Type
	if newAvailable < 0 {
		return nil, errors.ErrStockMovementInvalid
	}

	stockMovementToCreate := &entity.StockMovement{
		ProductID:   newStockMovement.ProductID,
		WarehouseID: warehouse.WarehouseID,
		Quantity:    newStockMovement.Quantity,
		Available:   newAvailable,
		Type:        newStockMovement.Type,
		Concept:     newStockMovement.Concept,
	}

	return rStockMovement.Create(stockMovementToCreate)
}

/*
checkQuantity gets the max allowed stock movement quantity and returns a boolean indicating if the quantity is valid
*/
func (s *StockMovementService) checkQuantity(quantity int) (bool, error) {
	maxQty, err := getParamInt("max_stock_to_move")
	if err != nil {
		return false, err
	}
	return quantity <= maxQty, nil
}

/*
Create takes a newStockMovementDTO, validates it, locks and validate the warehouse and product, and then create the stock movement.
*/
func (s *StockMovementService) Create(newStockMovement *dto.NewStockMovement) (*entity.StockMovement, error) {
	err := newStockMovement.Validate()
	if err != nil {
		return nil, err
	}
	isTransfer := newStockMovement.WarehouseID != 0
	if isTransfer && newStockMovement.Type != constant.OutcomeType {
		return nil, errors.ErrStockMovementInvalid
	}
	validQty, err := s.checkQuantity(newStockMovement.Quantity)
	if err != nil {
		return nil, err
	}
	if !validQty {
		return nil, errors.ErrStockMovementInvalid
	}
	rWarehouse := s.rWarehouse.Clone().(interfaces.IWarehouseRepository)
	rCWarehouse := s.rCWarehouse.Clone().(commonInterfaces.IWarehouseRepository)
	rProduct := s.rProduct.Clone().(commonInterfaces.IProductRepository)
	rStockMovement := s.rStockMovement.Clone().(interfaces.IStockMovementRepository)

	tx := rStockMovement.Begin(nil)
	rCWarehouse.Begin(tx)
	rProduct.Begin(tx)
	rWarehouse.Begin(tx)
	defer rStockMovement.Rollback()

	warehouse, err := rWarehouse.FindAndLockByUserID(newStockMovement.UserID)
	if err != nil {
		return nil, err
	}
	if isTransfer && warehouse.WarehouseID == newStockMovement.WarehouseID {
		return nil, errors.ErrStockMovementInvalid
	}
	product, err := rProduct.FindAndLockByID(newStockMovement.ProductID)
	if err != nil {
		return nil, err
	}
	if !*product.Enabled {
		return nil, errors.ErrProductNotEnabled
	}
	createdStockMovement, err := s.save(rStockMovement, newStockMovement, warehouse)
	if err != nil {
		return nil, err
	}
	if isTransfer {
		_, err = s.createTransfer(newStockMovement, warehouse, rCWarehouse, rStockMovement)
		if err != nil {
			return nil, err
		}
	}
	err = rStockMovement.Commit()
	if err != nil {
		return nil, err
	}
	return createdStockMovement, nil
}

/*
createTransfer creates a new stock movement for a transfer
*/
func (s *StockMovementService) createTransfer(newStockMovement *dto.NewStockMovement, warehouse *entity.Warehouse, rCWarehouse commonInterfaces.IWarehouseRepository, rStockMovement interfaces.IStockMovementRepository) (*entity.StockMovement, error) {
	warehouseTarget, err := rCWarehouse.FindAndLockByID(newStockMovement.WarehouseID)
	if err != nil {
		return nil, err
	}
	newTargetStockMovement := &dto.NewStockMovement{
		ProductID:   newStockMovement.ProductID,
		WarehouseID: newStockMovement.WarehouseID,
		Quantity:    newStockMovement.Quantity,
		Type:        1,
		Concept:     "Transfer from " + warehouse.Name,
	}
	receiverMovement, err := s.save(rStockMovement, newTargetStockMovement, warehouseTarget)
	if err != nil {
		return nil, err
	}
	return receiverMovement, nil
}
