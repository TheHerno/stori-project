package movement

import (
	goerrors "errors"
	"stori-service/src/environments/client/resources/interfaces"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
)

/*
Struct that implements IMovementService
*/
type MovementService struct {
	rMovement interfaces.IMovementRepository
	rUser     interfaces.IUserRepository
}

/*
	NewMovementService creates a new service, receives repository by dependency injection
	and returns IRepositoryService, so it needs to implement all its methods
*/
func NewMovementService(rMovement interfaces.IMovementRepository, rUser interfaces.IUserRepository) interfaces.IMovementService {
	return &MovementService{rMovement, rUser}
}

/*
getAvailableStock finds the last stock movement and returns the available stock
*/
func (s *MovementService) getAvailableStock(userID int, productID int) (int, error) {
	lastMovement, err := s.rMovement.FindLastMovement(userID, productID)
	if goerrors.Is(err, errors.ErrNotFound) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return lastMovement.Available, nil
}

/*
save calculates the new available, validates it,  and saves the new stock movement
*/
func (s *MovementService) save(rMovement interfaces.IMovementRepository, newMovement *dto.NewMovement, user *entity.User) (*entity.Movement, error) {
	lastAvailable, err := s.getAvailableStock(user.UserID, newMovement.ProductID)
	if err != nil {
		return nil, err
	}

	newAvailable := lastAvailable + newMovement.Quantity*newMovement.Type
	if newAvailable < 0 {
		return nil, errors.ErrMovementInvalid
	}

	movementToCreate := &entity.Movement{
		UserID:    user.UserID,
		Quantity:  newMovement.Quantity,
		Available: newAvailable,
		Type:      newMovement.Type,
	}

	return rMovement.Create(movementToCreate)
}

/*
Create takes a newMovementDTO, validates it, locks and validate the user and product, and then create the stock movement.
*/
func (s *MovementService) Create(newMovement *dto.NewMovement) (*entity.Movement, error) {
	err := newMovement.Validate()
	if err != nil {
		return nil, err
	}

	rUser := s.rUser.Clone().(interfaces.IUserRepository)
	rMovement := s.rMovement.Clone().(interfaces.IMovementRepository)

	tx := rMovement.Begin(nil)
	rUser.Begin(tx)
	defer rMovement.Rollback()

	user, err := rUser.FindAndLockByUserID(newMovement.UserID)
	if err != nil {
		return nil, err
	}

	createdMovement, err := s.save(rMovement, newMovement, user)
	if err != nil {
		return nil, err
	}
	err = rMovement.Commit()
	if err != nil {
		return nil, err
	}
	return createdMovement, nil
}
