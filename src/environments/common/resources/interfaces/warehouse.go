package interfaces

import "stori-service/src/environments/common/resources/entity"

/*
	IWarehouseRepository to interact with entity and database
*/
type IWarehouseRepository interface {
	ITransactionalRepository
	FindByID(id int) (*entity.Warehouse, error)
	FindAndLockByID(id int) (*entity.Warehouse, error)
}
