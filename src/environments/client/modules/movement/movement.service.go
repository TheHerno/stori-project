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
	rCustomer interfaces.ICustomerRepository
}

/*
	NewMovementService creates a new service, receives repository by dependency injection
	and returns IRepositoryService, so it needs to implement all its methods
*/
func NewMovementService(rMovement interfaces.IMovementRepository, rCustomer interfaces.ICustomerRepository) interfaces.IMovementService {
	return &MovementService{rMovement, rCustomer}
}

/*
getAvailableStock finds the last stock movement and returns the available stock
*/
func (s *MovementService) getAvailableStock(customerid int) (int, error) {
	lastMovement, err := s.rMovement.FindLastMovement(customerid)
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
func (s *MovementService) save(rMovement interfaces.IMovementRepository, newMovement *dto.NewMovement, customer *entity.Customer) (*entity.Movement, error) {
	lastAvailable, err := s.getAvailableStock(customer.CustomerID)
	if err != nil {
		return nil, err
	}

	newAvailable := lastAvailable + newMovement.Quantity*newMovement.Type
	if newAvailable < 0 {
		return nil, errors.ErrMovementInvalid
	}

	movementToCreate := &entity.Movement{
		CustomerID: customer.CustomerID,
		Quantity:   newMovement.Quantity,
		Available:  newAvailable,
		Type:       newMovement.Type,
	}

	return rMovement.Create(movementToCreate)
}

/*
Create takes a newMovementDTO, validates it, locks and validate the customer and product, and then create the stock movement.
*/
func (s *MovementService) Create(newMovement *dto.NewMovement) (*entity.Movement, error) {
	err := newMovement.Validate()
	if err != nil {
		return nil, err
	}

	rCustomer := s.rCustomer.Clone().(interfaces.ICustomerRepository)
	rMovement := s.rMovement.Clone().(interfaces.IMovementRepository)

	tx := rMovement.Begin(nil)
	rCustomer.Begin(tx)
	defer rMovement.Rollback()

	customer, err := rCustomer.FindAndLockByCustomerID(newMovement.CustomerID)
	if err != nil {
		return nil, err
	}

	createdMovement, err := s.save(rMovement, newMovement, customer)
	if err != nil {
		return nil, err
	}
	err = rMovement.Commit()
	if err != nil {
		return nil, err
	}
	return createdMovement, nil
}
