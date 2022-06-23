package warehouse

import (
	goerrors "errors"
	"stori-service/src/environments/admin/resources/interfaces"
	"stori-service/src/environments/common/resources/entity"
	commonInterfaces "stori-service/src/environments/common/resources/interfaces"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
)

/*
	Struct that implements the IWarehouseService interface
*/
type warehouseService struct {
	rWarehouse     interfaces.IWarehouseRepository
	rCWarehouse    commonInterfaces.IWarehouseRepository
	rStockMovement interfaces.IStockMovementRepository
}

/*
	NewWarehouseService creates a new service, receives repository by dependency injection
	and returns IRepositoryService, so it needs to implement all its methods
*/
func NewWarehouseService(rWarehouse interfaces.IWarehouseRepository, rCWarehouse commonInterfaces.IWarehouseRepository, rStockMovement interfaces.IStockMovementRepository) interfaces.IWarehouseService {
	return &warehouseService{rWarehouse, rCWarehouse, rStockMovement}
}

/*
	Index receives pagination and calls Index from repository
*/
func (s *warehouseService) Index(pagination *dto.Pagination) (*[]entity.Warehouse, error) {
	return s.rWarehouse.Index(pagination)
}

/*
	Update receives an UpdateWarehouse dto and validates it
	Finds Warehouse by ID, replace new values and calls Update from repository
	If there is an error, returns it as a second result
*/
func (s *warehouseService) Update(updateWarehouse *dto.UpdateWarehouse) (*entity.Warehouse, error) {
	if updateWarehouse.WarehouseID < 1 {
		return nil, errors.ErrNotFound
	}
	if err := updateWarehouse.Validate(); err != nil {
		return nil, err
	}
	warehouse, err := s.FindByID(updateWarehouse.WarehouseID)
	if err != nil {
		return nil, err
	}
	warehouse.Name = updateWarehouse.Name
	warehouse.Address = updateWarehouse.Address
	return s.rWarehouse.Update(warehouse)
}

/*
	Create receives a warehouse, validates and creates it
*/
func (s *warehouseService) Create(createWarehouse *dto.CreateWarehouse) (*entity.Warehouse, error) {
	if err := createWarehouse.Validate(); err != nil {
		return nil, err
	}
	warehouse := createWarehouse.ParseToWarehouse()
	return s.rWarehouse.Create(warehouse)
}

/*
	Delete receives an ID and calls Delete from repository
*/
func (s *warehouseService) Delete(id int) error {
	if id < 1 {
		return errors.ErrNotFound
	}
	rStockMovement := s.rStockMovement.Clone().(interfaces.IStockMovementRepository)
	rWarehouse := s.rWarehouse.Clone().(interfaces.IWarehouseRepository)
	rCWarehouse := s.rCWarehouse.Clone().(commonInterfaces.IWarehouseRepository)
	tx := rWarehouse.Begin(nil)
	rCWarehouse.Begin(tx)
	rStockMovement.Begin(tx)
	defer rWarehouse.Rollback()

	_, err := rCWarehouse.FindAndLockByID(id)
	if err != nil {
		return err
	}
	movement, err := rStockMovement.FindLastMovementByWarehouseID(id)
	if movement != nil || !goerrors.Is(err, errors.ErrNotFound) { // si hay movimiento no se borra
		return errors.ErrCantDeleteWarehouse
	}

	err = rWarehouse.Delete(id)
	if err != nil {
		return err
	}
	err = rWarehouse.Commit()
	if err != nil {
		return err
	}
	return nil
}

/*
FindByID receives an ID and calls FindByID from repository
*/
func (s *warehouseService) FindByID(id int) (*entity.Warehouse, error) {
	if id < 1 {
		return nil, errors.ErrNotFound
	}
	return s.rCWarehouse.FindByID(id)
}
