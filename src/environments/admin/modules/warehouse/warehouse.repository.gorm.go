package warehouse

import (
	"stori-service/src/environments/admin/resources/interfaces"
	database "stori-service/src/environments/common/resources/database/transaction"
	"stori-service/src/environments/common/resources/entity"
	"stori-service/src/libs/dto"
	"stori-service/src/libs/errors"
	"strings"

	"gorm.io/gorm"
)

/*
struct that implements IWarehouseRepository
*/
type warehouseGormRepo struct {
	database.TransactionalGORMRepository
}

/*
NewWarehouseGormRepo creates a new repo and returns IWarehouseRepository,
so it needs to implement all its methods
*/
func NewWarehouseGormRepo(gormDb *gorm.DB) interfaces.IWarehouseRepository {
	rWarehouse := &warehouseGormRepo{}
	rWarehouse.DB = gormDb
	return rWarehouse
}

/*
Index receives pagination data, finds and counts activities, then returns them
If there is an error, returns it as a second result
*/
func (r *warehouseGormRepo) Index(pagination *dto.Pagination) (*[]entity.Warehouse, error) {
	warehouses := &[]entity.Warehouse{}
	err := r.DB.Model(warehouses).
		Count(&pagination.TotalCount).
		Offset(pagination.Offset()).
		Limit(pagination.PageSize).
		Order("name asc").
		Find(&warehouses).Error
	if err != nil {
		return nil, err
	}
	return warehouses, nil
}

/*
Update receives the warehouse to be updated with new values
Finds by id and updates the name and address
If there is an error, returns it as a second result
*/
func (r *warehouseGormRepo) Update(warehouse *entity.Warehouse) (*entity.Warehouse, error) {
	result := r.DB.Model(warehouse).
		Updates(map[string]interface{}{
			"name":    warehouse.Name,
			"address": warehouse.Address,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.ErrNotFound
	}
	return warehouse, nil
}

/*
Create receives the warehouse to be created and creates it
If there is an error, returns it as a second result
*/
func (r *warehouseGormRepo) Create(warehouse *entity.Warehouse) (*entity.Warehouse, error) {
	err := r.DB.Create(warehouse).Error
	if err != nil {
		if e := err.Error(); strings.Contains(e, "23505") {
			err = errors.ErrUserIDDuplicated
		}
		return nil, err
	}

	return warehouse, nil
}

/*
Delete recieves the id of the warehouse to be deleted
and deletes it.
If there is an error, returns it as result
*/
func (r *warehouseGormRepo) Delete(warehouseId int) error {
	result := r.DB.Delete(&entity.Warehouse{}, warehouseId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

/*
Clone returns a new instance of the repository
*/
func (r *warehouseGormRepo) Clone() interface{} {
	return NewWarehouseGormRepo(r.DB)
}
