package user

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
struct that implements IUserRepository
*/
type userGormRepo struct {
	database.TransactionalGORMRepository
}

/*
NewUserGormRepo creates a new repo and returns IMovementRepository,
so it needs to implement all its methods
*/
func NewUserGormRepo(gormDb *gorm.DB) interfaces.IUserRepository {
	rUser := &userGormRepo{}
	rUser.DB = gormDb
	return rUser
}

/*
findByUserIDAndMayLock finds a user by its UserID and locks it if lock is true
*/
func (r *userGormRepo) findByUserIDAndMayLock(userID int, lock bool) (*entity.User, error) {
	db := r.DB
	if lock {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var user entity.User
	err := db.
		Where("user_id", userID).
		Take(&user).Error
	if goerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, err
}

/*
FindAndLockByUserID finds and locks a user by its ID
*/
func (r *userGormRepo) FindAndLockByUserID(userID int) (*entity.User, error) {
	return r.findByUserIDAndMayLock(userID, true)
}

/*
FindByUserID finds a user by its ID
*/
func (r *userGormRepo) FindByUserID(userID int) (*entity.User, error) {
	return r.findByUserIDAndMayLock(userID, false)
}

/*
Clone returns a new instance of the repository
*/
func (r *userGormRepo) Clone() interface{} {
	return NewUserGormRepo(r.DB)
}
