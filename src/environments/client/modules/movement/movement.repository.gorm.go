package stock

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
FindLastMovement finds the last stock movement of a product and a user
*/
func (r *movementGormRepo) FindLastMovement(userID int) (*entity.Movement, error) {
	var movement entity.Movement
	err := r.DB.Scopes(scopes.MovementByUserID(userID)).
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
getStockCountByUser returns the stock count of a product in a user
*/
func (r *movementGormRepo) getStockCountByUser(userID int) (int64, error) {
	var count int64
	// No encontramos como hacer que gorm nos deje hacer esto con subquery.
	err := r.DB.Raw(`SELECT COUNT(*) FROM (
		SELECT DISTINCT ON
		(movement.product_id) product.product_id,
		movement.available as stock,
		product.name,
		product.slug,
		product.description FROM "movement"
		JOIN product
		ON product.product_id = movement.product_id
		WHERE movement.user_id = ?
		AND (product.enabled = true AND product.deleted_at IS NULL)
		AND "movement"."deleted_at" IS NULL
		ORDER BY movement.product_id ASC,
		movement_id DESC
		) AS count`, userID).Scan(&count).Error
	if err != nil {
		return int64(0), err
	}
	return count, nil
}

/*
Clone returns a new instance of the repository
*/
func (r *movementGormRepo) Clone() interface{} {
	return NewMovementGormRepo(r.DB)
}
