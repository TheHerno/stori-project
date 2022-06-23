package product

import (
	goerrors "errors"
	"stori-service/src/environments/admin/resources/interfaces"
	database "stori-service/src/environments/common/resources/database/transaction"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"

	"gorm.io/gorm"
)

/*
struct that implements IProductRepository
*/
type productGormRepo struct {
	database.TransactionalGORMRepository
}

/*
NewProductGormRepo creates a new repo and returns IProductRepository,
so it needs to implement all its methods
*/
func NewProductGormRepo(gormDb *gorm.DB) interfaces.IProductRepository {
	rProduct := &productGormRepo{}
	rProduct.DB = gormDb
	return rProduct
}

/*
Index receives pagination data, finds and counts activities, then returns them
If there is an error, returns it as a second result
*/
func (r *productGormRepo) Index(pagination *dto.Pagination) ([]entity.Product, error) {
	products := []entity.Product{}
	err := r.DB.Model(products).
		Count(&pagination.TotalCount).
		Offset(pagination.Offset()).
		Limit(pagination.PageSize).
		Order("name asc").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

/*
FindByID receives and id and finds the product
*/
func (r *productGormRepo) FindByID(id int) (*entity.Product, error) {
	product := &entity.Product{}
	err := r.DB.Take(product, id).Error
	if goerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return product, nil
}

/*
Update receives the product to be updated with new values
Finds by id and updates the name and address
If there is an error, returns it as a second result
*/
func (r *productGormRepo) Update(product *entity.Product) (*entity.Product, error) {
	result := r.DB.Model(product).
		Updates(map[string]interface{}{
			"name":        product.Name,
			"description": product.Description,
			"enabled":     product.Enabled,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.ErrNotFound
	}
	return product, nil
}

/*
Create receives the product to be created and creates it
If there is an error, returns it as a second result
*/
func (r *productGormRepo) Create(product *entity.Product) (*entity.Product, error) {
	err := r.DB.Create(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

/*
Delete recieves the id of the product to be deleted
and deletes it.
If there is an error, returns it as result
*/
func (r *productGormRepo) Delete(productId int) error {
	result := r.DB.Delete(&entity.Product{}, productId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

/*
Count by slug counts the number of products with the same slug
*/
func (r *productGormRepo) CountBySlug(slug string) (int64, error) {
	var count int64
	err := r.DB.Model(&entity.Product{}).
		Where("slug LIKE ?", slug+"%").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

/*
Clone returns a new instance of the repository
*/
func (r *productGormRepo) Clone() interface{} {
	return NewProductGormRepo(r.DB)
}
