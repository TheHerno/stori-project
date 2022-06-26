package interfaces

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/environments/common/resources/interfaces"
)

/*
	ICustomerRepository to interact with entity and database
*/
type ICustomerRepository interface {
	interfaces.ITransactionalRepository
	FindAndLockByCustomerID(id int) (*entity.Customer, error)
	FindByCustomerID(id int) (*entity.Customer, error)
}
