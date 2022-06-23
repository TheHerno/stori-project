package interfaces

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/environments/common/resources/interfaces"
)

/*
	IWarehouseRepository to interact with entity and database
*/
type IWarehouseRepository interface {
	interfaces.ITransactionalRepository
	FindAndLockByUserID(id int) (*entity.Warehouse, error)
	FindByUserID(id int) (*entity.Warehouse, error)
}
