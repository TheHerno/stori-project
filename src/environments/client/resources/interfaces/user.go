package interfaces

import (
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/environments/common/resources/interfaces"
)

/*
	IUserRepository to interact with entity and database
*/
type IUserRepository interface {
	interfaces.ITransactionalRepository
	FindAndLockByUserID(id int) (*entity.User, error)
	FindByUserID(id int) (*entity.User, error)
}
