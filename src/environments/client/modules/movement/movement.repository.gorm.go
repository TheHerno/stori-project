package movement

import (
	goerrors "errors"
	"stori-service/src/environments/client/resources/interfaces"
	database "stori-service/src/environments/common/resources/database/transaction"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/database/scopes"
	"stori-service/src/libs/errors"

	"gorm.io/gorm"
)

/*
struct that implements IMovementRepository
*/
type movementGormRepo struct {
	database.TransactionalGORMRepository
}

/*
NewMovementGormRepo creates a new repo and returns IMovementRepository,
so it needs to implement all its methods
*/
func NewMovementGormRepo(gormDb *gorm.DB) interfaces.IMovementRepository {
	rMovement := &movementGormRepo{}
	rMovement.DB = gormDb
	return rMovement
}

/*
Create receives the movement to be created and creates it
If there is an error, returns it as a second result
*/
func (r *movementGormRepo) Create(movement *entity.Movement) (*entity.Movement, error) {
	err := r.DB.Create(movement).Error
	if err != nil {
		return nil, err
	}

	return movement, nil
}

/*
BulkCreate receives a list of movements to be created and creates them
*/
func (r *movementGormRepo) BulkCreate(movements []entity.Movement) error {
	return r.DB.Create(&movements).Error
}

/*
FindLastMovementByCustomerID finds the last stock movement of a user
*/
func (r *movementGormRepo) FindLastMovementByCustomerID(customerid int) (*entity.Movement, error) {
	var movement entity.Movement
	err := r.DB.Scopes(scopes.MovementByCustomerID(customerid)).
		Order("movement_id DESC").
		Take(&movement).Error

	if goerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &movement, nil
}

/*
Clone returns a new instance of the repository
*/
func (r *movementGormRepo) Clone() interface{} {
	return NewMovementGormRepo(r.DB)
}
