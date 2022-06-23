package product

import (
	goerrors "errors"
	database "stori-service/src/environments/common/resources/database/transaction"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/environments/common/resources/interfaces"
	"stori-service/src/libs/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
struct that implements IProductRepository
*/
type productGormRepo struct {
	database.TransactionalGORMRepository
}

/*
NewProductGormRepo creates a new repo and returns IStockMovementRepository,
so it needs to implement all its methods
*/
func NewProductGormRepo(gormDb *gorm.DB) interfaces.IProductRepository {
	rProduct := &productGormRepo{}
	rProduct.DB = gormDb
	return rProduct
}

/*
FindAndLockByID finds and locks a product by its ID
*/
func (r *productGormRepo) FindAndLockByID(id int) (*entity.Product, error) {
	db := r.DB
	db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	var product entity.Product
	err := db.Take(&product, id).Error
	if goerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &product, err
}

/*
Clone returns a new instance of the repository
*/
func (r *productGormRepo) Clone() interface{} {
	return NewProductGormRepo(r.DB)
}
