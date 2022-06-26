package interfaces

import (
	"stori-service/src/environments/common/resources/entity"
	commonInterfaces "stori-service/src/environments/common/resources/interfaces"
)

/*
IMovementRepository to interact with entity and database
*/
type IMovementRepository interface {
	commonInterfaces.ITransactionalRepository
	BulkCreate(movements []entity.Movement) error
	FindLastMovementByCustomerID(customerid int) (*entity.Movement, error)
}

/*
	IMovementService methods with bussiness logic
*/
type IMovementService interface {
	BulkSave(movements []entity.Movement) error
}

/*
	IMovementController methods to handle requests and responses
*/
type IMovementController interface {
}
