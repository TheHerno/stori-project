package interfaces

import (
	"net/http"
	"stori-service/src/environments/common/resources/entity"
	commonInterfaces "stori-service/src/environments/common/resources/interfaces"
	"stori-service/src/libs/dto"
)

/*
IMovementRepository to interact with entity and database
*/
type IMovementRepository interface {
	commonInterfaces.ITransactionalRepository
	BulkCreate(movements []entity.Movement) error
	GetLastMovementByCustomerID(customerID int) (*entity.Movement, error)
}

/*
	IMovementService methods with bussiness logic
*/
type IMovementService interface {
	ProcessFile(customerID int) (*dto.MovementList, error)
}

/*
	IMovementController methods to handle requests and responses
*/
type IMovementController interface {
	ProcessFile(response http.ResponseWriter, request *http.Request)
}
