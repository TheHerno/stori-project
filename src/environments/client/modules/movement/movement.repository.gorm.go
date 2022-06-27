package movement

import (
	goerrors "errors"
	"stori-service/src/environments/client/resources/interfaces"
	database "stori-service/src/environments/common/resources/database/transaction"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
BulkCreate receives a list of movements to be created and creates them
*/
func (r *movementGormRepo) BulkCreate(movements []entity.Movement) error {
	return r.DB.Create(&movements).Error
}

/*
GetLastMovementByCustomerID receives a customerID, locks the table and returns the last movement
*/
func (r *movementGormRepo) GetLastMovementByCustomerID(customerID int) (*entity.Movement, error) {
	var movement entity.Movement
	db := r.DB.Clauses(clause.Locking{Strength: "UPDATE"})
	err := db.Model(&entity.Movement{}).
		Where(&entity.Movement{CustomerID: customerID}).
		Order("date DESC").
		Take(&movement).Error
	if goerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &movement, err
}

/*
Clone returns a new instance of the repository
*/
func (r *movementGormRepo) Clone() interface{} {
	return NewMovementGormRepo(r.DB)
}
