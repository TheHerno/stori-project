package interfaces

import (
	"stori-service/src/environments/common/resources/entity"
)

/*
	IProductRepository to interact with entity and database
*/
type IProductRepository interface {
	ITransactionalRepository
	FindAndLockByID(id int) (*entity.Product, error)
}
