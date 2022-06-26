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
	FindAndLockByCustomerid(id int) (*entity.Customer, error)
	FindByCustomerid(id int) (*entity.Customer, error)
}
