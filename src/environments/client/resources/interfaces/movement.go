package interfaces

import (
	"stori-service/src/environments/common/resources/entity"
	commonInterfaces "stori-service/src/environments/common/resources/interfaces"
	"stori-service/src/libs/dto"
)

/*
IMovementRepository to interact with entity and database
*/
type IMovementRepository interface {
	commonInterfaces.ITransactionalRepository
	Create(movement *entity.Movement) (*entity.Movement, error)
	FindLastMovement(userID int, productID int) (*entity.Movement, error)
}

/*
	IMovementService methods with bussiness logic
*/
type IMovementService interface {
	Create(movement *dto.NewMovement) (*entity.Movement, error)
}

/*
	IMovementController methods to handle requests and responses
*/
type IMovementController interface {
}
